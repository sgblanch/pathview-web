package model

import (
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/sgblanch/pathview-web/internal/config"

	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID            int64     `db:"id" json:"-" `
	Name          string    `db:"name" json:"name,omitempty"` // TODO: Nullable
	Email         string    `db:"email" json:"email"`
	EmailVerified bool      `db:"email_verified" json:"-"` // TODO: Nullable
	Provider      string    `db:"provider" json:"-"`
	ProviderID    string    `db:"provider_id" json:"-"`
	CreatedAt     time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt     time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

func init() {
	gob.Register(User{})
}

func (p *User) Login(claims jwt.MapClaims) error {

	var (
		ok   bool
		err  error
		user = new(User)
	)

	user.Provider, ok = claims["iss"].(string)
	if !ok {
		return fmt.Errorf("no issuer found")
	}

	user.ProviderID, ok = claims["sub"].(string)
	if !ok {
		return fmt.Errorf("no subject found")
	}

	// Optional
	user.Name = claims["name"].(string)

	// Technically Optional, but requiring
	if user.Email, ok = claims["email"].(string); !ok {
		return fmt.Errorf("no email found")
	}

	// Optional
	user.EmailVerified = claims["email_verified"].(bool)

	err = config.Get().DB.NamedGet(p, _sql["user_select"], user)
	if errors.Is(err, sql.ErrNoRows) {
		p = user
		return p.Create()
	} else if err != nil {
		return err
	}

	var update bool

	if p.Email != user.Email {
		p.Email = user.Email
		update = true
	}
	if p.EmailVerified != user.EmailVerified {
		p.EmailVerified = user.EmailVerified
		update = true
	}
	if p.Name != user.Name {
		p.Name = user.Name
		update = true
	}

	if update {
		return p.Update()
	}

	return nil
}

func (p *User) Create() error {
	return config.Get().DB.NamedGet(p, _sql["user_create"], p)
}

func (p *User) Update() error {
	return config.Get().DB.NamedGet(p, _sql["user_update"], p)
}
