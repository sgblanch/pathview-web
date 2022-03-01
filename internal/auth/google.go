package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/sgblanch/pathview-web/internal/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func Google(g *gin.RouterGroup) {
	c := config.Get().Google
	google := &Config{
		Config: &oauth2.Config{
			ClientID:     c.ClientID,
			ClientSecret: c.ClientSecret,
			RedirectURL:  c.Redirect,
			Scopes: []string{
				"openid",
				"email",
				"profile",
			},
			Endpoint: google.Endpoint,
		},
		OpenID: New("https://accounts.google.com/.well-known/openid-configuration"),
	}

	g.GET("google", google.Login)
	g.GET("google/callback", google.Callback)
}
