package auth

import (
	"database/sql"
	"errors"
	"log"

)

//StoreUser save new user
func storeUser(db *sql.DB, user *User) error {
	var id int
	err := db.QueryRow("SELECT id FROM user WHERE username = ?",
		user.Username).Scan(&id)
	if err == sql.ErrNoRows {
		user.Password = HashPassword(user.Password)
		log.Printf("user ins %x", user.Password)
		row, err := db.Exec("INSERT INTO user (username, password) VALUES (?,?)",
			user.Username, user.Password)
		if err != nil {
			return err
		}

		user.ID, err = row.LastInsertId()
		if err != nil {
			return err
		}

		return nil
	}
	return errors.New("username already registered")
}

//GetUser get user from db
func getUser(user *User, db *sql.DB) error {
	var passDb string
	return db.QueryRow(`SELECT id, username, password FROM user
		 WHERE username = ? AND password = ?`, user.Username, user.Password).Scan(
		&user.ID, &user.Username, &passDb)
}
