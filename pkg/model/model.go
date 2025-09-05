package model

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
