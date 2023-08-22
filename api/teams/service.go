package teams_api

import (
	"context"

	"github.com/akshitbansal-1/async-testing/be/app"
)

type Service interface {
	GetTeams(ctx context.Context) ([]Team, error)
	AddMember(ctx context.Context, teamId string, userEmail string) error
}

type service struct {
	app app.App
}

func NewService(app app.App) Service {
	return &service{app}
}

func (s *service) GetTeams(ctx context.Context) ([]Team, error) {
	return GetTeams(ctx, s.app)
}

func (s *service) AddMember(ctx context.Context, teamId string, userEmail string) error {
	return AddMember(ctx, s.app, teamId, userEmail)
}
