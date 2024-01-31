package db

import (
	"context"
)

type User struct {
    Id int64
    Name string
    Email string
    Password string
    Status string
}

func AddUser(ctx context.Context, user User) (int64, error) {
    query := `
    insert into
    user(name, email, password, status)
    values(?, ?, ?, ?)
    `
    return commonExec(ctx, query,
        user.Name,
        user.Email,
        hashString(user.Password),
        user.Status,
    )
}

func DelUser(ctx context.Context, id int64) (int64, error) {
    query := `delete from user where id = ?`
    return commonExec(ctx, query, id)
}

func SetUser(ctx context.Context, user User) (int64, error) {
    query := `
    update user set
    name = ?, email = ?, password = ?, status = ?
    where id = ?
    `
    return commonExec(ctx, query,
        user.Name,
        user.Email,
        user.Password,
        user.Status,
    )
}

func GetUserById(ctx context.Context, id int64) (User, error) {
    query := `
    select
    id, name, email, password, status
    from user where id = ?
    `
    row := adaptiveQueryRow(ctx, query, id)

    user := User{Id: -1}
    err := commonScanRow(row,
        &user.Id,
        &user.Name,
        &user.Email,
        &user.Password,
        &user.Status,
    )
    return user, err
}

func GetUsers(
        ctx    context.Context,
        limit  int64,
        offset int64,
        where  string,
        args   ...any,
    ) ([]User, error) {

    query := `select id, name, email, status from user`
    rows, err := getRows(ctx, query, limit, offset, where, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []User
    for rows.Next() {
        var item User
        if err := rows.Scan(
            &item.Id,
            &item.Name,
            &item.Email,
            &item.Status,
        ); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func GetUsersByPage(
        ctx       context.Context,
        ix_page   int64,
        page_size int64,
    ) ([]User, error) {

    limit := page_size
    offset := ix_page * page_size
    return GetUsers(ctx, limit, offset, "")
}

func GetUsersTotal(ctx context.Context) (int64, error) {
    return getTotal(ctx, "user")
}

func GetUsersPageTotal(ctx context.Context, size int64) (int64, error) {
    return getPageTotal(ctx, "user", size)
}
