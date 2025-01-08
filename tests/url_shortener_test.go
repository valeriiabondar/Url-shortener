package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gavv/httpexpect/v2"

	"urlShortener/internal/http-server/handlers"
	"urlShortener/internal/utils/random"
)

const (
	host = "localhost:8080"
)

func TestUrlShortenerSave(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	e.POST("/url/").
		WithJSON(handlers.Request{
			Url:   gofakeit.URL(),
			Alias: random.NewRandomAlias(10),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(http.StatusCreated).
		JSON().Object().
		ContainsKey("alias")
}

func TestUrlShortenerSaveRedirect(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		alias       string
		status      int
		error       string
		aliasExists bool
		urlNotFound bool
	}{
		{
			name:   "Full",
			url:    "https://example.com",
			alias:  random.NewRandomAlias(10),
			status: http.StatusCreated,
		},
		{
			name:   "Empty alias",
			url:    "https://example.com",
			alias:  "",
			status: http.StatusCreated,
		},
		{
			name:        "Existing alias",
			url:         gofakeit.URL(),
			alias:       random.NewRandomAlias(10),
			aliasExists: true,
			error:       "alias already exists",
			status:      http.StatusConflict,
		},
		{
			name:   "Invalid Url",
			url:    "invalid url",
			alias:  random.NewRandomAlias(10),
			error:  "validation error: [field Url is not a valid Url]",
			status: http.StatusBadRequest,
		},
		{
			name:        "Url not found",
			url:         gofakeit.URL(),
			alias:       random.NewRandomAlias(10),
			error:       "url not found",
			status:      http.StatusNotFound,
			urlNotFound: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			if tc.aliasExists {
				e.POST("/url/").
					WithJSON(handlers.Request{
						Url:   tc.url,
						Alias: tc.alias,
					}).
					WithBasicAuth("admin", "admin").
					Expect().
					Status(http.StatusCreated).
					JSON().Object().ContainsKey("alias")
			}

			if tc.urlNotFound {
				e.GET(fmt.Sprintf("/%s", tc.alias)).
					Expect().
					Status(http.StatusNotFound)
				return
			}

			// Save
			resp := e.POST("/url/").
				WithJSON(handlers.Request{
					Url:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("admin", "admin").
				Expect().
				Status(tc.status).
				JSON().Object()

			if tc.error != "" {
				resp.NotContainsKey("alias")

				resp.Value("error").String().IsEqual(tc.error)

				return
			}

			var alias string
			if tc.alias != "" {
				resp.Value("alias").String().IsEqual(tc.alias)
				alias = tc.alias
			} else {
				resp.Value("alias").String().NotEmpty()
				alias = resp.Value("alias").String().Raw()
			}

			// Redirect
			e.GET(fmt.Sprintf("/%s", alias)).
				Expect().
				Status(http.StatusOK)
		})
	}
}

func TestUrlShortenerSaveDelete(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	// Save
	resp := e.POST("/url/").
		WithJSON(handlers.Request{
			Url:   gofakeit.URL(),
			Alias: random.NewRandomAlias(10),
		}).
		WithBasicAuth("admin", "admin").
		Expect().
		Status(http.StatusCreated).
		JSON().Object().ContainsKey("alias")

	// Delete
	alias := resp.Value("alias").String().Raw()

	e.DELETE(fmt.Sprintf("/url/%s", alias)).
		WithBasicAuth("admin", "admin").
		Expect().Status(http.StatusNoContent)
}
