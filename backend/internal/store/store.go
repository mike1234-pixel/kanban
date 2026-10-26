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

// ErrColumnNotFound is returned when a card is moved to a column that doesn't
// exist or belongs to another user.
var ErrColumnNotFound = errors.New("column not found")

type Store interface {
	//#region Card Operations
	SaveCard(ctx context.Context, card models.Card, ownerID string) (bool, error)
	GetCards(ctx context.Context, ownerID string) ([]models.Card, error)
	GetCardByID(ctx context.Context, id, ownerID string) (models.Card, error)
	UpdateCard(ctx context.Context, id string, updated models.Card, ownerID string) (bool, error)
	MoveCard(ctx context.Context, id, columnID, ownerID string) (bool, error)
	DeleteCard(ctx context.Context, id, ownerID string) (bool, error)
	//#endregion

	//#region Column Operations
	SaveColumn(ctx context.Context, col models.Column, ownerID string) (bool, error)
	GetColumnsWithCards(ctx context.Context, boardID, ownerID string) ([]models.Column, error)
	GetColumnByID(ctx context.Context, id, ownerID string) (models.Column, error)
	UpdateColumn(ctx context.Context, col models.Column, ownerID string) (bool, error)
	DeleteColumn(ctx context.Context, id, ownerID string) (bool, error)
	//#endregion

	//#region Board operations
	SaveBoard(ctx context.Context, board models.Board) (time.Time, error)
	GetBoards(ctx context.Context, ownerID string) ([]models.Board, error)
	GetBoardByID(ctx context.Context, id, ownerID string) (models.Board, error)
	UpdateBoard(ctx context.Context, board models.Board) (bool, error)
	DeleteBoard(ctx context.Context, id, ownerID string) (bool, error)
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
	DeleteSession(ctx context.Context, tokenHash string) error
	//#endregion
}
