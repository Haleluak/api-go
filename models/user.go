package models

import (
	"database/sql"
	u "github.com/Haleluak/kb-backend/config/util"
	"github.com/dgrijalva/jwt-go"
	"log"
	"strings"
)

/*
JWT claims struct
*/
type Token struct {
	UserId uint
	jwt.StandardClaims
}

//a struct to rep user account
type User struct {
	ID                     int     `json:"user_id"`
	FacebookId             *string `json:"facebook_id"`
	Name                   string  `json:"name"`
	ProfilePicUrl          *string  `json:"profile_pic_url"`
	ProfilePicThumbnailUrl *string `json:"profile_pic_thumbnail_url"`
	CredentialsName        *string `json:"credentials_name"`
	Email                  string  `json:"email"`
	Password               string  `json:"password"`
	ExpFrom                *string `json:"exp_from"`
	ExpTo                  *string `json:"exp_to"`
	UpdatedAt              string  `json:"updated_at"`
	CreatedAt              string  `json:"created_at"`
	Token                  string  `json:"token";sql:"-"`
}

func (user *User) Validate() (map[string]interface{}, bool) {
	if ! strings.Contains(user.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(user.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	//Email must be unique
	temp := &User{}

	dbClient, err := DatabaseMysql.ConnectDB()
	defer dbClient.Close()

	stmtOut, err := dbClient.Prepare("SELECT email FROM user WHERE email = ?")
	if err != nil {
		log.Println("db.go Func LoginUser PrePare Error:" + err.Error()) // proper error handling instead of panic in your app
		return u.Message(false, "PrePare Error "+err.Error()), false
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(user.Email).Scan(&temp.Email)
	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement passed"), true
}

func CheckLoginFb(facebook_id string) (User, bool) {
	dbClient, err := DatabaseMysql.ConnectDB()
	if err != nil {
		return User{}, false
	}
	defer dbClient.Close()

	stmtOut, err := dbClient.Prepare("SELECT user_id, name, email, facebook_id, created_at, updated_at FROM user WHERE facebook_id  = ?")
	if err != nil {
		return User{}, false
	}
	defer stmtOut.Close()

	var temp User
	err = stmtOut.QueryRow(facebook_id).Scan(&temp.ID, &temp.Name, &temp.Email, &temp.FacebookId, &temp.CreatedAt, &temp.UpdatedAt)
	if err == sql.ErrNoRows {
		return User{}, false
	} else if err != nil {
		return User{}, false
	}
	return temp, true
}
