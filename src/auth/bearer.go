package auth

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/anyuan-chen/hackathon/.gen/hackathon/public/model"
	. "github.com/anyuan-chen/hackathon/.gen/hackathon/public/table"
	. "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/google/uuid"
)

func GetBearerToken(r *http.Request) (string, error) {
	auth := r.Header.Get("Authorization")
	//format should be Bearer {token}, so need to discard Bearer
	auth_array := strings.Split(auth, " ")
	log.Println("auth array", len(auth_array), "auth string:", auth)
	if len(auth_array) < 2 {
		return "", errors.New("out of bounds access error")
	}
	return auth_array[1], nil
}

/*
gets the user's id from their token
*/
func GetUserFromBearerToken(db *sql.DB, token string) (int32, error) {
	statement := SELECT(Tokens.ID).FROM(Tokens).WHERE(Tokens.BearerToken.EQ(String(token)))
	log.Println("token: ", token)
	var bearer []model.Tokens
	statement.Query(db, &bearer)
	if len(bearer) != 1 || bearer[0].ID == 0 {
		return 0, errors.New("no bearer token")
	}
	return bearer[0].ID, nil
}

func GetBearerTokenByUserId(id int32, db *sql.DB) (model.Tokens, error) {
	token_exists_statement := SELECT(Tokens.ID).DISTINCT().FROM(Tokens).WHERE(Tokens.ID.EQ(Int32(id)))
	var currentToken []model.Tokens
	var bearer model.Tokens
	err := token_exists_statement.Query(db, &currentToken)
	log.Println("bearer token geeneration")
	if err == qrm.ErrNoRows || len(currentToken) == 0 {
		log.Println("starting generation")
		bearer, err = generateBearerToken(id, db)
	} else {
		log.Println("starting refresh")
		bearer, err = refreshBearerToken(id, db)
	}
	if err != nil {
		return bearer, err
	}
	return bearer, nil
}

/*
generates the bearer token for a user
*/
func generateBearerToken(id int32, db *sql.DB) (model.Tokens, error) {
	b_token := uuid.New()
	// 8 days and 8 hours in the future
	time := time.Now().Unix() + 720000
	insert_statement := Tokens.INSERT(Tokens.BearerToken, Tokens.ID, Tokens.ExpiryTime).VALUES(b_token, id, time).RETURNING(Tokens.AllColumns)
	var bearer []model.Tokens
	insert_statement.Query(db, &bearer)
	return bearer[0], nil
}

/*
refreshes the bearer token for a user
*/
func refreshBearerToken(id int32, db *sql.DB) (model.Tokens, error) {
	b_token := uuid.New()
	// 8 days and 8 hours in the future
	time := time.Now().Unix() + 720000
	update_statement := Tokens.UPDATE(Tokens.BearerToken, Tokens.ExpiryTime).SET(b_token, time).WHERE(Tokens.ID.EQ(Int32(id))).RETURNING(Tokens.AllColumns)
	var bearer []model.Tokens
	update_statement.Query(db, &bearer)
	if len(bearer) != 1 {
		log.Println("query returned nothing: ", update_statement.DebugSql())
	}
	return bearer[0], nil
}

// new table:
// bearer_token
// user_id

// make account:
// create a user endpoint
// users supplies a password
// password is hashed with a secret
// hashed string is stored in db

// login:
// user supplies id + pw
// pw gets hashed, compared to hashed pw in db
// return random string as bearer token (expires in 3600)

// every request:
//

// server has a centralized hashing token
