package db

import (
	"context"
)

type Meeting struct {
    Id       int64
    Title    string
    OwnerId int64
    Content  string
    Status   string
}

func AddMeeting(ctx context.Context, meeting Meeting) (int64, error) {
    query := `
    insert into
    meeting(title, owner_id, content, status)
    values(?, ?, ?, ?)
    `
    return commonExec(ctx, query,
        meeting.Title,
        meeting.OwnerId,
        meeting.Content,
        meeting.Status,
    )
}

func DelMeeting(ctx context.Context, id int64) (int64, error) {
    query := `delete from meeting where id = ?`
    return commonExec(ctx, query, id)
}

func SetMeeting(ctx context.Context, meeting Meeting) (int64, error) {
    query := `
    update meeting set
    title = ?, owner_id = ?, content = ?, status = ?
    where id = ?
    `
    return commonExec(ctx, query,
        meeting.Title,
        meeting.OwnerId,
        meeting.Content,
        meeting.Status,
    )
}

func GetMeetingById(ctx context.Context, id int64) (Meeting, error) {
    query := `
    select
    id, title, owner_id, content, status
    from meeting where id = ?
    `
    row := adaptiveQueryRow(ctx, query, id)

    meeting := Meeting{Id: -1}
    err := commonScanRow(row,
        &meeting.Id,
        &meeting.Title,
        &meeting.OwnerId,
        &meeting.Content,
        &meeting.Status,
    )
    return meeting, err
}

func GetMeetings(
        ctx    context.Context,
        limit  int64,
        offset int64,
        where  string,
        args   ...any,
    ) ([]Meeting, error) {

    query := `select id, title, owner_id, content, status from meeting`
    rows, err := getRows(ctx, query, limit, offset, where, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []Meeting
    for rows.Next() {
        var item Meeting
        if err := rows.Scan(
            &item.Id,
            &item.Title,
            &item.OwnerId,
            &item.Content,
            &item.Status,
        ); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func GetMeetingsByPage(
        ctx       context.Context,
        ix_page   int64,
        page_size int64,
    ) ([]Meeting, error) {

    limit := page_size
    offset := ix_page * page_size
    return GetMeetings(ctx, limit, offset, "")
}

func GetMeetingsTotal(ctx context.Context) (int64, error) {
    return getTotal(ctx, "meeting")
}

func GetMeetingsPageTotal(ctx context.Context, size int64) (int64, error) {
    return getPageTotal(ctx, "meeting", size)
}
