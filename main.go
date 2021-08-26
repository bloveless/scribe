package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s %v %v\n", r.RemoteAddr, r.Method, r.URL, r.Header, r.URL.Query())
		handler.ServeHTTP(w, r)
	})
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		type request struct {
			RequestType string `json:"type"`
			Token       string `json:"token"`
			Challenge   string `json:"challenge"`
		}

		str, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return
		}

		log.Printf("Request body: %s", string(str))

		var req request
		err = json.NewDecoder(strings.NewReader(string(str))).Decode(&req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Unable to decode request body"))
			return
		}

		log.Printf("Request body: %+v", req)

		w.Write([]byte(req.Challenge))
	})

	srv := &http.Server{
		Addr:              ":8080",
		ReadTimeout:       1 * time.Second,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
		Handler:           r,
	}

	log.Printf("Listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Server failed: %s\n", err)
	}
}
