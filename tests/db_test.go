package test

import (
    "log"
    "testing"
    "nicemeeting/db"
)

var sample_users = []db.User{
    {Name: "kk", Email: "kk@kmail.com", Password: "1234567k"},
    {Name: "jojo", Email: "jojo@kmail.com", Password: "1234567j"},
}

func TestInit(t *testing.T) {
    db.SetDBFile("/tmp/_test_nicemeeting.db")
    err := db.Init(sample_users[0])
    if err != nil {
        t.Error(err)
    }
}

func TestConnect(t *testing.T) {
    db.SetDBFile("/tmp/_test_nicemeeting.db")
    err := db.Connect()
    if err != nil {
        t.Error(err)
    }
}

func TestAddUser(t *testing.T) {
    id, err := db.AddUser(nil, sample_users[1])
    if err != nil {
        t.Error(err)
    }
    log.Printf("insert an user[id=%d]\n", id)
}

func TestGetUserById(t *testing.T) {
    user, err := db.GetUserById(nil, 1)
    if err != nil {
        t.Error(err)
    }
    log.Println(user)
}

func TestGetUsers(t *testing.T) {
    users, err := db.GetUsers(nil, 0, 0, "")
    if err != nil {
        t.Error(err)
    }
    for ix, user := range users {
        log.Println(ix, user)
        reverse_ix := len(sample_users) - ix - 1
        if user.Email != sample_users[reverse_ix].Email {
            t.Errorf("data of users not marched")
        }
    }
    if len(users) != 2 {
        t.Errorf("number of users not marched")
    }
}

func TestGetUsersTotal(t *testing.T) {
    total, err := db.GetUsersTotal(nil)
    if err != nil {
        t.Error(err)
    }
    if total != int64(len(sample_users)) {
        t.Errorf("total of users not marched")
    }
}
