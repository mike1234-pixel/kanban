package store

import (
	"context"
	"errors"
	"kanban/internal/models"
	"time"
)

// ErrEmailTaken is returned when a user is saved with an email that already belongs to another user.
var ErrEmailTaken = errors.New("email already in use")

// ErrBoardExists is returned when a board is saved with an ID that is already in use.
var ErrBoardExists = errors.New("board already exists")

type Store interface {
	//#region Card Operations
	SaveCard(ctx context.Context, card models.Card) error
	GetCards(ctx context.Context) ([]models.Card, error)
	GetCardByID(ctx context.Context, id string) (models.Card, error)
	UpdateCard(ctx context.Context, id string, updated models.Card) (bool, error)
	MoveCard(ctx context.Context, id string, columnID string) (bool, error)
	DeleteCard(ctx context.Context, id string) (bool, error)
	//#endregion

	//#region Column Operations
	SaveColumn(ctx context.Context, col models.Column) error
	GetColumnsWithCards(ctx context.Context, boardID string) ([]models.Column, error)
	GetColumnByID(ctx context.Context, id string) (models.Column, error)
	UpdateColumn(ctx context.Context, col models.Column) (bool, error)
	DeleteColumn(ctx context.Context, id string) (bool, error)
	//#endregion

	//#region Board operations
	SaveBoard(ctx context.Context, board models.Board) (time.Time, error)
	GetBoards(ctx context.Context) ([]models.Board, error)
	GetBoardByID(ctx context.Context, id string) (models.Board, error)
	UpdateBoard(ctx context.Context, board models.Board) (bool, error)
	DeleteBoard(ctx context.Context, id string) (bool, error)
	//#endregion

	//#region User operations
	SaveUser(ctx context.Context, user models.User) (time.Time, error)
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id string) (models.User, error)
	UpdateUser(ctx context.Context, user models.User) (bool, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	//#endregion

	//#region Session operations
	CreateSession(ctx context.Context, tokenHash, userID string, expiresAt time.Time) error
	GetSessionUser(ctx context.Context, tokenHash string) (models.User, error)
	//#endregion
}
