package db

import (
	"context"
)

type FeedBack struct {
    Id         int64
    AudienceId int64
    MeetingId  int64
    Type       string
    Value      int64
    
}

func AddFeedBack(ctx context.Context, feedBack FeedBack) (int64, error) {
    query := `
    insert into
    feedback(audience_id, meeting_id, type, value)
    values(?, ?, ?, ?)
    `
    return commonExec(ctx, query,
        feedBack.AudienceId,
        feedBack.MeetingId,
        feedBack.Type,
        feedBack.Value,
    )
}

func DelFeedBack(
        ctx         context.Context,
        audience_id int64,
        meeting_id  int64,
    ) (int64, error) {

    query := `delete from feedBack where audience_id = ? and meeting_id = ?`
    return commonExec(ctx, query, audience_id, meeting_id)
}

func GetFeedBacks(
        ctx    context.Context,
        limit  int64,
        offset int64,
        where  string,
        args   ...any,
    ) ([]FeedBack, error) {

    query := `select id, audience_id, meeting_id, type, value from feedback`
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

func GetFeedBacksByAudienceAndMeeting(
        ctx         context.Context,
        audience_id int64,
        meeting_id  int64,
    ) ([]FeedBack, error) {

    where := `audience_id = ? and meeting_id = ?`
    return GetFeedBacks(ctx, 0, 0, where, audience_id, meeting_id)
}
