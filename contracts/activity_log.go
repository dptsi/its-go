package contracts

import "github.com/dptsi/its-go/models"

type Activity interface {
	CauserUserId() string
	ImpersonatorUserId() *string
}

type ActivityLogBuilder interface {
	Causer(user *models.User) ActivityLogBuilder
	Build() Activity
}

type ActivityLogService interface {
	Log(activity Activity, action string)
	Builder() ActivityLogBuilder
}
