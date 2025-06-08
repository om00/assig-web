package handler

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/om00/assig-web/models"
	"github.com/om00/assig-web/psqldb"
)

type App struct {
	Db *psqldb.DbIns
}

func (app *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		username := r.FormValue("username")
		password := r.FormValue("password")

		name := "om"
		dbPassword := "om12345"

		// Check if the user exists and password is correct
		// var dbPassword string
		// err := db.QueryRow("SELECT password FROM users WHERE username = $1", username).Scan(&dbPassword)
		// if err != nil {
		// 	http.Error(w, "User not found", http.StatusUnauthorized)
		// 	return
		// }

		if username == name && password == dbPassword {
			// If credentials are correct
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		} else {
			// Incorrect password
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		}
	}
}

func ShowDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("app/webpage/dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func (app *App) ShowAllUser(w http.ResponseWriter, r *http.Request) {
	var request models.UserRequest
	request.Name = r.URL.Query().Get("username")

	request.Email = r.URL.Query().Get("email")
	request.Status = r.URL.Query().Get("status")
	request.BlockReason = r.URL.Query().Get("blockReason")

	request.PhoneStr = r.URL.Query().Get("phone")
	if request.PhoneStr != "" {
		request.Phone = strings.Split(request.PhoneStr, ",")
	}

	err := request.HandleIntFields()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	users, err := app.Db.GetAllUsers(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reasons := map[int]string{
		1:  "Fraudulent Activity",
		2:  "Used in Multiple Locations",
		3:  "Suspected Criminal Activity",
		99: "Other", // Added Other as an option with ID 0
	}

	userData := models.UserData{
		Users:   users,
		Filter:  request,
		Reasons: reasons,
	}
	tmpl, err := template.ParseFiles("app/webpage/user.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, userData)
}

func (app *App) BlockUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var request models.UserRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if request.PhoneStr != "" {
		request.Phone = strings.Split(request.PhoneStr, ",")
	}

	err = request.HandleIntFields()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = app.Db.BlockUser(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "success",
		"message": "User information updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (app *App) UnblockUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}

	var request models.UserRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = request.HandleIntFields()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	err = app.Db.UnblockUser(request)
	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"status":  "success",
		"message": "User information updated successfully",
	}
	json.NewEncoder(w).Encode(response)
}

func (app *App) Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("app/webpage/adminLogin.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, "")
}
