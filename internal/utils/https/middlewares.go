package https

import "net/http"

const ContentTypeJson string = "application/json"

func NewMiddlewareContentType(contentType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", contentType)

			next.ServeHTTP(w, r)
		})
	}
}
