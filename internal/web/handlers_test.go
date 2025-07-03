package web

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWebHandlers(t *testing.T) {
	handlers := NewWebHandlers()
	assert.NotNil(t, handlers)
}

func TestWebHandlers_HandleHome(t *testing.T) {
	handlers := NewWebHandlers()

	// Create static directory and index.html for test
	err := os.MkdirAll("static", 0755)
	require.NoError(t, err)
	defer os.RemoveAll("static")

	err = os.WriteFile("static/index.html", []byte("<!DOCTYPE html><html><body>Test</body></html>"), 0644)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Set up a proper HTTP request for the context
	req := httptest.NewRequest("GET", "/", http.NoBody)
	c.Request = req

	handlers.HandleHome(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "<!DOCTYPE html>")
	assert.Contains(t, w.Body.String(), "Test")
}
