package handlers

import (
	"encoding/json"
	"net/http"

	authComponents "go-huginn-clone/components/auth"
	"go-huginn-clone/middleware"
	"go-huginn-clone/models"
)

func AuthLoginPage(w http.ResponseWriter, r *http.Request) {
	// If already logged in, redirect to home
	if middleware.GetCurrentUser(r) != nil {
		redirect(w, r, "/")
		return
	}
	props := makeProps(w, r, "Log in")
	authComponents.LoginPage(props, nil).Render(r.Context(), w)
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	login := r.FormValue("user[login]")
	password := r.FormValue("user[password]")

	// Mock auth: any non-empty credentials work
	if login != "" && password != "" {
		middleware.LoginUser(w, r)
		middleware.SetFlash(w, r, "notice", "Signed in successfully.")
		redirect(w, r, "/")
		return
	}

	props := makeProps(w, r, "Log in")
	authComponents.LoginPage(props, []string{"Invalid login or password."}).Render(r.Context(), w)
}

func AuthLogout(w http.ResponseWriter, r *http.Request) {
	middleware.LogoutUser(w, r)
	middleware.SetFlash(w, r, "notice", "Signed out successfully.")
	redirect(w, r, "/")
}

func AuthRegisterPage(w http.ResponseWriter, r *http.Request) {
	props := makeProps(w, r, "Sign up")
	authComponents.RegisterPage(props, nil).Render(r.Context(), w)
}

func AuthRegister(w http.ResponseWriter, r *http.Request) {
	middleware.LoginUser(w, r)
	middleware.SetFlash(w, r, "notice", "Welcome! You have signed up successfully.")
	redirect(w, r, "/")
}

func AuthAccountEdit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetCurrentUser(r)
	if user == nil {
		redirect(w, r, "/users/sign_in")
		return
	}
	props := makeProps(w, r, "Account")
	authComponents.AccountEditPage(props, *user, nil).Render(r.Context(), w)
}

func AuthAccountUpdate(w http.ResponseWriter, r *http.Request) {
	middleware.SetFlash(w, r, "notice", "Your account has been updated successfully.")
	redirect(w, r, "/users/edit")
}

func WorkerStatus(w http.ResponseWriter, r *http.Request) {
	status := models.WorkerStatus{
		Pending:        2,
		AwaitingRetry:  1,
		RecentFailures: 0,
		EventsSince:    5,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
