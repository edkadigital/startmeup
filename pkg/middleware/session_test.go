package middleware

import (
	"testing"

	"github.com/edkadigital/startmeup/pkg/session"
	"github.com/edkadigital/startmeup/pkg/tests"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSession(t *testing.T) {
	ctx, _ := tests.NewContext(c.Web, "/")
	_, err := session.Get(ctx, "test")
	assert.Equal(t, session.ErrStoreNotFound, err)

	store := sessions.NewCookieStore([]byte("secret"))
	err = tests.ExecuteMiddleware(ctx, Session(store))
	require.NoError(t, err)

	_, err = session.Get(ctx, "test")
	assert.NotEqual(t, session.ErrStoreNotFound, err)
}
