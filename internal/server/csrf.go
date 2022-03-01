package server

import (
	"context"
	"log"
	"net/http"

	"github.com/sgblanch/pathview-web/internal/config"
	"github.com/sgblanch/pathview-web/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

func CSRF() gin.HandlerFunc {
	cookieName := "pathview-csrf"
	if !gin.IsDebugging() {
		cookieName = "__Host-" + cookieName
	}

	handlers := new(csrfMW)

	csrf_key, err := util.ParseKey(config.Get().CSRFKey, 32)
	if err != nil {
		log.Fatalf("csrf key: %q", err)
	}

	CSRF := csrf.Protect(
		csrf_key,
		csrf.CookieName(cookieName),
		csrf.FieldName("csrf-token"),
		csrf.MaxAge(0), // Session only
		csrf.ErrorHandler(http.HandlerFunc(handlers.errorHandler)),
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Secure(!gin.IsDebugging()),
	)

	return func(c *gin.Context) {
		state := &csrfCtx{ctx: c}
		ctx := context.WithValue(c.Request.Context(), handlers, state)
		CSRF(handlers).ServeHTTP(c.Writer, c.Request.WithContext(ctx))

		if !state.called {
			c.Abort()
		}
	}
}

type csrfMW struct{}

func (p *csrfMW) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	state := r.Context().Value(p).(*csrfCtx)
	state.ctx.Request = r
	w.Header().Set("X-CSRF-Token", csrf.Token(r))
	state.called = true
	state.ctx.Next()
}

func (p *csrfMW) errorHandler(w http.ResponseWriter, r *http.Request) {
	state := r.Context().Value(p).(*csrfCtx)
	state.failed = true
	state.ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": csrf.FailureReason(r)})
}

type csrfCtx struct {
	ctx    *gin.Context
	called bool
	failed bool
}
