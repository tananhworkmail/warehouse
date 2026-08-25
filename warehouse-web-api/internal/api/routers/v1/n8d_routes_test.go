package router_v1

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestN8dRoutesAreRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	Register(engine)

	routes := make(map[string]struct{})
	for _, route := range engine.Routes() {
		routes[route.Method+" "+route.Path] = struct{}{}
	}

	expected := []string{
		"POST /api/v1/auth/login",
		"POST /api/v1/scan/Login",
		"POST /api/v1/scan/login",
		"GET /api/v1/scan/number-scan/inventory",
		"POST /api/v1/scan/number-scan/inventory",
		"GET /api/v1/scan/number-scan/inventory-month",
		"POST /api/v1/scan/number-scan/inventory-month",
		"GET /api/v1/scan/number-scan/rack-list",
		"POST /api/v1/scan/tracking/meterialimport",
		"POST /api/v1/scan/tracking/meterialexport",
		"POST /api/v1/scan/tracking/merterialprint",
		"POST /api/v1/scan/tracking/updatekcbhbybarcode",
		"POST /api/v1/scan/tracking/getkcllscan_rfss",
		"POST /api/v1/scan/tracking/getkcllscan_rfsss",
		"POST /api/v1/scan/tracking/updatekcllscan_rfss",
		"POST /api/v1/scan/tracking/updatekcrkscan_rfss",
		"POST /api/v1/scan/tracking/updatekcllscan_rfsss",
	}

	for _, route := range expected {
		if _, ok := routes[route]; !ok {
			t.Errorf("missing N8d route %s", route)
		}
	}
}

func TestN8dLoginRouteRejectsInvalidPayloadInsteadOfReturningNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	Register(engine)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/scan/Login",
		strings.NewReader(`{}`),
	)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()

	engine.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected invalid login payload to return 400, got %d", response.Code)
	}
}
