package handlers

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    db "nicemeeting/db"
)

func HandlerRegister(r *gin.RouterGroup) {
    r.POST("/register", func (c *gin.Context) {
        args := db.User{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        _, err := db.AddUser(c, args)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}

func HandlerLogin(r *gin.RouterGroup) {
    type Args struct {
        Email    string `json:"email"`
        Password string `json:"password"`
    }
    r.POST("/login", func (c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        user, err := db.GetUserForLogin(c, args.Email, args.Password)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        if user.Id == -1 {
            c.JSON(http.StatusNotFound, gin.H{"desc": "email or password do not march"})
            return
        }
        session := sessions.Default(c)
        session.Set("hash_id", hashId(user.Id))
        session.Set("login", "yes")
        if err = session.Save(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
        }
        c.JSON(http.StatusOK, user)
    })
}

func HandlerLogout(r *gin.RouterGroup) {
    r.GET("/logout", func (c *gin.Context) {
        session := sessions.Default(c)
        session.Clear()
        if err := session.Save(); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}
