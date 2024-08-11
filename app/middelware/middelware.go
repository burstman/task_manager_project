package middelware

import (
	"context"
	"net/http"
	"taskManager/app/types"
)

func HtmxMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isHtmxRequest := r.Header.Get("HX-Request") != ""
		ctx := context.WithValue(r.Context(), types.IsHtmxRequestKey, isHtmxRequest)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
