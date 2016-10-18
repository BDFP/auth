# Auth
This handles authentication for your application.

## Features
* username-password based login and register
* schema creation for user
* jwt based token management

## Usage
In your main function, call `auth.Setup()` as shown,
```go
package main

import "github.com/bdfp/auth"

func main() {
    db, err := sql.Open("mysql", "root:morning_star@/diary")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	log.Println("Database connection opened")

	// Set up auth
	auth.Setup(db)
}
```

## Database changes
the following user table is created in the passwed database

    | user        |
    | ------------- |
    | id
    | username      |
    | password |

 ## API
This will also register login and register http handlers as described

    ### Path /login
    * **Method** POST
    * **Request Object**
        ```json
        {
            "username": "your_username",
            "password": "your_password"
        }
        ```
    * **Response Object**
        ```json
        {
            "token": "eiofgheriuygij wer",
            "user": {
                "username": "your_username"
            }
        }
        ```

    ### Path  /register
    * **Method** POST
    * **Request Object**
        ```json
        {
            "username": "your_username",
            "password": "your_password"
        }
        ```
    * **Response Object**
        ```json
        {
        "message": "Success message"
    }
    ```

## Middleware
In order to secure any route with jwt, just use the `auth.Secure(yourHandler)` while registering your router
Example:
```go
func startHTTPServer(db *sql.DB) {
	router := httprouter.New()
	router.POST("/tasks", auth.Secure(myHandler))
	log.Fatal(http.ListenAndServe(":8484", router))
}
```

