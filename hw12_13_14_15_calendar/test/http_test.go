package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/app"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/config"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/baserepo"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/db/inmemory"
	"github.com/widrik/hw/hw12_13_14_15_calendar/internal/logging"
	httpserver "github.com/widrik/hw/hw12_13_14_15_calendar/internal/server/http"
)

func makeRequest(handler http.Handler, method, path string, body io.Reader) *httptest.ResponseRecorder {
	request, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, request)

	return w
}

func TestHttpEndpoints(t *testing.T) {
	configuration, _ := config.Init("testdata/config/correctData.json")
	_ = logging.Init(configuration.Logging.Level, configuration.Logging.File)
	repo := new(inmemory.Repo)
	calenderApp := app.Calendar{
		Repository: repo,
	}
	handler := httpserver.CreateHandler(&calenderApp)

	startAt := time.Date(2021, 1, 1, 10, 10, 0, 0, time.UTC)
	finishedAt := time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)
	notifyAt := time.Date(2021, 1, 1, 12, 0, 0, 0, time.UTC)

	event := baserepo.Event{
		Title:       "Встреча очень важная",
		Description: "Тестовая очень важная встреча",
		StartedAt:   startAt,
		FinishedAt:  finishedAt,
		NotifyAt:    notifyAt,
		UserID:      1,
	}
	jsonEvent, _ := json.Marshal(event)

	t.Run("main page ok", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		request := makeRequest(handler, "GET", "/", nil)
		response := request.Result()

		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("add wrong method", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		request := makeRequest(handler, "GET", "/add", bytes.NewReader(jsonEvent))
		response := request.Result()

		require.Equal(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("add success", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		request := makeRequest(handler, "POST", "/add", bytes.NewReader(jsonEvent))
		response := request.Result()

		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("get list wrong method", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		request := makeRequest(handler, "POSt", "/getlist", nil)
		response := request.Result()

		require.Equal(t, http.StatusNotFound, response.StatusCode)
	})

	t.Run("get list success", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		request := makeRequest(handler, "GET", "/getlist", nil)
		response := request.Result()

		require.Equal(t, http.StatusOK, response.StatusCode)
	})

	t.Run("delete success", func(t *testing.T) {
		gin.SetMode(gin.ReleaseMode)
		request := makeRequest(handler, "POST", "/add", bytes.NewReader(jsonEvent))

		resEvent := baserepo.Event{}
		_ = json.Unmarshal(request.Body.Bytes(), &resEvent)

		request = makeRequest(handler, "DELETE", "/delete/" + resEvent.ID.String(), bytes.NewReader(jsonEvent))
		response := request.Result()

		require.Equal(t, http.StatusOK, response.StatusCode)
	})
}
