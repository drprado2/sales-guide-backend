package repositories

import (
	"context"
	"github.com/drprado2/react-redux-typescript/internal/models"
	playerModels "github.com/drprado2/react-redux-typescript/internal/models/players"
	apptracer2 "github.com/drprado2/react-redux-typescript/pkg/apptracer"
	logs2 "github.com/drprado2/react-redux-typescript/pkg/logs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"runtime/debug"
)

type PlayerRepositoryInterface interface {
	Create(ctx context.Context, user *playerModels.CreatePlayerRequest) (*playerModels.Player, error)
	CreateTx(ctx context.Context, tx pgx.Tx, user *playerModels.CreatePlayerRequest) (*playerModels.Player, error)
	Update(ctx context.Context, playerId string, user *playerModels.UpdatePlayerRequest) (*playerModels.Player, error)
	UpdateTx(ctx context.Context, tx pgx.Tx, playerId string, user *playerModels.UpdatePlayerRequest) (*playerModels.Player, error)
	Delete(ctx context.Context, playerId string) error
	DeleteTx(ctx context.Context, tx pgx.Tx, playerId string) error
	GetById(ctx context.Context, playerId string) (*playerModels.Player, error)
	GetPaged(ctx context.Context, pagination *models.PaginationParameters, filter *playerModels.PlayerFilter) (*models.PaginationResponse, error)
}

type PlayerRepository struct {
	DbPool *pgxpool.Pool
	Tracer apptracer2.TracerService
}

func (r PlayerRepository) Create(ctx context.Context, user *playerModels.CreatePlayerRequest) (*playerModels.Player, error) {
	result := &playerModels.Player{}
	if err := r.DbPool.QueryRow(ctx, "INSERT INTO player (name, image) VALUES ($1, $2) RETURNING id, name, image, created_at, updated_at", user.Name, user.Image).
		Scan(&result.ID, &result.Name, &result.Image, &result.CreatedAt, &result.UpdatedAt); err != nil {
		logs2.Logger(ctx).WithError(err).Errorf("Error updating player by id, %v", err)
		return nil, err
	}
	return result, nil
}

func (r PlayerRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *playerModels.CreatePlayerRequest) (*playerModels.Player, error) {
	result := &playerModels.Player{}
	if err := tx.QueryRow(ctx, "INSERT INTO player (name, image) VALUES ($1, $2) RETURNING id, name, image, created_at, updated_at", user.Name, user.Image).
		Scan(&result.ID, &result.Name, &result.Image, &result.CreatedAt, &result.UpdatedAt); err != nil {
		logs2.Logger(ctx).WithError(err).Errorf("Error updating player by id, %v", err)
		return nil, err
	}
	return result, nil
}

func (r PlayerRepository) Update(ctx context.Context, playerId string, user *playerModels.UpdatePlayerRequest) (*playerModels.Player, error) {
	result := &playerModels.Player{}
	if err := r.DbPool.QueryRow(ctx, "UPDATE player SET name = $1, image = $2 WHERE id = $3 RETURNING id, name, image, created_at, updated_at", user.Name, user.Image, playerId).
		Scan(&result.ID, &result.Name, &result.Image, &result.CreatedAt, &result.UpdatedAt); err != nil {
		logs2.Logger(ctx).WithError(err).WithField("player_id", playerId).Errorf("Error updating player by id, %v", err)
		return nil, err
	}
	return result, nil
}

func (r PlayerRepository) UpdateTx(ctx context.Context, tx pgx.Tx, playerId string, user *playerModels.UpdatePlayerRequest) (*playerModels.Player, error) {
	result := &playerModels.Player{}
	if err := tx.QueryRow(ctx, "UPDATE player SET name = $1, image = $2 WHERE id = $3 RETURNING id, name, image, created_at, updated_at", user.Name, user.Image, playerId).
		Scan(&result.ID, &result.Name, &result.Image, &result.CreatedAt, &result.UpdatedAt); err != nil {
		logs2.Logger(ctx).WithError(err).WithField("player_id", playerId).Errorf("Error updating player by id, %v", err)
		return nil, err
	}
	return result, nil
}

func (r PlayerRepository) Delete(ctx context.Context, playerId string) error {
	_, err := r.DbPool.Exec(ctx, "DELETE FROM player WHERE id = $1", playerId)
	if err != nil {
		logs2.Logger(ctx).WithError(err).WithField("player_id", playerId).Errorf("Error deleting player by id, %v", err)
		return err
	}
	return nil
}

func (r PlayerRepository) DeleteTx(ctx context.Context, tx pgx.Tx, playerId string) error {
	_, err := tx.Exec(ctx, "DELETE FROM player WHERE id = $1", playerId)
	if err != nil {
		logs2.Logger(ctx).WithError(err).WithField("player_id", playerId).Errorf("Error deleting player by id, %v", err)
		return err
	}
	return nil
}

func (r PlayerRepository) GetById(ctx context.Context, playerId string) (*playerModels.Player, error) {
	result := &playerModels.Player{}
	if err := r.DbPool.QueryRow(ctx, "SELECT id, name, image, created_at, updated_at FROM player WHERE id=$1", playerId).
		Scan(&result.ID, &result.Name, &result.Image, &result.CreatedAt, &result.UpdatedAt); err != nil {
		logs2.Logger(ctx).WithError(err).WithField("player_id", playerId).Errorf("Error getting player by id, %v", err)
		return nil, err
	}
	return result, nil
}

func (r PlayerRepository) GetPaged(ctx context.Context, pagination *models.PaginationParameters, filter *playerModels.PlayerFilter) (*models.PaginationResponse, error) {
	span, ctx := r.Tracer.SpanFromContext(ctx)
	defer span.Finish()

	players := make([]*playerModels.Player, 0)
	var count int
	if err := r.DbPool.QueryRow(ctx, `SELECT count(*) 
FROM player WHERE ($1 = '' OR name like '%$1%') AND ($2 = '' OR $2 = id)
`, filter.Name).Scan(&count); err != nil {
		logs2.Logger(ctx).WithError(err).Errorf("Error count on player paged query, %v", err)
		span.Tag("error", err.Error())
		span.Tag("errorStack", string(debug.Stack()))
		return nil, err
	}

	rows, err := r.DbPool.Query(ctx, `SELECT id, name, image, created_at, updated_at 
FROM player 
WHERE ($1 = '' OR name like '%$1%') AND ($2 = '' OR $2 = id)
ORDER BY id
OFFSET $3
LIMIT $4`, filter.Name, filter.ID, (pagination.CurrentPage-1)*pagination.ItemsByPage, pagination.ItemsByPage)
	if err != nil {
		logs2.Logger(ctx).WithError(err).Errorf("Error player paged query, %v", err)
		return nil, err
	}

	for rows.Next() {
		p := &playerModels.Player{}
		rows.Scan(&p.ID, &p.Name, &p.Image, &p.CreatedAt, &p.UpdatedAt)
		players = append(players, p)
	}

	result := &models.PaginationResponse{
		Data:       players,
		TotalItems: count,
	}

	return result, nil
}
