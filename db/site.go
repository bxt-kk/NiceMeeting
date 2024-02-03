package db

import (
	"context"
)

type Site struct {
    Id          int64 `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

func SetSite(ctx context.Context, site Site) (int64, error) {
    count, err := getTotal(ctx, "site")
    if err != nil {
        return 0, err
    }

    query := `update site set
    name = ?, description = ?
    where id = 1
    `
    if count < 1 {
        query = `insert into
        site(name, description)
        values(?, ?)
        `
    }
    return commonExec(ctx, query,
        site.Name,
        site.Description,
    )
}

func GetSite(ctx context.Context) (Site, error) {
    query := `select
    name, description
    from site where id = 1
    `
    row := adaptiveQueryRow(ctx, query)

    var site Site
    err := commonScanRow(row,
        &site.Name,
        &site.Description,
    )
    return site, err
}
