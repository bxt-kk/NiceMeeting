package handlers

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    db "nicemeeting/db"
)

func HandlerTags(r *gin.RouterGroup) {
    type Args struct {
        Size int64 `uri:"size"`
        Page int64 `uri:"page"`
    }
    r.GET("/tags/:size/:page", func (c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindUri(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        tags, err := db.GetTagsByPage(c, args.Size, args.Page)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, tags)
    })
}

func HandlerAddTag(r *gin.RouterGroup) {
    r.POST("/add/tag", func (c *gin.Context) {
        args := db.Tag{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        meeting_id := args.MeetingId
        meeting, err := db.GetMeetingById(c, meeting_id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        if meeting.Id == -1 {
            c.JSON(http.StatusNotFound, gin.H{"desc": "not found meeting by id"})
            return
        }

        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(meeting.OwnerId) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user is not marched"})
            return
        }

        if _, err := db.AddTag(c, args); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}

func HandlerDelTag(r *gin.RouterGroup) {
    r.POST("/del/tag", func (c *gin.Context) {
        args := db.Tag{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        meeting_id := args.MeetingId
        meeting, err := db.GetMeetingById(c, meeting_id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        if meeting.Id == -1 {
            c.JSON(http.StatusNotFound, gin.H{"desc": "not found meeting by id"})
            return
        }

        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(meeting.OwnerId) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user do not march"})
            return
        }

        if _, err := db.DelTag(c, args.Id, meeting_id); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}
