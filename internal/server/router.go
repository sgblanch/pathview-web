package server

import (
	"fmt"
	"net/http"

	"github.com/sgblanch/pathview-web/internal/auth"
	"github.com/sgblanch/pathview-web/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type Router struct {
	config *config.Config
}

func (p *Router) sessionSetup() gin.HandlerFunc {
	store := p.config.RedisStore
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   2592000,
		Secure:   !gin.IsDebugging(),
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	})

	sessionName := "pathview"
	if !gin.IsDebugging() {
		sessionName = "__Host-" + sessionName
	}

	return sessions.Sessions(sessionName, store)
}

func (p *Router) securitySetup() gin.HandlerFunc {
	security := secure.DefaultConfig()
	security.ContentSecurityPolicy = "default-src 'self'; frame-ancestors 'none'; img-src 'self' data:; script-src 'unsafe-eval' 'self';"
	security.SSLRedirect = !gin.IsDebugging()
	security.IsDevelopment = false

	return secure.New(security)
}

func (p *Router) headerSetup() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Writer.Header().Set("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), battery=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=()")
	}
}

func (p *Router) Router() *gin.Engine {
	p.config = config.Get()

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(p.securitySetup())
	router.Use(p.sessionSetup())
	router.Use(p.headerSetup())
	// router.Use(CSRF())

	router.SetFuncMap(funcMap())
	router.LoadHTMLGlob("template/*.go.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("count")
		var count int
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()
		c.String(http.StatusOK, fmt.Sprintf("Welcome Gin Server: %v", count))
	})

	pathview := router.Group("pathview")
	{
		p.pathviewGroup(pathview)
	}

	v1 := router.Group("api/v1")
	{
		kegg := v1.Group("kegg")
		{
			p.KeggGroup(kegg)
		}
	}

	login := router.Group("login")
	{
		// g.GET("/", func(c *gin.Context) {
		// 	c.HTML(200, "login.go.html", pathview)
		// })
		auth.Google(login)
	}

	return router
}
