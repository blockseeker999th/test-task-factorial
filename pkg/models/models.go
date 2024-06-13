package models

import "time"

type Factorial struct {
	Id        int64     `json:"id"`
	ValueA    int       `json:"a" validate:"gte=0"`
	ValueB    int       `json:"b" validate:"gte=0"`
	CreatedAt time.Time `json:"createdAt"`
}
