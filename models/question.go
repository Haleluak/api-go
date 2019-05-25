package models

import (
	"database/sql"
	u "github.com/Haleluak/kb-backend/config/util"
)

type Question struct {
	ID           int    `json:"question_id"`
	QuestionerId int    `json:"questioner_id"`
	QuestionText string `json:"question_text"`
	CountView    int    `json:"k_views"`
	Slug         string `json:"slug"`
	Answer       []Answer `json:"answer"`
	User         User   `json:"user"`
	UpdatedAt    string `json:"updated_at"`
	CreatedAt    string `json:"created_at"`
}

func (question *Question) Validate() (map[string]interface{}, bool) {
	if len(question.QuestionText) == 0 {
		return u.Message(false, "Question text is required"), false
	}
	return u.Message(true, "success"), true
}

func GetQuetion(dbClient *sql.DB, question_id int) (idUser int, returnFalse bool) {
	stmtOut, err := dbClient.Prepare("SELECT questioner_id FROM questions WHERE question_id  = ?")
	if err != nil {
		return idUser, returnFalse
	}
	defer stmtOut.Close()
	err = stmtOut.QueryRow(question_id).Scan(&idUser)
	if err == sql.ErrNoRows {
		return idUser, returnFalse
	} else if err != nil {
		return idUser, returnFalse
	}
	return idUser, true

}
