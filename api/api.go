package api

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/parish/notes/store"
)

//New configures application resources and routers.
func New(enableCors bool, store store.Store, logger *log.Logger) (*chi.Mux, error) {
	noteResource := NewNoteResource(store, logger)

	r := chi.NewRouter()
	r.Use(contentTypeMiddleware("application/json"))
	r.Use(loggingMiddleware(logger))

	if enableCors {
		r.Use(corsConfig().Handler)
	}

	r.Mount("/api/v1/", noteResource.Router())

	r.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "User-agent: *\nDisallow: /")
	})
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("%s %s\n", method, route)
		return nil
	}

	if err := chi.Walk(r, walkFunc); err != nil {
		fmt.Printf("Logging err:%s\n", err.Error())
	}

	return r, nil
}

// contentType is a middleware that forces response Content-Type.
func contentTypeMiddleware(contentType string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-type", "application/json")
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

// LoggingMiddleware logs the incoming HTTP request & its duration.
func loggingMiddleware(logger *log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					logger.Printf(
						"err %s trace %s", err, debug.Stack(),
					)
				}
			}()

			start := time.Now()
			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			logger.Println(
				"status", wrapped.status,
				"method", r.Method,
				"path", r.URL.EscapedPath(),
				"duration", time.Since(start),
			)
		}

		return http.HandlerFunc(fn)
	}
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

// responseWriter is a minimal wrapper for http.ResponseWriter that allows the
// written HTTP status code to be captured for logging.
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func corsConfig() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Authorization", "Accept-Encoding", "Content-Type", "Content-Length", "X-CSRF-Token"},
		ExposedHeaders:     []string{"Link"},
		AllowCredentials:   true,
		MaxAge:             86400,
		OptionsPassthrough: false,
	})
}
