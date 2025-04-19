package middleware

import (
    "log"
    "net/http"
    "time"
)

// Logger middleware
func Logger(next func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        timestamp := time.Now().Format("2006-01-02 15:04:05")
        log.Printf("[%s] %s -> %s", r.Method, r.URL.Path, timestamp)
        next(w, r)
    }
}


