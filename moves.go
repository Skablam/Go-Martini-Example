// moves
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

type move struct {
	Name       string
	Category   string
	Tags       string
	Video_link string
	Username   string
	Id         int
}

type get_move_list_resp struct {
	Moves []move
}

type get_id_resp struct {
	Id int
}

func GetMoveList(res http.ResponseWriter, req *http.Request, db *sql.DB) {

	username := req.FormValue("username")

	var move_list []move

	rows, statementError := db.Query("SELECT trim(name), trim(category), trim(tags), trim(video_link), id FROM moves WHERE username = $1", username)

	if statementError != nil {
		panic(statementError)
	}

	defer rows.Close()

	move_list = make([]move, 0)

	for rows.Next() {

		var name, category, tags, video_link string
		var id int

		err := rows.Scan(&name, &category, &tags, &video_link, &id)

		if err != nil {
			fmt.Printf("rows.Scan error: %v\n", err)
		}

		current_move := move{name, category, tags, video_link, username, id}

		move_list = append(move_list, current_move)

	}

	resp := new(get_move_list_resp)
	resp.Moves = move_list

	//convert resp struct to json
	jsonResp, encError := json.Marshal(resp)
	if encError != nil {
		fmt.Println(encError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	res.Write(jsonResp)
}

func AddMove(res http.ResponseWriter, req *http.Request, db *sql.DB) {

	decoder := json.NewDecoder(req.Body)

	newmove := move{}

	err := decoder.Decode(&newmove)
	if err != nil {
		panic(err.Error())
	}

	var insert_statement string = "insert into moves (name, category, tags, video_link, username) values ($1, $2, $3, $4, $5)"

	stmt, err := db.Prepare(insert_statement)

	if err != nil {
		panic(err.Error())
	}

	queryResult, err := stmt.Exec(newmove.Name, newmove.Category, newmove.Tags, newmove.Video_link, newmove.Username)

	if err != nil || queryResult == nil {
		panic(err.Error())
	}

	stmt.Close()

	res.Write([]byte("Move added"))
}

func GetMoveID(res http.ResponseWriter, req *http.Request, db *sql.DB) {

	username := req.FormValue("username")
	movename := req.FormValue("movename")

	select_statement := "SELECT id FROM moves WHERE username = $1 and name = $2"

	stmt, err := db.Prepare(select_statement)

	if err != nil {
		panic(err.Error())
	}

	queryResult, err := stmt.Query(username, movename)

	if err != nil || queryResult == nil {
		panic(err.Error())
	}

	queryResult.Next()

	var id int

	queryResult.Scan(&id)

	if err != nil {
		fmt.Printf("rows.Scan error: %v\n", err)
	}

	resp := new(get_id_resp)
	resp.Id = id

	//convert resp struct to json
	jsonResp, encError := json.Marshal(resp)
	if encError != nil {
		fmt.Println(encError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	res.Write(jsonResp)

}

func DeleteMove(res http.ResponseWriter, req *http.Request, db *sql.DB, params martini.Params) {

	move_id := params["moveid"]

	var delete_statement string = "delete from moves where id = $1"

	stmt, err := db.Prepare(delete_statement)

	if err != nil {
		panic(err.Error())
	}

	queryResult, err := stmt.Exec(move_id)

	if err != nil || queryResult == nil {
		panic(err.Error())
	}

	stmt.Close()

	res.Write([]byte("Move deleted"))

}

func GetMove(res http.ResponseWriter, req *http.Request, db *sql.DB, params martini.Params) {

	move_id, err := strconv.Atoi(params["moveid"])

	select_statement := "SELECT trim(name), trim(category), trim(tags), trim(video_link), trim(username) FROM moves WHERE id = $1"

	stmt, err := db.Prepare(select_statement)

	if err != nil {
		panic(err.Error())
	}

	queryResult, err := stmt.Query(move_id)

	if err != nil || queryResult == nil {
		panic(err.Error())
	}

	queryResult.Next()

	var name, category, tags, video_link, username string

	queryResult.Scan(&name, &category, &tags, &video_link, &username)

	//create move struct
	resp := move{name, category, tags, video_link, username, move_id}

	//convert resp struct to json
	jsonResp, encError := json.Marshal(resp)
	if encError != nil {
		fmt.Println(encError)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	res.Write(jsonResp)

}

func AmendMove(res http.ResponseWriter, req *http.Request, db *sql.DB, params martini.Params) {

	move_id := params["moveid"]

	decoder := json.NewDecoder(req.Body)

	current_move := move{}

	err := decoder.Decode(&current_move)
	if err != nil {
		panic(err.Error())
	}

	var update_statement string = "update moves set name = $1, category = $2, tags = $3, video_link = $4 where id = $5"

	stmt, err := db.Prepare(update_statement)

	if err != nil {
		panic(err.Error())
	}

	queryResult, err := stmt.Exec(current_move.Name, current_move.Category, current_move.Tags, current_move.Video_link, move_id)

	if err != nil || queryResult == nil {
		panic(err.Error())
	}

	stmt.Close()

	res.Write([]byte("Move updated"))

}
