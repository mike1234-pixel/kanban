package models

import "time"

// #region Board
type Board struct {
	ID        string    `json:"id" validate:"required,uuid"`
	Title     string    `json:"title" validate:"required,min=3,max=100"`
	Columns   []Column  `json:"columns" validate:"dive"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

//#endregion

// #region Column
type Column struct {
	ID        string    `json:"id" validate:"required,uuid"`
	BoardID   string    `json:"board_id" validate:"required,uuid"`
	Title     string    `json:"title" validate:"required,min=1,max=100"`
	Position  int       `json:"position" validate:"gte=0"`
	Cards     []Card    `json:"cards,omitempty" validate:"dive"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

//#endregion

// #region Card
type Card struct {
	ID          string `json:"id" validate:"required,uuid"`
	ColumnID    string `json:"column_id" validate:"required,uuid"`
	Title       string `json:"title" validate:"required,min=1,max=150"`
	Description string `json:"description" validate:"max=1000"`
	Order       int    `json:"order" validate:"gte=0"`
}

//#endregion

// #region User
type User struct {
	ID        string    `json:"id" validate:"required,uuid"`
	Name      string    `json:"name" validate:"required,min=1,max=100"`
	Email     string    `json:"email" validate:"required,email,max=254"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

//#endregion
