package tests

import (
	"net/url"
	"testing"

	"github.com/DimsFromDergachy/Url-Shortener/internal/http-server/handlers/url/save"
	"github.com/DimsFromDergachy/Url-Shortener/internal/lib/random"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

const (
	host = "localhost:8082"
)

func TestURLShortener_HappyPath(t *testing.T) {
	u := url.URL {
		Scheme: "http",
		Host:   host,
	}

	e := httpexpect.Default(t, u.String())

	e.POST("/").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		// WithBasicAuth("myuser", "mypass").
		Expect().
		Status(200).
		JSON().Object().
		ContainsKey("alias")
}