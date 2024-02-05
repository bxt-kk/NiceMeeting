package main

import (
    "fmt"
    "flag"
    "net/http"

    "github.com/gin-gonic/gin"

    db "nicemeeting/db"
    api "nicemeeting/api"
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
    router.Use(api.Sessions())

    router.LoadHTMLGlob("./templates/*.html")
    router.Static("/static", "./static")

    page := router.Group("/")
    // authorized := api.Group("/auth", api.Authorization())

    api.Setup(router, "api")
    page.GET("/", func (c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{
            "title": "index", "try": []string{"a", "b", "c"}})
    })

    page.GET("/try", func (c *gin.Context) {
        c.YAML(http.StatusOK, []struct{
            Name string `yaml:"name"`
            Age int64 `yaml:"age"`}{
            {"kk", 13},
            {"jojo", 14},
        })})

    router.Run("localhost:8080")
}
