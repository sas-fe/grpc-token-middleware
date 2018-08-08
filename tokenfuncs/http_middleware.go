package tokenfuncs

import (
	"net/http"

	"google.golang.org/grpc/metadata"
)

type mdIncomingKey struct{}

// Middleware is an http middleware for tokenfuncs in lieu of gRPC
func (t *TokenFuncs) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiToken := r.Header.Get("Authorization")
		if len(apiToken) == 0 {
			http.Error(w, "Authorization header doesn't exist", http.StatusTooManyRequests)
			return
		}

		md := map[string][]string{
			"authorization": []string{apiToken},
		}

		ctx := metadata.NewIncomingContext(r.Context(), md)

		allowed, err := t.CheckValidity(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if !allowed {
			http.Error(w, "API limit exceeded", http.StatusTooManyRequests)
			return
		}

		if t.IsAsync() {
			go t.AsyncIncrementUsage(ctx)
		} else {
			_, err = t.IncrementUsage(ctx)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
