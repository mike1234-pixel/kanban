package store

import (
	"context"
	"kanban/internal/models"
)

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
	SaveBoard(ctx context.Context, board models.Board) error
	GetBoards(ctx context.Context) ([]models.Board, error)
	GetBoardByID(ctx context.Context, id string) (models.Board, error)
	UpdateBoard(ctx context.Context, board models.Board) (bool, error)
	DeleteBoard(ctx context.Context, id string) (bool, error)
	//#endregion
}
