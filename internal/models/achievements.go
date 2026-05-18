package models

type Achievement struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IconData    string `json:"icon_data"`
}
