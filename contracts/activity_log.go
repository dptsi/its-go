package contracts

import (
	"context"

	"github.com/dptsi/its-go/models"
)

type Activity interface {
	CauserUserId() string
	ImpersonatorUserId() *string
}

type ActivityLogBuilder interface {
	Causer(user *models.User) ActivityLogBuilder
	Build() Activity
}

type ActivityLogService interface {
	Log(ctx context.Context, activity Activity, action string) error
	Builder() ActivityLogBuilder
}
