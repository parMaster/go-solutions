// main.go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type contextKey string

const (
	requestIDKey  contextKey = "requestID"
	companyIdKey  contextKey = "companyID"
	CustomDataKey contextKey = "customData"
)

// ResponseWriter wrapper to capture the response
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter, r *http.Request) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (key contextKey) Read(ctx context.Context) any {
	data := ctx.Value(CustomDataKey)
	if data == nil {
		return nil
	}

	if value, ok := data.(map[contextKey]any)[key]; ok {
		return value
	}
	return nil
}

func (key contextKey) Write(ctx context.Context, value any) context.Context {
	data := ctx.Value(CustomDataKey)
	if data == nil {
		data = map[contextKey]any{}
	}

	data.(map[contextKey]any)[key] = value
	return context.WithValue(ctx, CustomDataKey, data)
}

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Before Logger: %p\n", r)
		start := time.Now()

		ctx := context.WithValue(r.Context(), CustomDataKey, map[contextKey]any{})
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		fmt.Printf("After Logger: %p\n", r)

		// Get custom data from context
		if companyID, ok := companyIdKey.Read(r.Context()).(int); ok {
			log.Printf("CompanyID in logger (after): %d", companyID)
		} else {
			log.Printf("No companyID in logger (after)")
		}

		log.Printf("Request took: %v", time.Since(start))
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Before Auth: %p\n", r)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)

		if companyID, ok := companyIdKey.Read(r.Context()).(int); ok {
			log.Printf("CompanyID in auth (after): %d", companyID)
		} else {
			log.Printf("No companyID in auth (after)")
		}

		fmt.Printf("After Auth: %p\n", r)
	})
}

func ApiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("In ApiHandler before setting context: %p\n", r)

	// Get existing context and add company ID
	if data, ok := r.Context().Value(CustomDataKey).(map[contextKey]any); ok {
		data[companyIdKey] = 5
	}

	fmt.Printf("In ApiHandler after setting context: %p\n", r)
	if companyID, ok := companyIdKey.Read(r.Context()).(int); ok {
		log.Printf("CompanyID in handler (after setting): %d", companyID)
	}

	fmt.Fprintf(w, "Request processed\n")
}

func main() {
	apiHandler := http.HandlerFunc(ApiHandler)

	handler := LoggerMiddleware(
		AuthMiddleware(
			apiHandler,
		),
	)

	http.Handle("/api", handler)
	log.Println("Server starting on :8080...")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
