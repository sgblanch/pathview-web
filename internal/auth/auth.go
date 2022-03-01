package auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/sgblanch/pathview-web/internal/model"
	"github.com/sgblanch/pathview-web/internal/util"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/MicahParks/keyfunc"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func init() {
	jwt.TimeFunc = func() time.Time {
		return time.Now().UTC().Add(time.Second * 5)
	}
}

var _once sync.Once

type Config struct {
	*oauth2.Config
	*OpenID
}

func (p *Config) Login(c *gin.Context) {
	session := sessions.Default(c)
	token, err := util.RandomToken(32)
	if err != nil {
		log.Printf("[callback] generating token: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	session.Set("oauth-state-token", token)
	err = session.Save()
	if err != nil {
		log.Printf("[callback] save session: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	url := p.AuthCodeURL(*token)
	c.Redirect(http.StatusFound, url)
}

func (p *Config) Callback(c *gin.Context) {
	var validate string

	session := sessions.Default(c)
	if v := session.Get("oauth-state-token"); v == nil {
		log.Printf("Unable to retrieve oauth-state-token from session")
		c.AbortWithStatus(http.StatusInternalServerError)
	} else {
		validate = v.(string)
	}

	state := c.Request.FormValue("state")
	if state != validate {
		log.Print("[callback] unable to validate state")
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "state token failed validation"})
		return
	}

	code := c.Request.FormValue("code")
	if code == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": c.Request.FormValue("error_reason")})
		return
	}

	token, err := p.Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err})
		return
	}

	id_token, err := jwt.Parse(token.Extra("id_token").(string), p.JWKS.Keyfunc)

	claims, ok := id_token.Claims.(jwt.MapClaims)
	if !ok || !id_token.Valid {
		log.Printf("[callback] parsing %v", err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message:": "um, nope"})
		return
	}

	user := new(model.User)
	err = user.Login(claims)
	if err != nil {
		log.Printf("[callback] login user: %v", err)
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message:": "um, nope"})
		return
	}
	session.Set("user", *user)
	err = session.Save()
	if err != nil {
		log.Printf("[callback] save session: %v", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Redirect(http.StatusFound, "/")
}

type OpenID struct {
	URL    string
	Config gin.H
	JWKS   *keyfunc.JWKS
}

func New(url string) *OpenID {
	openid := &OpenID{
		URL: url,
	}

	r, err := http.Get(openid.URL)
	cobra.CheckErr(err)
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&openid.Config)
	cobra.CheckErr(err)

	jwks := openid.Config["jwks_uri"].(string)
	openid.JWKS, err = keyfunc.Get(jwks, keyfunc.Options{})
	cobra.CheckErr(err)

	return openid
}
