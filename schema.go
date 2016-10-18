package auth

import "database/sql"

//CreateSchema Creates schema for the db
func createSchema(db *sql.DB) error {
	userStmt := `
	CREATE TABLE IF NOT EXISTS user (
		id INT NOT NULL AUTO_INCREMENT,
		username VARCHAR(256),
		password VARCHAR(256),
		PRIMARY KEY(id)
	);
	`
	FireSingleStmt(userStmt, db, nil)

	return nil
}
