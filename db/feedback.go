package db

import (
	"context"
)

type FeedBack struct {
    Id         int64 `json:"id"`
    AudienceId int64 `json:"audience_id"`
    MeetingId  int64 `json:"meeting_id"`
    Type       string `json:"type"`
    Value      int64 `json:"value"`
    Time       int64 `json:"time"`
    
}

func AddFeedBack(ctx context.Context, feedBack FeedBack) (int64, error) {
    query := `insert into
    feedback(audience_id, meeting_id, type, value)
    values(?, ?, ?, ?)
    `
    return commonExec(ctx, query,
        feedBack.AudienceId,
        feedBack.MeetingId,
        feedBack.Type,
        feedBack.Value,
        nowUnix(),
    )
}

func DelFeedBack(
        ctx         context.Context,
        audience_id int64,
        meeting_id  int64,
        _type       string,
    ) (int64, error) {

    query := `delete from feedBack where
    audience_id = ? and meeting_id = ? and type = ?`
    return commonExec(ctx, query, audience_id, meeting_id, _type)
}

func GetFeedBacks(
        ctx    context.Context,
        limit  int64,
        offset int64,
        where  string,
        args   ...any,
    ) ([]FeedBack, error) {

    query := `select
    id, audience_id, meeting_id, type, value, time from feedback`
    rows, err := getRows(ctx, query, limit, offset, where, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []FeedBack
    for rows.Next() {
        var item FeedBack
        if err := rows.Scan(
            &item.Id,
            &item.AudienceId,
            &item.MeetingId,
            &item.Type,
            &item.Value,
            &item.Time,
        ); err != nil {
            return nil, err
        }
        items = append(items, item)
    }
    return items, rows.Err()
}

func GetFeedBacksByAudience(
        ctx         context.Context,
        limit       int64,
        offset      int64,
        audience_id int64,
    ) ([]FeedBack, error) {

    where := `audience_id = ?`
    return GetFeedBacks(ctx, limit, offset, where, audience_id)
}

func GetFeedBacksByMeeting(
        ctx        context.Context,
        limit      int64,
        offset     int64,
        meeting_id int64,
    ) ([]FeedBack, error) {

    where := `meeting_id = ?`
    return GetFeedBacks(ctx, limit, offset, where, meeting_id)
}

func GetFeedBacksByIds(
        ctx         context.Context,
        limit       int64,
        offset      int64,
        audience_id int64,
        meeting_id  int64,
    ) ([]FeedBack, error) {

    if audience_id == -1 && meeting_id == -1 {
        return GetFeedBacks(ctx, limit, offset, "")
    } else if audience_id != -1 && meeting_id == -1 {
        return GetFeedBacksByAudience(ctx, limit, offset, audience_id)
    } else if audience_id == -1 && meeting_id != -1 {
        return GetFeedBacksByMeeting(ctx, limit, offset, meeting_id)
    } 
    where := `audience_id = ? and meeting_id = ?`
    return GetFeedBacks(ctx, limit, offset, where, audience_id, meeting_id)
}

func GetFeedBacksByPage(
        ctx         context.Context,
        page_size   int64,
        ix_page     int64,
        audience_id int64,
        meeting_id  int64,
    ) ([]FeedBack, error) {

    limit := page_size
    offset := ix_page * page_size
    return GetFeedBacksByIds(ctx, limit, offset, audience_id, meeting_id)
}

func GetFeedBacksTotal(ctx context.Context) (int64, error) {
    return getTotal(ctx, "feedback")
}

func GetFeedBacksPageTotal(ctx context.Context, size int64) (int64, error) {
    return getPageTotal(ctx, "feedback", size)
}
