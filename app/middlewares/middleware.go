package middlewares

import "net/http"

type TMiddleware func(http.HandlerFunc) http.HandlerFunc

func Chain(f http.HandlerFunc, ms ...TMiddleware) http.HandlerFunc  {
	for _, m := range ms {
		f = m(f)
	}
	return f
}
