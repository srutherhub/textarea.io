package handlers

import (
	"app/components"
	"net/http"
)

func UIForgetArea() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		components.ForgetArea().Render(r.Context(), w)
	}
}

func UIDeleteArea() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		components.DeleteArea().Render(r.Context(), w)
	}
}

func UIClear() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		return
	}
}
