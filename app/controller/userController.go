package controller

import (
	"encoding/json"
	"github.com/Haleluak/kb-backend/app/repository"
	u "github.com/Haleluak/kb-backend/config/util"
	"github.com/Haleluak/kb-backend/models"
	"io/ioutil"
	"net/http"
	"time"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	currentTime := time.Now()
	account.CreatedAt = currentTime.Format("2006-01-02 15:04:05")
	account.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")
	resp := repository.Create(account) //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.User{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	//resp := models.Login(account.Email, account.Password)
	resp := repository.Login(account.Email, account.Password)
	u.Respond(w, resp)
}

var LoginFacebook = func(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	value, ok := data["access_token"]
	if !ok {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := repository.LoginFacebook(value.(string))
	u.Respond(w, resp)
}

var UpdateUser = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	account := models.User{}
	err := json.NewDecoder(r.Body).Decode(&account) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	account.ID = int(user)
	isUpdated := repository.UpdateUser(account)
	if isUpdated {
		u.Respond(w, u.Message(true, "update success "))
	} else {
		u.Respond(w, u.Message(false, "error"))
	}
}
