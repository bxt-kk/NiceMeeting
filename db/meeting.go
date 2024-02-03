package db

import (
	"context"
)

type Meeting struct {
    Id             int64 `json:"id"`
    Title          string `json:"title"`
    OwnerId        int64 `json:"owner_id"`
    Content        string `json:"content"`
    Status         string `json:"status"`
    CreationTime   int64 `json:"creation_time"`
    LastEditedTime int64 `json:"last_edited_time"`
}

func AddMeeting(ctx context.Context, meeting Meeting) (int64, error) {
    query := `insert into
    meeting(title, owner_id, content, status, creation_time)
    values(?, ?, ?, ?, ?)
    `
    return commonExec(ctx, query,
        meeting.Title,
        meeting.OwnerId,
        meeting.Content,
        meeting.Status,
        nowUnix(),
    )
}

func DelMeeting(ctx context.Context, id, owner_id int64) (int64, error) {
    query := `delete from meeting where id = ? and owner_id = ?`
    return commonExec(ctx, query, id, owner_id)
}

func SetMeeting(ctx context.Context, meeting Meeting) (int64, error) {
    query := `update meeting set
    title = ?, content = ?, status = ?, last_edited_time = ?
    where id = ? and owner_id = ?
    `
    return commonExec(ctx, query,
        meeting.Title,
        meeting.Content,
        meeting.Status,
        nowUnix(),
        meeting.Id,
        meeting.OwnerId,
    )
}

func GetMeetingById(ctx context.Context, id int64) (Meeting, error) {
    query := `select
    id, title, owner_id, content, status, creation_time, last_edited_time
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
        &meeting.CreationTime,
        &meeting.LastEditedTime,
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

    query := `select
    id, title, owner_id, status, creation_time, last_edited_time
    from meeting`
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
            &item.Status,
            &item.CreationTime,
            &item.LastEditedTime,
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

func GetMeetingsByTag(
        ctx    context.Context,
        limit  int64,
        offset int64,
        tag    string,
    ) ([]Meeting, error) {

    query := `select
    id, title, owner_id, status, creation_time, last_edited_time
    from meeting inner join tag on
    meeting.id = tag.meeting_id and tag.label = ?`
    rows, err := getRows(ctx, query, limit, offset, "", tag)
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
            &item.Status,
            &item.CreationTime,
            &item.LastEditedTime,
        ); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func GetTagMeetingsByPage(
        ctx       context.Context,
        page_size int64,
        ix_page   int64,
        tag       string,
    ) ([]Meeting, error) {

    limit := page_size
    offset := ix_page * page_size
    return GetMeetingsByTag(ctx, limit, offset, tag)
}
