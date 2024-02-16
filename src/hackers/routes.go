package hackers

import (
	"database/sql"

	"github.com/gorilla/mux"
)

func AddRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/health", handleHealth())
	r.Handle("/users", adminOnly(db, handleGetAllUsers(db)))
	r.Handle("/users/{id}", selfOrAdmin(db, handleGetOneUser(db)))
}
