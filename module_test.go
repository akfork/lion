package lion

import (
	"net/http"
	"testing"

	"golang.org/x/net/context"
)

type testmodule struct {
	base string
}

func (m testmodule) Routes(r *Router) {

}

func (m testmodule) Base() string {
	return m.base
}

func (m testmodule) Requires() []string {
	return []string{"auth", "jwt"}
}

func (m testmodule) Uses() (mws Middlewares) {
	return mws
}

func (m testmodule) Get(c context.Context, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("getmodule"))
}

func TestModule(t *testing.T) {
	l := New()
	l.DefineFunc("auth", func(next Handler) Handler {
		return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			w.Header().Set("auth", "authmw")
			next.ServeHTTPC(c, w, r)
		})
	})

	l.DefineFunc("jwt", func(next Handler) Handler {
		return HandlerFunc(func(c context.Context, w http.ResponseWriter, r *http.Request) {
			w.Header().Set("token", "jwtmw")
			next.ServeHTTPC(c, w, r)
		})
	})

	l.Module(testmodule{"/admin"})

	expectHeader(t, l, "GET", "/admin", "auth", "authmw")
	expectHeader(t, l, "GET", "/admin", "token", "jwtmw")
	expectBody(t, l, "GET", "/admin", "getmodule")
}
