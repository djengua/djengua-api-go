package router

import (
	"net/http"
	"time"

	"github.com/djengua/djengua-api-go/internal/adapters/http/handlers"
	appmw "github.com/djengua/djengua-api-go/internal/adapters/http/middleware"
	"github.com/djengua/djengua-api-go/internal/core/usecase/categories"
	"github.com/djengua/djengua-api-go/internal/core/usecase/collections"
	"github.com/djengua/djengua-api-go/internal/core/usecase/products"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(
	productsUC *products.Service,
	categoriesUC *categories.Service,
	collectionsUC *collections.Service,
) http.Handler {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))

	appmw.RegisterMetrics(prometheus.DefaultRegisterer)
	r.Use(appmw.PrometheusMetrics)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	r.Handle("/metrics", promhttp.Handler())

	ph := handlers.NewProductsHandler(productsUC)
	ch := handlers.NewCategoriesHandler(categoriesUC)
	coh := handlers.NewCollectionsHandler(collectionsUC)

	r.Route("/api/v1", func(api chi.Router) {
		api.Route("/products", func(rr chi.Router) {
			rr.Get("/", ph.List)
			rr.Post("/", ph.Create)
			rr.Route("/{id}", func(rid chi.Router) {
				rid.Get("/", ph.Get)
				rid.Put("/", ph.Put)
				rid.Patch("/", ph.Patch)
				rid.Delete("/", ph.Delete)
			})
		})

		api.Route("/categories", func(rr chi.Router) {
			rr.Get("/", ch.List)
			rr.Post("/", ch.Create)
			rr.Route("/{id}", func(rid chi.Router) {
				rid.Get("/", ch.Get)
				rid.Put("/", ch.Put)
				rid.Patch("/", ch.Patch)
				rid.Delete("/", ch.Delete)
			})
		})

		api.Route("/collections", func(rr chi.Router) {
			rr.Get("/", coh.List)
			rr.Post("/", coh.Create)
			rr.Route("/{id}", func(rid chi.Router) {
				rid.Get("/", coh.Get)
				rid.Put("/", coh.Put)
				rid.Patch("/", coh.Patch)
				rid.Delete("/", coh.Delete)
			})
		})
	})

	return r
}
