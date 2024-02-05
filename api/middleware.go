package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"

    db "nicemeeting/db"
)

func Sessions() gin.HandlerFunc {
    store := cookie.NewStore([]byte(db.SECURITY_KEY))
    return sessions.Sessions("nms", store)
}

func Authorization() gin.HandlerFunc {
    return func (c *gin.Context) {
        if IsUnauthorized(c) {
            return
        }
        c.Next()
    }
}

func IsUnauthorized(c *gin.Context) bool {
    session := sessions.Default(c)
    if session.Get("login") != "yes" {
        c.AbortWithStatus(http.StatusUnauthorized)
        return true
    }
    return false
}
