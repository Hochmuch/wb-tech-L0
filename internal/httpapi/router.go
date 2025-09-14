package httpapi

import (
	"html/template"
	"log"
	"net/http"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /order/{id}", func(w http.ResponseWriter, r *http.Request) {
		h.GetOrder(w, r)
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			return
		}
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Query().Get("order_uid") == "" {
			tmpl := template.Must(template.ParseFiles("templates/index.html"))
			err := tmpl.Execute(w, nil)
			if err != nil {
				log.Println(err)
				return
			}
			return
		}

		orderUID := r.URL.Query().Get("order_uid")
		http.Redirect(w, r, "/order/"+orderUID, http.StatusFound)
	})

	return mux
}
