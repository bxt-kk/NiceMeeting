package main

import (
	"flag"
	"fmt"
	"log"

	// "net/http"

	"github.com/gin-gonic/gin"

	api "nicemeeting/api"
	db "nicemeeting/db"
	page "nicemeeting/page"
)

func main() {
    init_db  := flag.Bool("init-db", false, "To initializing database file")
    username := flag.String("username", "kk", "Set an user name")
    email    := flag.String("email", "kk@kmail.com", "Set an user mail")
    password := flag.String("password", "12345678k", "Set an user password")

    flag.Parse()

    fmt.Printf("init-db: %v\n", *init_db)
    if *init_db {
        db.Init(db.User{
            Name: *username,
            Email: *email,
            Password: *password,
            Level: 2,
        })
    }

    fmt.Println("Hello Nice-Meeting!")

    err := db.Connect()
    if err != nil {
        log.Fatal(err)
    }

    router := gin.Default()
    router.SetTrustedProxies(nil)
    router.Use(api.Sessions())

    api.Setup(router, "api")
    page.Setup(router, "/")

    router.Run("localhost:8080")
}
