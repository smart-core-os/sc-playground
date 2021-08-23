package run

import "net/http"

type HttpInterceptor func(w http.ResponseWriter, r *http.Request, next http.Handler)

func InterceptHttp(underlying http.Handler, interceptors ...HttpInterceptor) http.Handler {
	root := underlying
	for i := len(interceptors) - 1; i >= 0; i-- {
		interceptor := interceptors[i]
		lastRoot := root
		root = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			interceptor(writer, request, lastRoot)
		})
	}
	return root
}
