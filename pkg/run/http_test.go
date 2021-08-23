package run

import (
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestInterceptHttp(t *testing.T) {
	t.Run("in correct order", func(t *testing.T) {
		items := &collector{}
		handler := InterceptHttp(
			handler(items.addValue("root")),
			interceptor(items.addValue("one")),
			interceptor(items.addValue("two")),
			interceptor(items.addValue("three")),
		)
		handler.ServeHTTP(nil, nil) // don't care about values
		items.want(t, "one", "two", "three", "root")
	})
}

func handler(ran func()) http.Handler {
	return http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		ran()
	})
}

func interceptor(ran func()) HttpInterceptor {
	return func(w http.ResponseWriter, r *http.Request, next http.Handler) {
		ran()
		next.ServeHTTP(w, r)
	}
}

type collector []string

func (c *collector) addValue(val string) func() {
	return func() {
		*c = append(*c, val)
	}
}

func (c collector) want(t *testing.T, want ...string) {
	if diff := cmp.Diff(want, ([]string)(c)); diff != "" {
		t.Errorf("unexpected items (-want, +got)\n%v", diff)
	}
}
