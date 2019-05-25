package models

import (
	"database/sql"
	u "github.com/Haleluak/kb-backend/config/util"
)

type Answer struct {
	ID           int    `json:"answer_id"`
	QuestionId   int    `json:"question_id"`
	QuestionerId int    `json:"questioner_id"`
	AnswererId   int    `json:"answerer_id"`
	Flag         int    `json:"f_flag"`
	User         User   `json:"user"`
	AnswerText   string `json:"answer_text"`
	UpdatedAt    string `json:"updated_at"`
	CreatedAt    string `json:"created_at"`
}

func (answer *Answer) Validate() (map[string]interface{}, bool) {
	if answer.AnswerText == "" {
		return u.Message(false, "Answer text is required"), false
	}

	if answer.QuestionId == 0 {
		return u.Message(false, "question id is required"), false
	}
	return u.Message(true, "success"), true
}

func GetAnswerByQuestionId(dbClient *sql.DB, question_id uint) ([]Answer, error) {
	stmtOut, err := dbClient.Prepare(`SELECT 
											a.answer_id, a.question_id, a.questioner_id, a.answerer_id, a.f_flag, a.answer_text, a.created_at, a.updated_at,
											user.user_id, user.name, user.credentials_name
											FROM answers as a
											JOIN user ON user.user_id = a.answerer_id
											WHERE a.question_id  = ?`)
	if err != nil {
		return nil, err
	}
	defer stmtOut.Close()

	rowQuery, err := stmtOut.Query(question_id)
	if err != nil {
		return nil, err
	}
	result := make([] Answer, 0)
	for rowQuery.Next() {
		a := Answer{}
		err = stmtOut.QueryRow(question_id).Scan(
			&a.ID,
			&a.QuestionId,
			&a.QuestionerId,
			&a.AnswererId,
			&a.Flag,
			&a.AnswerText,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.User.ID,
			&a.User.Name,
			&a.User.CredentialsName,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil

}
