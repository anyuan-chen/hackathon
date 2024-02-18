package hackers

import (
	"database/sql"
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

func getCurrentUser(db *sql.DB, w http.ResponseWriter, r *http.Request) (model.Users, error) {
	token, err := auth.GetBearerToken(r)
	if err != nil {
		print(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + " no bearer token"))
		return model.Users{}, err
	}

	user, err := auth.GetUserFromBearerToken(db, token)
	log.Println("token: ", token)
	if err != nil {
		print(err.Error())
		w.Write([]byte(err.Error() + " no user associated with token"))
		return model.Users{}, err
	}

	statement := SELECT(
		Users.AllColumns,
	).FROM(Users).WHERE(Users.ID.EQ(Int32(int32(user))))

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
	// log.Println(statement.DebugSql(), dest[0].ID)
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
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		log.Println("id: ", id)
		if err != nil {
			return
		}
		user, err := getCurrentUser(db, w, r)
		log.Println("current user: ", id, user.ID)
		if err != nil {
			return
		}
		if *(user.Role) != "admin" && id != int(user.ID) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("only admins or owner of account may call this"))
			return
		}
		h.ServeHTTP(w, r)
	})
}
