package repository

import (
	"fmt"
	u "github.com/Haleluak/kb-backend/config/util"
	"github.com/Haleluak/kb-backend/models"
)

func CreateQuestion(question *models.Question) (map[string]interface{}) {
	if res, ok := question.Validate(); !ok {
		return res
	}

	dbClient, err := models.DatabaseMysql.ConnectDB()
	defer dbClient.Close()

	stmtIns, err := dbClient.Prepare("INSERT INTO questions ( questioner_id, question_text, slug, created_at, updated_at) VALUES( ?, ?, ?, ?, ?)")
	if err != nil {
		return u.Message(false, "error when insert "+err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	result, err := stmtIns.Exec(question.QuestionerId, question.QuestionText, question.Slug, question.CreatedAt, question.UpdatedAt)
	if err != nil {
		return u.Message(false, "Error Exec when insert user "+err.Error())
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return u.Message(false, "Failed to create account, connection error")
	}
	response := u.Message(true, "Question has been created")
	question.ID = int(lastInsertID)
	response["question"] = question
	return response
}

func FetchQuestion() (map[string]interface{}) {
	dbClient, err := models.DatabaseMysql.ConnectDB()
	defer dbClient.Close()

	stmtOut, err := dbClient.Prepare(`SELECT
											q.question_id , q.questioner_id, q.question_text, q.slug , q.k_views, q.created_at, q.updated_at,
											COALESCE(a.answer_id, 0), COALESCE(a.answerer_id, 0), COALESCE(a.question_id, 0), COALESCE(a.questioner_id, 0),
											COALESCE(a.f_flag, 0), COALESCE(a.answer_text, '') as answer_text , COALESCE(a.created_at, '') as created_at,
											user.user_id, user.name, user.credentials_name, COALESCE(a.user_id, 0), COALESCE(a.name,''), COALESCE(a.credentials_name,'')
											FROM questions as q 
											JOIN user ON user.user_id = questioner_id
											LEFT JOIN (
											  SELECT answers.*, ua.user_id, ua.name, ua.credentials_name FROM answers 
											  JOIN user as ua ON ua.user_id = answers.answerer_id
											  WHERE answers.f_flag = 1 OR answers.created_at in (select max(answers.created_at) FROM answers GROUP BY answers.question_id)
											  ) as a ON q.question_id = a.question_id  ORDER BY q.question_id desc`)

	if err != nil {
		return u.Message(false, "error"+err.Error())
	}
	defer stmtOut.Close()

	rowQuery, err := stmtOut.Query()
	if err != nil {
		return u.Message(false, "Failed to get question"+err.Error())
	}
	result := make([] *models.Question, 0)
	for rowQuery.Next() {
		t := new(models.Question)
		i := models.Answer{}
		//count := 0
		err = rowQuery.Scan(
			&t.ID,
			&t.QuestionerId,
			&t.QuestionText,
			&t.Slug,
			&t.CountView,
			&t.CreatedAt,
			&t.UpdatedAt,
			&i.ID,
			&i.AnswererId,
			&i.QuestionId,
			&i.QuestionerId,
			&i.Flag,
			&i.AnswerText,
			&i.CreatedAt,
			&t.User.ID,
			&t.User.Name,
			&t.User.CredentialsName,
			&i.User.ID,
			&i.User.Name,
			&i.User.CredentialsName,
		)
		if err != nil {
			return u.Message(false, "Failed to get question, connection error"+err.Error())
		}
		//t.CountView = count
		t.Answer = append(t.Answer, i)
		result = append(result, t)
	}
	response := u.Message(true, "success")
	response["data"] = result
	return response
}
func GetByID(nId uint) (map[string]interface{}) {
	question := &models.Question{}
	dbClient, err := models.DatabaseMysql.ConnectDB()
	stmtOut, err := dbClient.Prepare(`SELECT 
											q.question_id , q.questioner_id, q.question_text, q.slug , q.k_views, q.created_at, q.updated_at,
											user.user_id, user.name, user.credentials_name
											FROM questions as q
											JOIN user ON user.user_id = q.questioner_id
											WHERE q.question_id  = ?`)
	if err != nil {
		return u.Message(false, "Error "+err.Error())
	}
	defer dbClient.Close()
	err = stmtOut.QueryRow(nId).Scan(
		&question.ID,
		&question.QuestionerId,
		&question.QuestionText,
		&question.Slug,
		&question.CountView,
		&question.CreatedAt,
		&question.UpdatedAt,
		&question.User.ID,
		&question.User.Name,
		&question.User.CredentialsName,
		)
	if err != nil {
		return u.Message(false, "Failed to get question, connection error "+err.Error())
	}
	question.Answer, err = models.GetAnswerByQuestionId(dbClient, nId)
	if err!= nil{
		fmt.Print("error" + err.Error())
	}
	response := u.Message(true, "success")
	response["data"] = question
	return response
}
