package activitylog

import (
	"context"
	"fmt"

	"github.com/dptsi/its-go/contracts"
	"github.com/dptsi/its-go/models"
)

type Activity struct {
	causerUserId       string
	impersonatorUserId *string
}

func (a *Activity) CauserUserId() string {
	return a.causerUserId
}

func (a *Activity) ImpersonatorUserId() *string {
	return a.impersonatorUserId
}

type Builder struct {
	user *models.User
}

func (b *Builder) Causer(user *models.User) contracts.ActivityLogBuilder {
	b.user = user
	return b
}

func (b *Builder) Build() contracts.Activity {
	return &Activity{
		causerUserId:       b.user.Id(),
		impersonatorUserId: b.user.ImpersonatorId(),
	}

}

type Service struct {
	logger  contracts.LoggingService
	builder *Builder
}

func NewService(logger contracts.LoggingService) *Service {
	return &Service{
		logger:  logger,
		builder: &Builder{},
	}
}

func (s *Service) Builder() contracts.ActivityLogBuilder {
	return s.builder
}

func (s *Service) Log(ctx context.Context, activity contracts.Activity, action string) error {
	if activity.ImpersonatorUserId() != nil {
		return s.logger.Info(ctx, fmt.Sprintf("%s impersonated by %s executed action %s", activity.CauserUserId(), *activity.ImpersonatorUserId(), action))
	}

	return s.logger.Info(ctx, fmt.Sprintf("%s executed action %s", activity.CauserUserId(), action))
}
