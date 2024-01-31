package main

import (
    "fmt"
    "flag"
    "net/http"

    "github.com/gin-gonic/gin"

    db "nicemeeting/db"
)

type album struct {
    ID     string  `json:"id"`
    Title  string  `json:"title"`
    Artist string  `json:"artist"`
    Price  float64 `json:"price"`
}

var albums = []album{
    {ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
    {ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
    {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
}

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
    router.GET("/albums", func (c *gin.Context) {
        c.IndentedJSON(http.StatusOK, albums)
    })

    router.Run("localhost:8080")
}
