package main

import (
    "fmt"
    "flag"
    "net/http"

    "github.com/gin-gonic/gin"

    db "nicemeeting/db"
    hd "nicemeeting/handlers"
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
        })
    }

    fmt.Println("Hello Nice-Meeting!")

    db.Connect()

    router := gin.Default()
    router.SetTrustedProxies(nil)
    router.Use(hd.Sessions())

    common := router.Group("/")
    authorized := router.Group("/auth", hd.Authorization())

    hd.HandlerUser(common)
    hd.HandlerTags(common)
    hd.HandlerMeetings(common)
    hd.HandlerMeetingsByTag(common)
    hd.HandlerLinksByIds(common)

    authorized.GET("/try", func (c *gin.Context) {
        c.IndentedJSON(http.StatusOK, gin.H{"hello": "world"})
    })

    router.Run("localhost:8080")
}
