package api

import (
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"

    db "nicemeeting/db"
)

func GetMeeting(r *gin.RouterGroup) {
    type Args struct {
        Id int64 `uri:"id"`
    }
    r.GET("/meeting/:id", func (c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindUri(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        meeting, err := db.GetMeetingById(c, args.Id)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        if meeting.Id == -1 {
            c.JSON(http.StatusNotFound, gin.H{"desc": "not found meeting by id"})
            return
        }
        c.JSON(http.StatusOK, meeting)
    })
}

func GetMeetings(r *gin.RouterGroup) {
    type Args struct {
        Size int64 `json:"size"`
        Page int64 `json:"page"`
    }
    r.GET("/meetings/:size/:page", func(c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindUri(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        meetings, err := db.GetMeetingsByPage(c, args.Size, args.Page)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, meetings)
    })
}

func GetMeetingsByTag(r *gin.RouterGroup) {
    type Args struct {
        Size int64 `json:"size"`
        Page int64 `json:"page"`
        Tag  string `json:"tag"`
    }
    r.GET("/meetings/:size/:page/:tag", func(c *gin.Context) {
        args := Args{}
        if err := c.ShouldBindUri(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        meetings, err := db.GetTagMeetingsByPage(
            c, args.Size, args.Page, args.Tag)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, meetings)
    })
}

func SetMeeting(r *gin.RouterGroup) {
    r.POST("/meeting", func (c *gin.Context) {
        args := db.Meeting{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        user_id := args.OwnerId
        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(user_id) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user is not marched"})
            return
        }

        if _, err := db.SetMeeting(c, args); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}

func AddMeeting(r *gin.RouterGroup) {
    r.POST("/meeting", func (c *gin.Context) {
        args := db.Meeting{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        user_id := args.OwnerId
        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(user_id) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user is not marched"})
            return
        }

        if _, err := db.AddMeeting(c, args); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}

func DelMeeting(r *gin.RouterGroup) {
    r.POST("/meeting", func (c *gin.Context) {
        args := db.Meeting{}
        if err := c.ShouldBindJSON(&args); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"desc": err})
            return
        }

        user_id := args.OwnerId
        session := sessions.Default(c)
        if session.Get("hash_id") != hashId(user_id) {
            c.JSON(http.StatusForbidden, gin.H{"desc": "user do not march"})
            return
        }

        if _, err := db.DelMeeting(c, args.Id, user_id); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"desc": "internal server error"})
            log.Panic(err)
            return
        }
        c.JSON(http.StatusOK, gin.H{"desc": ""})
    })
}
