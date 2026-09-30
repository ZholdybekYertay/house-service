package model

import "time"

type ServiceRequest struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	Problem   string    `json:"problem"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateRequestInput struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Problem string `json:"problem"`
}

type UpdateStatusInput struct {
	Status string `json:"status"`
}
