package tequilapi

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCorsHeadersAreAppliedToResponse(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/not-important", nil)
	assert.NoError(t, err)

	respRecorder := httptest.NewRecorder()

	mock := &mockedHTTPHandler{}

	ApplyCors(mock).ServeHTTP(respRecorder, req)

	assert.NotEmpty(t, respRecorder.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, respRecorder.Header().Get("Access-Control-Allow-Methods"))
	assert.True(t, mock.wasCalled)
}

func TestPreflightCorsCheckIsHandled(t *testing.T) {
	req, err := http.NewRequest(http.MethodOptions, "/not-important", nil)
	assert.NoError(t, err)
	req.Header.Add("Origin", "Original site")
	req.Header.Add("Access-Control-Request-Method", "POST")
	req.Header.Add("Access-Control-Request-Headers", "origin, x-requested-with")

	respRecorder := httptest.NewRecorder()

	mock := &mockedHTTPHandler{}

	ApplyCors(mock).ServeHTTP(respRecorder, req)

	assert.NotEmpty(t, respRecorder.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, respRecorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "origin, x-requested-with", respRecorder.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, 0, respRecorder.Body.Len())
	assert.False(t, mock.wasCalled)
}

func TestDeleteCorsPreflightCheckIsHandledCorrectly(t *testing.T) {
	req, err := http.NewRequest(http.MethodOptions, "/not-important", nil)
	assert.NoError(t, err)
	req.Header.Add("Origin", "Original site")
	req.Header.Add("Access-Control-Request-Method", "DELETE")

	respRecorder := httptest.NewRecorder()

	mock := &mockedHTTPHandler{}

	ApplyCors(mock).ServeHTTP(respRecorder, req)

	assert.NotEmpty(t, respRecorder.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, respRecorder.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, 0, respRecorder.Body.Len())
	assert.False(t, mock.wasCalled)

}

func TestCacheControlHeadersAreAddedToResponse(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/not-important", nil)
	assert.NoError(t, err)
	respRecorder := httptest.NewRecorder()

	mock := &mockedHTTPHandler{}

	DisableCaching(mock).ServeHTTP(respRecorder, req)

	assert.Equal(
		t,
		[]string{
			"no-cache",
			"no-store",
			"must-revalidate",
		},
		respRecorder.HeaderMap["Cache-Control"],
	)
	assert.True(t, mock.wasCalled)

}

type mockedHTTPHandler struct {
	wasCalled bool
}

func (mock *mockedHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mock.wasCalled = true
}
