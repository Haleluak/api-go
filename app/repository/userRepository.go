package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	u "github.com/Haleluak/kb-backend/config/util"
	"github.com/Haleluak/kb-backend/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func Login(email, password string) (map[string]interface{}) {
	user := &models.User{}

	dbClient, err := models.DatabaseMysql.ConnectDB()
	defer dbClient.Close()

	stmtOut, err := dbClient.Prepare("SELECT user_id, name, profile_pic_url, profile_pic_thumbnail_url, exp_from, exp_to, credentials_name, created_at, updated_at, email, password FROM user WHERE email = ?")
	if err != nil {
		return u.Message(false, "PrePare Error "+err.Error())
	}
	defer stmtOut.Close()

	err = stmtOut.QueryRow(email).Scan(&user.ID, &user.Name, &user.ProfilePicUrl, &user.ProfilePicThumbnailUrl, &user.ExpFrom, &user.ExpTo, &user.CredentialsName, &user.CreatedAt, &user.UpdatedAt, &user.Email, &user.Password)
	if err == sql.ErrNoRows {
		return u.Message(false, "Email address not found")
	} else if err != nil {
		return u.Message(false, "Connection error. Please retry"+err.Error())
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	//Worked! Logged In
	user.Password = ""

	//Create JWT token
	tk := &models.Token{UserId: uint(user.ID)}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("thisIsTheJwtPassword"))
	user.Token = tokenString //Store the token in the response
	resp := u.Message(true, "Logged In")
	resp["account"] = user
	return resp
}

func Create(user *models.User) (map[string]interface{}) {
	if res, ok := user.Validate(); !ok {
		return res
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	dbClient, err := models.DatabaseMysql.ConnectDB()
	defer dbClient.Close()

	stmtIns, err := dbClient.Prepare("INSERT INTO user ( name, email, password, created_at, updated_at) VALUES( ?, ?, ?, ?, ?)")
	if err != nil {
		return u.Message(false, "error when inert s"+err.Error())
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	result, err := stmtIns.Exec(user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return u.Message(false, "Error Exec when insert user "+err.Error())
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return u.Message(false, "Failed to create account, connection error")
	}

	//Create new JWT token for the newly registered account
	tk := &models.Token{UserId: uint(lastInsertID)}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("thisIsTheJwtPassword"))
	user.ID = int(lastInsertID)
	user.Token = tokenString
	user.Password = "" //delete password
	response := u.Message(true, "Account has been created")
	response["account"] = user
	return response
}

func LoginFacebook(access_token string) (map[string]interface{}) {
	response, err := http.Get("https://graph.facebook.com/me?fields=email,name,id&access_token=" + access_token)
	if err != nil {
		return u.Message(false, "Failed to login facebook account, connection error")
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return u.Message(false, "facebook-auth ioutil.ReadAll... Fail")
	}

	//var user models.User
	data := make(map[string]interface{})
	err = json.Unmarshal(contents, &data)
	if err != nil {
		log.Println("facebook-auth json.Unmarshal... Fail", err)
		return u.Message(false, "facebook-auth json.Unmarshal... Fail")
	}

	//parse data from response facebook
	facebook_id := data["id"].(string)
	email := data["email"] .(string)
	fbusername, _ := data["name"] .(string)
	avatar := "https://graph.facebook.com/" + facebook_id + "/picture?type=large"
	currentTime := time.Now()
	created_at := currentTime.Format("2006-01-02 15:04:05")
	updated_at := currentTime.Format("2006-01-02 15:04:05")
	//create account facebook
	dbClient, err := models.DatabaseMysql.ConnectDB()
	defer dbClient.Close()

	//check login facebook already
	var lastID int
	user, ok := models.CheckLoginFb(facebook_id)
	if !ok{
		stmtIns, err := dbClient.Prepare("INSERT INTO user ( name, email, facebook_id, profile_pic_url, created_at, updated_at) VALUES( ?, ?, ?, ?, ?, ?)")
		if err != nil {
			return u.Message(false, "error when insert "+err.Error())
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

		result, err := stmtIns.Exec(fbusername, email, facebook_id, avatar, created_at, updated_at)
		if err != nil {
			return u.Message(false, "Error Exec when insert user "+err.Error())
		}
		lastInsertID, err := result.LastInsertId()
		if err != nil {
			return u.Message(false, "Failed to create account, connection error")
		}
		lastID = int(lastInsertID)
		user.ID = lastID
		user.Name = fbusername
		user.Email = email
		user.CreatedAt = created_at
		user.UpdatedAt = updated_at
	}
	//Create new JWT token for the newly registered account
	lastID = user.ID
	tk := &models.Token{UserId: uint(lastID)}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte("thisIsTheJwtPassword"))
	user.Token = tokenString

	results := u.Message(true, "Logged In")
	results["account"] = user
	return results
}

func UpdateUser(userDetails models.User) bool {
	dbClient, err := models.DatabaseMysql.ConnectDB()
	defer dbClient.Close()
	stmt, err := dbClient.Prepare("UPDATE user SET credentials_name=?,exp_from=?,exp_to=? WHERE user_id=?")
	if err != nil {
		fmt.Println(err)
		return false
	}
	_, queryError := stmt.Exec(userDetails.CredentialsName, userDetails.ExpFrom, userDetails.ExpTo, userDetails.ID)
	if queryError != nil {
		fmt.Println(queryError)
		return false
	}
	return true
}
