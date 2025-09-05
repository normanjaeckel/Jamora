package model

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type Model map[int64]Campaign

type Campaign struct {
	Id          int64
	Title       string
	Description string
}

const CampaignTableQuery = `
	CREATE TABLE IF NOT EXISTS campaigns (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT NOT NULL
	);`

func CampaignGetAll(ctx context.Context, db *sql.DB) ([]Campaign, error) {
	rows, err := db.QueryContext(ctx, "SELECT id, title, description FROM campaigns")
	if err != nil {
		return nil, fmt.Errorf("SELECT campaigns from database: %w", err)
	}
	defer rows.Close()

	var campaigns []Campaign
	for rows.Next() {
		var c Campaign
		if err := rows.Scan(&c.Id, &c.Title, &c.Description); err != nil {
			return nil, fmt.Errorf("scan campaigns from database query response: %w", err)
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, nil
}

func CampaignGet(ctx context.Context, db *sql.DB, idStr string) (Campaign, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Just send the error as if there is no object with this (invalid) id so we can send HTTP 400 later on.
		return Campaign{}, sql.ErrNoRows
	}
	row := db.QueryRowContext(ctx, "SELECT id, title, description FROM campaigns WHERE id=$1", id)
	var c Campaign
	if err := row.Scan(&c.Id, &c.Title, &c.Description); err != nil {
		return Campaign{}, fmt.Errorf("scan campaigns from database query response: %w", err)
	}
	return c, nil
}

type Group struct {
	Id         int64
	Title      string
	CampaignId int64
}

const GroupTableQuery = `
	CREATE TABLE IF NOT EXISTS groups (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		campaign_id INTEGER NOT NULL,
		FOREIGN KEY(campaign_id) REFERENCES campaigns(id) ON DELETE CASCADE
	);`
