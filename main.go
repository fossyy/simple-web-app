package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"os"

	"github.com/google/uuid"
)

type writerWrapper struct {
	http.ResponseWriter
	statusCode int
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	addr := fmt.Sprintf("localhost:%s", port)
	server := newServer(addr, middleware(setupRoute()))

	fmt.Printf("Listening on http://%s\n", addr)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting server:", err)
		return
	}
}

func newServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}

func setupRoute() *http.ServeMux {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		data := map[string]interface{}{
			"Request ID": uuid.New().String(),
			"Timestamp":  time.Now(),
		}
		json.NewEncoder(w).Encode(data)
	})
	return handler
}

func (w writerWrapper) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.statusCode = code
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := writerWrapper{w, http.StatusOK}
		log.Printf(fmt.Sprintf("%s %s %v \n", r.Method, r.RequestURI, wrapped.statusCode))
		next.ServeHTTP(wrapped, r)
	})
}
