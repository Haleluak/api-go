package repository

import (
	u "github.com/Haleluak/kb-backend/config/util"
	"github.com/Haleluak/kb-backend/models"
)

func CreateAnswer(answer *models.Answer) (map[string]interface{}) {
	if res, ok := answer.Validate(); !ok {
		return res
	}
	dbClient, err := models.DatabaseMysql.ConnectDB()
	defer dbClient.Close()

	idUser, ok := models.GetQuetion(dbClient, answer.QuestionId)
	if !ok {
		return u.Message(false, "Question not exit")
	}
	answer.QuestionerId = idUser
	stmtIns, err := dbClient.Prepare("INSERT INTO answers ( answerer_id, question_id, questioner_id, answer_text, created_at, updated_at) VALUES( ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return u.Message(false, "error when insert "+ err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	result, err := stmtIns.Exec(answer.AnswererId, answer.QuestionId, answer.QuestionerId, answer.AnswerText,answer.CreatedAt, answer.UpdatedAt)
	if err != nil {
		return u.Message(false, "Error Exec when insert user "+ err.Error())
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return u.Message(false, "Failed to create account, connection error")
	}
	response := u.Message(true, "Question has been created")
	answer.ID = int(lastInsertID)
	response["answer"] = answer
	return response
}
