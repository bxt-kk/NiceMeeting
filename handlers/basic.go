package handlers

import (
    "fmt"

    db "nicemeeting/db"
)

func hashId(Id int64) string {
    return db.HashString(fmt.Sprint(Id))
}
