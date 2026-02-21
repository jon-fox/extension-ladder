package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestApi_ValidURL(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><p>Mock API content</p></body></html>`))
	}))
	defer mockServer.Close()

	app := fiber.New()
	app.Get("/api/*", Api)

	req := httptest.NewRequest(http.MethodGet, "/api/"+mockServer.URL, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Mock API content")
}

func TestApi_InvalidURL(t *testing.T) {
	app := fiber.New()
	app.Get("/api/*", Api)

	req := httptest.NewRequest(http.MethodGet, "/api/invalid-url", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	// fetchSite returns 500 for unparseable/unreachable URLs
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
