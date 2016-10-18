package auth

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"crypto/sha1"
	"io"
	"fmt"
)

//HTTPResponse Format for error response
type HTTPResponse struct {
	Message string `json:"message"`
}

//WriteErrorResponse utility to return error response
func WriteErrorResponse(reason string, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusBadRequest)
	res := &HTTPResponse{
		Message: reason,
	}

	return json.NewEncoder(w).Encode(&res)
}

//HashPassword encrypte the password
func HashPassword(pass string) string {
	h := sha1.New()
	io.WriteString(h, pass)
	return fmt.Sprintf("%x", h.Sum(nil))
}

//FireSingleStmt executes a statement
func FireSingleStmt(stmt string, db *sql.DB, done chan bool) {
	_, err := db.Exec(stmt)
	if err != nil {
		log.Fatalf("Query %s \n error %s \n", stmt, err.Error())
	}

	if done != nil {
		done <- true
	}
}
