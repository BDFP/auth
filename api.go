package auth

import (
	"database/sql"

	"github.com/dgrijalva/jwt-go"
)

type DbConfig struct {
	DriverName string
	DataSourceName string
}

//User Representation of a user
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//TokenClaims details present in token
type TokenClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

//UserLoginResponse Response structure on login
type UserLoginResponse struct {
	id    string `json:"id"`
	Token string `json:"token"`
	User  User   `json:"user"`
}
