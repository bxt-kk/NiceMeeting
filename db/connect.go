package db

import (
	"context"
)

type Link struct {
    Id          int64
    FromId     int64
    ToId       int64
    Description string
    
}

func AddLink(ctx context.Context, link Link) (int64, error) {
    query := `
    insert into
    link(from_id, to_id, description)
    values(?, ?, ?)
    `
    return commonExec(ctx, query,
        link.FromId,
        link.ToId,
        link.Description,
    )
}

func DelLink(
        ctx     context.Context,
        from_id int64,
        to_id   int64,
    ) (int64, error) {

    query := `delete from link where from_id = ? and to_id = ?`
    return commonExec(ctx, query, from_id, to_id)
}

func GetLinks(
        ctx    context.Context,
        limit  int64,
        offset int64,
        where  string,
        args   ...any,
    ) ([]Link, error) {

    query := `select id, from_id, to_id, description from link`
    rows, err := getRows(ctx, query, limit, offset, where, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []Link
    for rows.Next() {
        var item Link
        if err := rows.Scan(
            &item.Id,
            &item.FromId,
            &item.ToId,
            &item.Description,
        ); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func GetLinksByFromId(
        ctx     context.Context,
        limit   int64,
        offset  int64,
        from_id int64,
    ) ([]Link, error) {

    where := `from_id = ?`
    return GetLinks(ctx, limit, offset, where, from_id)
}

func GetLinksByToId(
        ctx    context.Context,
        limit  int64,
        offset int64,
        to_id  int64,
    ) ([]Link, error) {

    where := `to_id = ?`
    return GetLinks(ctx, limit, offset, where, to_id)
}

func GetLinksByFromIdAndToId(
        ctx     context.Context,
        from_id int64,
        to_id   int64,
    ) ([]Link, error) {

    where := `from_id = ? and to_id = ?`
    return GetLinks(ctx, 0, 0, where, from_id, to_id)
}
