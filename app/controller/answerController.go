package controller

import (
	"encoding/json"
	"github.com/Haleluak/kb-backend/app/repository"
	u "github.com/Haleluak/kb-backend/config/util"
	"github.com/Haleluak/kb-backend/models"
	"net/http"
	"time"
)

var CreateAnswer = func(w http.ResponseWriter, r *http.Request) {
	user_id := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	answer := &models.Answer{}
	err := json.NewDecoder(r.Body).Decode(answer) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	currentTime := time.Now()
	answer.AnswererId = int(user_id)
	answer.CreatedAt = currentTime.Format("2006-01-02 15:04:05")
	answer.UpdatedAt = currentTime.Format("2006-01-02 15:04:05")
	resp := repository.CreateAnswer(answer) //Create account
	u.Respond(w, resp)
}
