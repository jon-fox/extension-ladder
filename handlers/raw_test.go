package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRaw_ValidURL(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`<html><body><p>Mock raw content</p></body></html>`))
	}))
	defer mockServer.Close()

	app := fiber.New()
	app.Get("/raw/*", Raw)

	req := httptest.NewRequest(http.MethodGet, "/raw/"+mockServer.URL, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.True(t, strings.Contains(string(body), "Mock raw content"))
}

func TestRaw_InvalidURL(t *testing.T) {
	app := fiber.New()
	app.Get("/raw/*", Raw)

	req := httptest.NewRequest(http.MethodGet, "/raw/invalid-url", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	body, _ := io.ReadAll(resp.Body)
	assert.NotEmpty(t, body) // should contain the error message
}
