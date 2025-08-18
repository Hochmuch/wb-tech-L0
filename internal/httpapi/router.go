package httpapi

import (
	"html/template"
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /order/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.GetOrder(w, r)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Query().Get("order_uid") == "" {
			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			tmpl.Execute(w, nil)
			return
		}

		orderUID := r.URL.Query().Get("order_uid")
		http.Redirect(w, r, "/order/"+orderUID, http.StatusFound)
	})

	return mux
}
