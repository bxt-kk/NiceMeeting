package db

import (
	"context"
)

type Link struct {
    Id          int64 `json:"id"`
    FromId      int64 `json:"from_id"`
    ToId        int64 `json:"to_id"`
    Description string `json:"description"`
    
}

func AddLink(ctx context.Context, link Link) (int64, error) {
    query := `insert into
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

func SetLink(ctx context.Context, link Link) (int64, error) {
    query := `update link set description = ?  where id = ?`
    return commonExec(ctx, query, link.Description, link.Id)
}

func GetLinks(
        ctx    context.Context,
        limit  int64,
        offset int64,
        where  string,
        args   ...any,
    ) ([]Link, error) {

    query := `select id, from_id, to_id, description from link`
    rows, err := getRows(ctx, query, limit, offset, "", where, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    items := []Link{}
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

func GetLinksByIds(
        ctx     context.Context,
        limit   int64,
        offset  int64,
        from_id int64,
        to_id   int64,
    ) ([]Link, error) {

    if from_id == -1 && to_id == -1 {
        return GetLinks(ctx, limit, offset, "")
    } else if from_id != -1 && to_id == -1 {
        return GetLinksByFromId(ctx, limit, offset, from_id)
    } else if from_id == -1 && to_id != -1 {
        return GetLinksByToId(ctx, limit, offset, from_id)
    } 
    where := `from_id = ? and to_id = ?`
    return GetLinks(ctx, limit, offset, where, from_id, to_id)
}

func GetLinksByPage(
        ctx       context.Context,
        page_size int64,
        ix_page   int64,
        from_id   int64,
        to_id     int64,
    ) ([]Link, error) {

    limit := page_size
    offset := ix_page * page_size
    return GetLinksByIds(ctx, limit, offset, from_id, to_id)
}

func GetLinksTotal(ctx context.Context) (int64, error) {
    return getTotal(ctx, "link")
}

func GetLinksPageTotal(ctx context.Context, size int64) (int64, error) {
    return getPageTotal(ctx, "link", size)
}
