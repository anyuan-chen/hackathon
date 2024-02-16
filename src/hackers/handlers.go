package hackers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/anyuan-chen/hackathon/.gen/hackathon/public/model"
	. "github.com/anyuan-chen/hackathon/.gen/hackathon/public/table"
	"github.com/anyuan-chen/hackathon/src/auth"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/gorilla/mux"
)

// Skill represents a skill with a name and rating
type InputSkill struct {
	Skill  *string `json:"skill"`
	Rating *int    `json:"rating"`
}

// Person represents a person with name, company, email, phone, and skills
type InputPerson struct {
	Name    *string       `json:"name"`
	Company *string       `json:"company"`
	Email   *string       `json:"email"`
	Phone   *string       `json:"phone"`
	Skills  *[]InputSkill `json:"skills"`
}

func readBody(r *http.Request, w http.ResponseWriter) (interface{}, error) {
	var resp interface{}
	var data *bytes.Buffer
	data.ReadFrom(r.Body)
	err := json.Unmarshal(data.Bytes(), &resp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("malformed json body"))
		return nil, errors.New("")
	}
	return resp, nil
}

func writeBody(r *http.Request, w http.ResponseWriter, resp interface{}) {
	json_resp, err := json.Marshal(resp)
	if err != nil {
		panic("coudln't serialize body")
	}
	w.Write(json_resp)
}

func databaseError(dest Statement, w http.ResponseWriter) {
	w.Write([]byte(dest.DebugSql() + " query failed"))
	w.WriteHeader(http.StatusInternalServerError)
}

func handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
func handleLogin(db *sql.DB) http.HandlerFunc {
	type response struct {
		Access_token string `json:"access_token"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		id, pass, err := auth.GetIdPasswordFromRequest(r)
		if err != nil {
			w.Write([]byte(err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		stmt := Users.SELECT(Users.HashedSecret, Users.Salt).WHERE(Users.ID.EQ(Int32(id)))
		var users []model.Users
		err = stmt.Query(db, &users)
		if err != nil || len(users) != 1 {
			databaseError(stmt, w)
			return
		}
		user := users[0]
		hashed_scrt := auth.VerifyHashedPassword(pass, *user.Salt)
		if *user.HashedSecret != hashed_scrt {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("user hash: ", *user.HashedSecret)
			log.Println("provided hash ", hashed_scrt)

			w.Write([]byte("incorrect pw for user"))
			return
		}
		token, err := auth.GetBearerTokenByUserId(id, db)
		if err != nil {
			w.Write([]byte("no auth provided"))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		/* do work here */
		resp := response{
			Access_token: *token.BearerToken,
		}
		log.Println("handleLogin", *token.BearerToken)
		var json_resp []byte
		json_resp, err = json.Marshal(resp)
		if err != nil {
			panic("bruh")
		}
		w.Write(json_resp)
	}
}

func handleGetAllUsers(db *sql.DB) http.HandlerFunc {
	type request struct {
	}
	type response struct {
		users model.Users
	}
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleGetOneUser(db *sql.DB) http.HandlerFunc {
	type request struct {
		id int32
	}
	type response = model.Users
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("id not integer"))
			return
		}
		var users []model.Users
		stmt := Users.SELECT(Users.AllColumns).WHERE(Users.ID.EQ(Int32(int32(id))))
		err = stmt.Query(db, &users)
		if err != nil {
			databaseError(stmt, w)
			return
		} else if len(users) == 0{
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(stmt.DebugSql()))
			return
		}
		user := users[0]
		log.Println("user from handler", user)
		user.HashedSecret = nil
		user.Salt = nil
		writeBody(r, w, user)
	}
}

func handleUpdateOneUser(db *sql.DB) http.HandlerFunc {
	type request = InputPerson
	type response = model.Users
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleGetAllSkills(db *sql.DB) http.HandlerFunc {
	type request struct {
		min_freq int `json:"min_freq"`
		max_freq int `json:"max_freq"`
	}
	type response = []model.Skills
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleCreateUser(db *sql.DB) http.HandlerFunc {
	type request = InputPerson
	type response = InputPerson
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func handleScanUser(db *sql.DB) http.HandlerFunc {
	type request struct {
		id string
	}
	type response struct {
		qrcode string `json:"qrcode"`
	}
	//hash the time so the QR code generated is different for a certain time, qr code generated will lose effect in x amount of time

	//anti-fraud for funsies + defence, ticketmaster has this!
	//if incorrect user, report as fraud
	//regenerate if wrong time
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

//if have time: implement
