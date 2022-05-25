package run

import (
	"net/http"
	"os"
	"path"
	"strings"
)

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

// See https://gist.github.com/lummie/91cd1c18b2e32fa9f316862221a6fd5c

// FSHandler404 provides the function signature for passing to the FileServerWith404
type FSHandler404 = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool)

// FileServerWith404 wraps the http.FileServer checking to see if the url path exists first.
// If the file fails to exist it calls the supplied FSHandle404 function
// The implementation can choose to either modify the request, e.g. change the URL path and return true to have the
// default FileServer handling to still take place, or return false to stop further processing, for example if you wanted
// to write a custom response
// e.g. redirects to root and continues the file serving handler chain
// 	func fileSystem404(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool) {
// 		//if not found redirect to /
// 		r.URL.Path = "/"
// 		return true
// 	}
// Use the same as you would with a http.FileServer e.g.
// 	r.Handle("/", http.StripPrefix("/", mw.FileServerWith404(http.Dir("./staticDir"), fileSystem404)))
func FileServerWith404(root http.FileSystem, handler404 FSHandler404) http.Handler {
	fs := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// make sure the url path starts with /
		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)

		// attempt to open the file via the http.FileSystem
		f, err := root.Open(upath)
		if err != nil {
			if os.IsNotExist(err) {
				// call handler
				if handler404 != nil {
					doDefault := handler404(w, r)
					if !doDefault {
						return
					}
				}
			}
		}

		// close if successfully opened
		if err == nil {
			f.Close()
		}

		// default serve
		fs.ServeHTTP(w, r)
	})
}
