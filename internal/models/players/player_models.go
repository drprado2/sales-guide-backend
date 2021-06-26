package players

import "time"

type Player struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePlayerRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type UpdatePlayerRequest struct {
	Name  string `json:"name"`
	Image string `json:"image"`
}

type PlayerFilter struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
