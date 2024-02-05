package api

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    db "nicemeeting/db"
)

func GetSite(r *gin.RouterGroup) {
    r.GET("/site", func (c *gin.Context) {
        site, err := db.GetSite(c)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        if site.Id == -1 {
            c.JSON(http.StatusNotFound, gin.H{"desc": "not found site by id"})
            return
        }
        c.JSON(http.StatusOK, site)
    })
}

func SetSite(r *gin.RouterGroup) {
    type Args struct {
        Site   db.Site `json:"site"`
        UserId int64 `json:"user_id"`
    }
    r.POST("/site", func (c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(args.UserId) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user is not marched"})
            return
        }
        user, err := db.GetUserByIdSafe(c, args.UserId)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        if user.Id == -1 {
            c.JSON(http.StatusNotFound, gin.H{"desc": "not found user by id"})
            return
        }
        // Note! to define levels
        if user.Level & 0b100 == 0 {
            c.AbortWithStatus(http.StatusUnauthorized)
            return
        }

        if _, err := db.SetSite(c, args.Site);err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}
