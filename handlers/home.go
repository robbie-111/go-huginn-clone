package handlers

import (
	"net/http"

	"go-huginn-clone/components/home"
)

func HomeIndex(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "")
	home.Index(props).Render(r.Context(), w)
}
