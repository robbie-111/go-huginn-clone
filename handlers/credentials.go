package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	credComponents "go-huginn-clone/components/user_credentials"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

func CredentialsIndex(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Credentials")
	creds := models.MockUserCredentials()
	pagination := mockPagination(len(creds), "/user_credentials")
	credComponents.Index(props, creds, pagination, "credential_name", "asc").Render(r.Context(), w)
}

func CredentialsNew(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Create Credential")
	props.LoadAceEditor = true
	cred := models.UserCredential{Mode: "text"}
	credComponents.NewPage(props, cred, nil).Render(r.Context(), w)
}

func CredentialsEdit(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, _ := strconv.Atoi(idStr)
	creds := models.MockUserCredentials()
	var cred models.UserCredential
	for _, c := range creds {
		if c.ID == id {
			cred = c
			break
		}
	}
	props := makeProps(w, r, cred.CredentialName)
	props.LoadAceEditor = true
	credComponents.EditPage(props, cred, nil).Render(r.Context(), w)
}

func CredentialsCreate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Credential was successfully created.")
	redirect(w, r, "/user_credentials")
}

func CredentialsUpdate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Credential was successfully updated.")
	redirect(w, r, "/user_credentials")
}

func CredentialsDestroy(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Credential was successfully deleted.")
	redirect(w, r, "/user_credentials")
}

func CredentialsImport(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Credentials imported successfully.")
	redirect(w, r, "/user_credentials")
}
