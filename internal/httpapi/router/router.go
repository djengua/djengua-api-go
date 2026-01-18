package router

import (
	"net/http"
	"time"

	"github.com/djengua/djengua-api-go/internal/application"
	"github.com/djengua/djengua-api-go/internal/httpapi/handlers"
	appmw "github.com/djengua/djengua-api-go/internal/httpapi/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(productsService *application.ProductsService, categoriesService *application.CategoriesService, collectionsService *application.CollectionsService) http.Handler {
	r := chi.NewRouter()

	// --- middleware ---
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(30 * time.Second))

	// Prometheus metrics (custom hist/counter)
	appmw.RegisterMetrics(prometheus.DefaultRegisterer)
	r.Use(appmw.PrometheusMetrics)

	// Health endpoints
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	// Prometheus scrape endpoint
	r.Handle("/metrics", promhttp.Handler())

	ph := handlers.NewProductsHandler(productsService)
	ch := handlers.NewCategoriesHandler(categoriesService)
	coh := handlers.NewCollectionsHandler(collectionsService)

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
