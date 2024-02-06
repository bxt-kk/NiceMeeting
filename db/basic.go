package db

import (
    "fmt"
    "log"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"strings"
    "time"
)

func adaptiveExec(ctx context.Context, query string, args ...any) (sql.Result, error) {
    if ctx == nil {
        return db.Exec(query, args...)
    }
    return db.ExecContext(ctx, query, args...)
}

func adaptiveQueryRow(ctx context.Context, query string, args ...any) *sql.Row {
    if ctx == nil {
        return db.QueryRow(query, args...)
    }
    return db.QueryRowContext(ctx, query, args...)
}

func adaptiveQuery(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
    if ctx == nil {
        return db.Query(query, args...)
    }
    return db.QueryContext(ctx, query, args...)
}

func commonExec(ctx context.Context, query string, args ...any) (int64, error) {
    result, err := adaptiveExec(ctx, query, args...)
    if err != nil {
        return 0, err
    }
    if strings.HasPrefix(strings.TrimLeft(query, " \t\n"), "insert") {
        return result.LastInsertId()
    }
    return result.RowsAffected()
}

func commonScanRow(row *sql.Row, dst ...any) error {
    err := row.Scan(dst...)
    if err == sql.ErrNoRows {
        err = nil
    }
    return err
}

func getRows(
        ctx      context.Context,
        query    string,
        limit    int64,
        offset   int64,
        order_by string,
        where    string,
        args     ...any,
    ) (*sql.Rows, error) {

    if where != "" {
        query = fmt.Sprintf("%s where %s", query, where)
    }
    if order_by == "" {
        order_by = "id"
    }
    query = fmt.Sprintf("%s order by %s desc", query, order_by)
    if limit != 0 {
        query = fmt.Sprintf("%s limit %d offset %d", query, limit, offset)
    }
    return adaptiveQuery(ctx, query, args...)
}

func getTotal(ctx context.Context, table string) (int64, error) {
    query := fmt.Sprintf("select count(*) from %s", table)
    row := adaptiveQueryRow(ctx, query)
    var count int64
    err := row.Scan(&count)
    return count, err
}

func getPageTotal(
        ctx   context.Context,
        table string,
        size  int64,
    ) (int64, error) {

    count, err := getTotal(ctx, table)
    if err != nil {
        return 0, err
    }
    if size < 1 {
        log.Println("The number of the pages is must be a positive integer")
        size = 1
    }
    total := count / size
    if count % size != 0 {
        total += 1
    }
    return total, nil
}

func HashString(text string) string {
    hasher := sha256.New()
    hasher.Write([]byte(text + SECURITY_KEY))
    return hex.EncodeToString(hasher.Sum(nil))
}

func nowUnix() int64 {
    return time.Now().Unix()
}
