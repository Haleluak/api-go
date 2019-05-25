package controller

import (
	"encoding/json"
	"github.com/Haleluak/kb-backend/app/repository"
	u "github.com/Haleluak/kb-backend/config/util"
	"github.com/Haleluak/kb-backend/models"
	"github.com/gorilla/mux"
	"github.com/gosimple/slug"
	"net/http"
	"strconv"
	"time"
)

var CreateQuestion = func(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	question := &models.Question{}
	err := json.NewDecoder(r.Body).Decode(question) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	currentTime := time.Now()
	question.Slug = slug.Make(question.QuestionText)
	question.QuestionerId = int(user_id)
	question.CreatedAt = currentTime.Format("2006-01-02 15:04:05")
	question.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")
	resp := repository.CreateQuestion(question) //Create account
	u.Respond(w, resp)
}

var GetQuestions = func(w http.ResponseWriter, r *http.Request) {
	resp := repository.FetchQuestion() //Create account
	u.Respond(w, resp)
}

var GetQuestion = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	nId, err := strconv.ParseUint(string(params["id"]), 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "ID is not valid"))
		return
	}
	resp := repository.GetByID(uint(nId)) //Create account
	u.Respond(w, resp)
}
