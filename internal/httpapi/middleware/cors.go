package middleware

import "net/http"

func CORS(next http.Handler, origins string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 🔐 Orígenes permitidos
		w.Header().Set("Access-Control-Allow-Origin", origins)
		// Si usas cookies o auth por headers específicos:
		// w.Header().Set("Access-Control-Allow-Origin", "https://tu-frontend.com")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Preflight request
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
