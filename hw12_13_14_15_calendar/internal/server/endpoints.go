package server

import (
	"net/http"
)

func (a *App) HelloWorld(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("hello - world")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}