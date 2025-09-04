package model

type Model map[int64]Campaign

type Campaign struct {
	Title       string
	Description string
	Groups      []Group
}

type Group struct {
	Title string
}
