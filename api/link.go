package api

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    db "nicemeeting/db"
)

func GetLinksByIds(r *gin.RouterGroup) {
    type Args struct {
        Size   int64 `uri:"size"`
        Page   int64 `uri:"page"`
        FromId int64 `uri:"from_id"`
        ToId   int64 `uri:"to_id"`
    }
    r.GET("/link/:size/:page/:from_id/:to_id", func (c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindUri(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        links, err := db.GetLinksByPage(
            c, args.Size, args.Page, args.FromId, args.ToId)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, links)
    })
}

func AddLink(r *gin.RouterGroup) {
    r.POST("/link", func (c *gin.Context) {
        args := db.Link{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        from_id := args.FromId
        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(from_id) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user do not march"})
            return
        }

        if _, err := db.AddLink(c, args); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}

func DelLink(r *gin.RouterGroup) {
    r.POST("/link", func (c *gin.Context) {
        args := db.Link{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        from_id := args.FromId
        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(from_id) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user do not march"})
            return
        }

        if _, err := db.DelLink(c, from_id, args.ToId); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}
