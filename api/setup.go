package api

import (
    "github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine, path string) {
    api := r.Group(path)
    add := api.Group("/add", Authorization())
    set := api.Group("/set", Authorization())
    del := api.Group("/del", Authorization())

    // admin
    GetSite(api)
    SetSite(set)

    // feedback
    GetFeedBackByIds(api)
    AddFeedBack(add)
    DelFeedBack(del)

    // link
    GetLinksByIds(api)
    AddLink(add)
    DelLink(del)

    // login
    Register(api)
    Login(api)
    Logout(api)

    // meeting
    GetMeeting(api)
    GetMeetings(api)
    GetMeetingsByTag(api)
    SetMeeting(set)
    AddMeeting(add)
    DelMeeting(del)

    // tag
    GetTags(api)
    AddTag(add)
    DelTag(del)

    // user
    GetUser(api)
    GetUsers(api)
    SetUser(set)
}
