package projects

import "time"

type Project struct {
	ID        string    `json:"id"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}
