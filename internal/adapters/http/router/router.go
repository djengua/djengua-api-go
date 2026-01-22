// internal/adapters/http/router/router.go
package router

import (
	"net/http"
	"time"

	"github.com/djengua/djengua-api-go/internal/adapters/http/handlers"
	appmw "github.com/djengua/djengua-api-go/internal/adapters/http/middleware"
	"github.com/djengua/djengua-api-go/internal/core/ports"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New(
	productsUC ports.ProductService,
	categoriesUC ports.CategoryService,
	collectionsUC ports.CollectionService,
	authUC ports.AuthService,
	ordersUC ports.OrderService,
	salesUC ports.SaleService,
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

	ah := handlers.NewAuthHandler(authUC)
	ph := handlers.NewProductsHandler(productsUC)
	cath := handlers.NewCategoriesHandler(categoriesUC)
	coh := handlers.NewCollectionsHandler(collectionsUC)

	// ✅ NUEVOS handlers
	oh := handlers.NewOrdersHandler(ordersUC)
	sh := handlers.NewSalesHandler(salesUC)

	r.Route("/api/v1", func(api chi.Router) {

		api.Route("/auth", func(rr chi.Router) {
			rr.Post("/register", ah.Register)
			rr.Post("/login", ah.Login)

			rr.Group(func(p chi.Router) {
				p.Use(appmw.RequireJWT(authUC))
				p.Get("/me", ah.Me)
			})
		})

		// -------- PRODUCTS --------
		api.Route("/products", func(rr chi.Router) {
			rr.Get("/", ph.List)
			rr.Get("/{id}", ph.Get)

			rr.Group(func(pr chi.Router) {
				pr.Use(appmw.RequireJWT(authUC))
				pr.Post("/", ph.Create)
				pr.Put("/{id}", ph.Put)
				pr.Patch("/{id}", ph.Patch)
				pr.Delete("/{id}", ph.Delete)
			})
		})

		// -------- CATEGORIES --------
		api.Route("/categories", func(rr chi.Router) {
			rr.Get("/", cath.List)
			rr.Get("/{id}", cath.Get)

			rr.Group(func(cr chi.Router) {
				cr.Use(appmw.RequireJWT(authUC))
				cr.Post("/", cath.Create)
				cr.Put("/{id}", cath.Put)
				cr.Patch("/{id}", cath.Patch)
				cr.Delete("/{id}", cath.Delete)
			})
		})

		// -------- COLLECTIONS --------
		api.Route("/collections", func(rr chi.Router) {
			rr.Get("/", coh.List)
			rr.Get("/{id}", coh.Get)

			rr.Group(func(cor chi.Router) {
				cor.Use(appmw.RequireJWT(authUC))
				cor.Post("/", coh.Create)
				cor.Put("/{id}", coh.Put)
				cor.Patch("/{id}", coh.Patch)
				cor.Delete("/{id}", coh.Delete)
			})
		})

		// -------- ORDERS (JWT) --------
		api.Route("/orders", func(rr chi.Router) {
			rr.Group(func(or chi.Router) {
				or.Use(appmw.RequireJWT(authUC))
				or.Post("/", oh.Create)
				or.Get("/me", oh.MyOrders) // listar mis órdenes
				// opcional:
				// or.Get("/{id}", oh.Get)
			})
		})

		// -------- SALES (JWT) --------
		api.Route("/sales", func(rr chi.Router) {
			rr.Group(func(sr chi.Router) {
				sr.Use(appmw.RequireJWT(authUC))
				sr.Post("/", sh.Register)
			})
		})
	})

	return r
}
