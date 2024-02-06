package db

import (
	"context"
)

type User struct {
    Id               int64 `json:"id"`
    Name             string `json:"name"`
    Email            string `json:"email"`
    Password         string `json:"password"`
    Status           string `json:"status"`
    RegistrationTime int64 `json:"registration_time"`
    LastLoginTime    int64 `json:"last_login_time"`
    Level            int64 `json:"level"`
}

func AddUser(ctx context.Context, user User) (int64, error) {
    query := `insert into
    user(name, email, password, status, registration_time, level)
    values(?, ?, ?, ?, ?, ?)
    `
    return commonExec(ctx, query,
        user.Name,
        user.Email,
        HashString(user.Password),
        user.Status,
        nowUnix(),
        user.Level,
    )
}

func DelUser(ctx context.Context, id int64) (int64, error) {
    query := `delete from user where id = ?`
    return commonExec(ctx, query, id)
}

func UserSetLowRisk(ctx context.Context, user User) (int64, error) {
    query := `update user set
    name = ?, email = ?
    where id = ?
    `
    return commonExec(ctx, query,
        user.Name,
        user.Email,
        user.Id,
    )
}

func UserSet(ctx context.Context, user User, last_password string) (int64, error) {
    query := `update user set
    name = ?, email = ?, password = ?
    where id = ? and password = ?
    `
    return commonExec(ctx, query,
        user.Name,
        user.Email,
        HashString(user.Password),
        user.Id,
        HashString(last_password),
    )
}

func SetUser(ctx context.Context, user User) (int64, error) {
    query := `update user set
    name = ?, email = ?, password = ?, status = ?, level = ?
    where id = ?
    `
    return commonExec(ctx, query,
        user.Name,
        user.Email,
        HashString(user.Password),
        user.Status,
        user.Level,
        user.Id,
    )
}

func UpdateLoginTime(ctx context.Context, user User) (int64, error) {
    query := `update user set last_login_time = ? where id = ?`
    return commonExec(ctx, query, nowUnix(), user.Id)
}

func GetUserById(ctx context.Context, id int64) (User, error) {
    query := `select
    id, name, email, password, status, registration_time, last_login_time, level
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
        &user.RegistrationTime,
        &user.LastLoginTime,
        &user.Level,
    )
    return user, err
}

func GetUserByIdSafe(ctx context.Context, id int64) (User, error) {
    user, err := GetUserById(ctx, id)
    user.Password = ""
    return user, err
}

func GetUserForLogin(ctx context.Context, email , password string) (User, error) {
    query := `select
    id, name, email, status, registration_time, last_login_time, level
    from user where email = ? and password = ?
    `
    row := adaptiveQueryRow(ctx, query, email, HashString(password))

    user := User{Id: -1}
    err := commonScanRow(row,
        &user.Id,
        &user.Name,
        &user.Email,
        &user.Status,
        &user.RegistrationTime,
        &user.LastLoginTime,
        &user.Level,
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

    query := `select
    id, name, email, status, registration_time, last_login_time
    from user`
    rows, err := getRows(ctx, query, limit, offset, "", where, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    items := []User{}
    for rows.Next() {
        var item User
        if err := rows.Scan(
            &item.Id,
            &item.Name,
            &item.Email,
            &item.Status,
            &item.RegistrationTime,
            &item.LastLoginTime,
        ); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func GetUsersByPage(
        ctx       context.Context,
        page_size int64,
        ix_page   int64,
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
