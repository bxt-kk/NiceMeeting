package db

import (
	"context"
)

type Tag struct {
    Id        int64 `json:"id"`
    MeetingId int64 `json:"meeting_id"`
    Label     string `json:"label"`
}

func AddTag(ctx context.Context, tag Tag) (int64, error) {
    query := `insert into tag(meeting_id, label) values(?, ?)`
    return commonExec(ctx, query, tag.MeetingId, tag.Label)
}

func DelTag(ctx context.Context, meeting_id int64, label string) (int64, error) {
    query := `delete from tag where meeting_id = ? and label = ?`
    return commonExec(ctx, query, meeting_id, label)
}

func GetTags(
        ctx    context.Context,
        limit  int64,
        offset int64,
        where  string,
        args   ...any,
    ) ([]Tag, error) {

    query := `select id, meeting_id, label from tag`
    rows, err := getRows(ctx, query, limit, offset, "", where, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    items := []Tag{}
    for rows.Next() {
        var item Tag
        if err := rows.Scan(
            &item.Id,
            &item.MeetingId,
            &item.Label,
        ); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func GetTagsByPage(
        ctx       context.Context,
        page_size int64,
        ix_page   int64,
    ) ([]Tag, error) {

    limit := page_size
    offset := ix_page * page_size
    return GetTags(ctx, limit, offset, "")
}

func GetTagsTotal(ctx context.Context) (int64, error) {
    return getTotal(ctx, "tag")
}

func GetTagsPageTotal(ctx context.Context, size int64) (int64, error) {
    return getPageTotal(ctx, "tag", size)
}
