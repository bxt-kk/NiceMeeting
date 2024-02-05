package api

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    db "nicemeeting/db"
)

func GetFeedBackByIds(r *gin.RouterGroup) {
    type Args struct {
        Size   int64 `uri:"size"`
        Page   int64 `uri:"page"`
        AudienceId int64 `uri:"audience_id"`
        MeetingId   int64 `uri:"meeting_id"`
    }
    r.GET("/feedback/:size/:page/:audience_id/:meeting_id", func (c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindUri(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        feedbacks, err := db.GetFeedBacksByPage(
            c, args.Size, args.Page, args.AudienceId, args.MeetingId)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, feedbacks)
    })
}

func AddFeedBack(r *gin.RouterGroup) {
    r.POST("/feedback", func (c *gin.Context) {
        args := db.FeedBack{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        user_id := args.AudienceId
        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(user_id) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user do not march"})
            return
        }

        if _, err := db.AddFeedBack(c, args); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}

func DelFeedBack(r *gin.RouterGroup) {
    r.POST("/feedback", func (c *gin.Context) {
        args := db.FeedBack{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        user_id := args.AudienceId
        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(user_id) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user do not march"})
            return
        }

        if _, err := db.DelFeedBack(c, user_id, args.MeetingId, args.Type); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}
