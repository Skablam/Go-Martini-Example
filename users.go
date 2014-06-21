// users
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"net/http"
)

type user struct {
	Id       string
	Password string
	Email    string
}

type get_user_resp struct {
	User   user
	Result string
}

func RegisterUser(res http.ResponseWriter, req *http.Request, db *sql.DB) {

	newuser := user{}

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&newuser)
	if err != nil {
		panic(err.Error())
	}

	var insert_statement string = "insert into user_accounts (username, password_hash, email) values ($1, $2, $3)"

	stmt, err := db.Prepare(insert_statement)

	if err != nil {
		panic(err.Error())
	}

	queryResult, err := stmt.Exec(newuser.Id, newuser.Password, newuser.Email)

	if err != nil || queryResult == nil {
		panic(err.Error())
	}

	stmt.Close()
}

func GetUser(res http.ResponseWriter, req *http.Request, parms martini.Params, db *sql.DB) {

	username := parms["username"]

	row := db.QueryRow("SELECT trim(email), password_hash FROM user_accounts WHERE username = $1", username)

	var result, email, password_hash string

	var user_info user

	err := row.Scan(&email, &password_hash)

	if err != nil {
		fmt.Printf("rows.Scan error: %v\n", err)
		result = "no user found"
	} else {
		user_info = user{username, password_hash, email}
		result = "user found"
	}

	resp := new(get_user_resp)
	resp.User = user_info
	resp.Result = result

	//convert resp struct to json
	jsonResp, encError := json.Marshal(resp)
	if encError != nil {
		fmt.Println(encError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	res.Write(jsonResp)

}
