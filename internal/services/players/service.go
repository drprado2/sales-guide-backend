package players

import (
	"context"
	"github.com/drprado2/react-redux-typescript/internal/apperrors"
	"github.com/drprado2/react-redux-typescript/internal/apptracer"
	"github.com/drprado2/react-redux-typescript/internal/logs"
	"github.com/drprado2/react-redux-typescript/internal/models"
	playerModels "github.com/drprado2/react-redux-typescript/internal/models/players"
	"github.com/drprado2/react-redux-typescript/internal/storage/repositories"
)

type UserServiceInterface interface{
	Create(ctx context.Context, request *playerModels.CreatePlayerRequest) (*playerModels.Player, error)
	Update(ctx context.Context, playerId string, request *playerModels.UpdatePlayerRequest) (*playerModels.Player, error)
	Delete(ctx context.Context, playerId string) error
	GetById(ctx context.Context, playerId string) (*playerModels.Player, error)
	GetPaged(ctx context.Context, pageParams *models.PaginationParameters, filter *playerModels.PlayerFilter) (*models.PaginationResponse, error)
}

type UserService struct {
	PlayerRepository repositories.PlayerRepositoryInterface
	Tracer apptracer.TracerService
}

func (u *UserService) Create(ctx context.Context, request *playerModels.CreatePlayerRequest) (*playerModels.Player, error) {
	if len(request.Name) == 0 || len(request.Image) == 0 {
		logs.Logger(ctx).Warn("create user with name or image empty")
		return nil, apperrors.InvalidRequestParameters
	}

	return u.PlayerRepository.Create(ctx, request)
}

func (u *UserService) Update(ctx context.Context, playerId string, request *playerModels.UpdatePlayerRequest) (*playerModels.Player, error){
	if len(playerId) == 0 || len(request.Name) == 0 || len(request.Image) == 0 {
		logs.Logger(ctx).Warn("update user with id or name or image empty")
		return nil, apperrors.InvalidRequestParameters
	}

	return u.PlayerRepository.Update(ctx, playerId, request)
}

func (u *UserService) Delete(ctx context.Context, playerId string) error {
	if len(playerId) == 0 {
		logs.Logger(ctx).Warn("delete user with id empty")
		return apperrors.InvalidRequestParameters
	}

	return u.PlayerRepository.Delete(ctx, playerId)
}

func (u *UserService) GetById(ctx context.Context, playerId string) (*playerModels.Player, error) {
	if len(playerId) == 0 {
		logs.Logger(ctx).Warn("get user with id empty")
		return nil, apperrors.InvalidRequestParameters
	}

	return u.PlayerRepository.GetById(ctx, playerId)
}

func (u *UserService) GetPaged(ctx context.Context, pageParams *models.PaginationParameters, filter *playerModels.PlayerFilter) (*models.PaginationResponse, error) {
	span, ctx := u.Tracer.SpanFromContext(ctx)
	defer span.Finish()

	if pageParams.ItemsByPage < 1 || pageParams.CurrentPage < 1 {
		logs.Logger(ctx).Warn("get user with id empty")
		return nil, apperrors.InvalidRequestParameters
	}

	return u.PlayerRepository.GetPaged(ctx, pageParams, filter)
}
