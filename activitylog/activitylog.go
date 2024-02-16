package activitylog

import (
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

func (s *Service) Log(activity contracts.Activity, action string) {
	if activity.ImpersonatorUserId() != nil {
		s.logger.Info("{{.Causer}} impersonated by {{.Impersonator}} executed action {{.Action}}", map[string]interface{}{
			"Action":       action,
			"Causer":       activity.CauserUserId(),
			"Impersonator": *activity.ImpersonatorUserId(),
		})
		return
	}

	s.logger.Info("{{.Causer}} executed action {{.Action}}", map[string]interface{}{
		"Causer": activity.CauserUserId(),
		"Action": action,
	})
}
