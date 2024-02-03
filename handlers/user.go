package handlers

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    db "nicemeeting/db"
)

func HandlerUser(r *gin.RouterGroup) {
    type Args struct {
        Id int64 `uri:"id"`
    }
    r.GET("/user/:id", func (c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindUri(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        user, err := db.GetUserByIdSafe(c, args.Id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        if user.Id == -1 {
            c.JSON(http.StatusNotFound, gin.H{"desc": "not found user by id"})
            return
        }
        c.JSON(http.StatusOK, user)
    })
}

func HandlerUsers(r *gin.RouterGroup) {
    type Args struct {
        Size int64 `uri:"size"`
        Page int64 `uri:"page"`
    }
    r.GET("/users/:size/:page", func(c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindUri(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        users, err := db.GetUsersByPage(c, args.Size, args.Page)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, users)
    })
}

func HandlerSetUser(r *gin.RouterGroup) {
    type Args struct {
        User db.User `json:"user"`
        LastPassword string `json:"last_password"`
    }
    r.POST("/set/user", func (c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        user_id := args.User.Id
        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(user_id) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user do not march"})
            return
        }

        var err error
        last_password := args.LastPassword
        if last_password == "" {
            _, err = db.UserSetLowRisk(c, args.User)
        } else {
            _, err = db.UserSet(c, args.User, last_password)
        }
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}
