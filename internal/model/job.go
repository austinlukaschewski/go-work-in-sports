package model

type Job struct {
	Link            string `json:"link"`
	Org             string `json:"org"`
	Title           string `json:"title"`
	Date            string `json:"date,omitempty"`
	SincePostedDate string `json:"since_posted_date,omitempty"`
}
