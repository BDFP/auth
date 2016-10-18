package auth

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/context"
	"github.com/julienschmidt/httprouter"
	"encoding/json"
	"time"
)


//UserEnv Environment required for user
type AuthEnv struct {
	DB *sql.DB
}


//Secure Secure middleware
func Secure(protectedPage httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(res http.ResponseWriter, req *http.Request,
	_ httprouter.Params) {
		//get authorization header
		header := req.Header.Get("Authorization")
		if header == "" {
			WriteErrorResponse("No token present", res)
			return
		}
		log.Println("Header is " + header)

		//parse jwt token
		token, err := jwt.ParseWithClaims(header, &TokenClaims{},
			func(token *jwt.Token) (interface{}, error) {
				//check signing method
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method %v", token.Header["alg"])
				}
				return []byte("secret"), nil
			})
		if err != nil {
			log.Println(err.Error())
		}

		if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
			context.Set(req, "Claims", claims)
		} else {
			http.NotFound(res, req)
			return
		}

		protectedPage(res, req, nil)
	})
}

func startHTTPServer(db *sql.DB) {
	uEnv := &AuthEnv{
		DB: db,
	}

	router := httprouter.New()

	//should test by writing request format and automate
	router.POST("/register", uEnv.register)
	router.POST("/login", uEnv.login)
	//sth like nodemon for go
	log.Fatal(http.ListenAndServe(":8484", router))
}

func Setup(db *sql.DB) {
	log.Println("Starting Auth.. Please wait")

	err := db.Ping()
	if err != nil {
		panic(err.Error())
	}
	log.Println("Database connection opened")

	log.Println("Creating schema")
	createSchema(db)
	log.Println("Schema Created")

	log.Println("Setting up HTTP Server")
	startHTTPServer(db)
}


//Register Signup Api
func (e *AuthEnv) register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	var user User

	//decode json request
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		WriteErrorResponse(err.Error(), w)
		return
	}

	err = storeUser(e.DB, &user)
	if err != nil {
		WriteErrorResponse(err.Error(), w)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(HTTPResponse{
		Message: "User Registered successfulyy",
	})
	return
}

//Login Login API
func (e *AuthEnv) login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		WriteErrorResponse(err.Error(), w)
		return
	}

	user.Password = HashPassword(user.Password)
	if err := getUser(&user, e.DB); err != nil {
		WriteErrorResponse(err.Error(), w)
		return
	}

	// Expires the token and cookie in 24 hours
	expireToken := time.Now().Add(time.Hour * 24).Unix()

	// We'll manually assign the claims but in production you'd insert values from a database
	claims := TokenClaims{
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: expireToken,
			Issuer:    "shakdwipeea.com",
		},
	}

	// Create the token using your claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Signs the token with a secret.
	signedToken, _ := token.SignedString([]byte("secret"))

	json.NewEncoder(w).Encode(UserLoginResponse{
		Token: signedToken,
		User:  user,
	})
}
