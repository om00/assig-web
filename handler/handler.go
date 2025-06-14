package handler

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/om00/assig-web/models"
	"github.com/om00/assig-web/psqldb"
	"golang.org/x/crypto/bcrypt"
)

type App struct {
	Db *psqldb.DbIns
}

func (app *App) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		var loginData models.AdminCred
		res := make(map[string]string)
		w.Header().Set("Content-type", "application/json")

		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			res["message"] = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
			return
		}

		admin, err := app.Db.GetAdmin(loginData)
		if err != nil {
			res["message"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
		}

		if err != nil && admin.Id == 0 {
			res["message"] = "Wrong Username .Please Enter valid username"
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(res)
		}

		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginData.Password))
		if err != nil {
			res["message"] = "Wrong Password. Please Enter a Valid password"
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(res)
		}

		userData, tmpl, err := app.ShowUserPage(models.UserRequest{})
		if err != nil {
			res["message"] = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res)
		}
		tmpl.Execute(w, userData)
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
