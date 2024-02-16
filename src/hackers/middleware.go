package hackers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/anyuan-chen/hackathon/.gen/hackathon/public/model"
	. "github.com/anyuan-chen/hackathon/.gen/hackathon/public/table"
	"github.com/anyuan-chen/hackathon/src/auth"
	. "github.com/go-jet/jet/v2/postgres"
)

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

func getCurrentUser(db *sql.DB, w http.ResponseWriter, r *http.Request) (model.Users, error) {
	token, err := auth.GetBearerToken(r)
	if err != nil {
		print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + " no bearer token"))
	}

	user, err := auth.GetUserFromBearerToken(db, token)
	if err != nil {
		print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + " no user associated with token"))
	}

	statement := SELECT(
		Users.Role,
	).FROM(Users).DISTINCT().WHERE(Users.ID.EQ(Int32(int32(user))))

	var dest []model.Users
	err = statement.Query(db, &dest)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("database failed to execute query"))
		return model.Users{}, errors.New("failed to execute q")
	}
	if len(dest) != 1 {
		return model.Users{}, errors.New("bad # of records")
	}
	return dest[0], nil
}

func adminOnly(db *sql.DB, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := getCurrentUser(db, w, r)
		if err != nil {
			return
		}
		if *(user.Role) != "admin" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("only admins may call this"))
			return
		}
		h.ServeHTTP(w, r)
	})
}

func selfOrAdmin(db *sql.DB, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		id_str := queryParams.Get("id")
		id, err := strconv.Atoi(id_str)

		user, err := getCurrentUser(db, w, r)
		if err != nil {
			return
		}
		if *(user.Role) != "admin" || id != int(user.ID) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("only admins or owner of account may call this"))
			return
		}
		h.ServeHTTP(w, r)
	})
}
