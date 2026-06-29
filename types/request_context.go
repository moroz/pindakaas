package types

import (
	"net/http"

	"github.com/moroz/pindakaas/config"
	"github.com/moroz/pindakaas/db/queries"
	"github.com/moroz/pindakaas/web/sessions"
)

type RequestContext struct {
	store   *sessions.Store
	Session sessions.Payload
	User    *queries.User
}

func NewRequestContext(store *sessions.Store) *RequestContext {
	return &RequestContext{
		store: store,
	}
}

func (c *RequestContext) SaveSession(w http.ResponseWriter) error {
	cookie, err := c.store.EncodeSession(c.Session)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.SessionCookieName,
		Value:    cookie,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
	})
	return nil
}
