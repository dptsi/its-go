package models

import "errors"

var (
	ErrUserDoesNotHaveRole = errors.New("user_does_not_have_role")
)

type Role struct {
	Id          string   `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
	IsDefault   bool     `json:"is_default"`
}

type User struct {
	id                 string
	name               string
	preferred_username string
	email              string
	picture            string
	activeRole         string
	activeRoleName     string
	roles              []Role
	hashedPassword     string
}

func NewUser(id string) *User {
	return &User{
		id:             id,
		activeRole:     "",
		activeRoleName: "",
		roles:          []Role{},
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) SetName(name string) {
	u.name = name
}

func (u *User) PreferredUsername() string {
	return u.preferred_username
}

func (u *User) SetPreferredUsername(preferredUsername string) {
	u.preferred_username = preferredUsername
}

func (u *User) Email() string {
	return u.email
}

func (u *User) SetEmail(email string) {
	u.email = email
}

func (u *User) Picture() string {
	return u.picture
}

func (u *User) SetPicture(picture string) {
	u.picture = picture
}

func (u *User) ActiveRole() string {
	return u.activeRole
}

func (u *User) ActiveRoleName() string {
	return u.activeRoleName
}

func (u *User) Roles() []Role {
	return u.roles
}

func (u *User) AddRole(id string, name string, permissions []string, isDefault bool) {
	u.roles = append(u.roles, Role{
		Id:          id,
		Name:        name,
		Permissions: permissions,
		IsDefault:   isDefault,
	})

	if isDefault || u.activeRole == "" {
		u.SetActiveRole(id)
	}
}

func (u *User) SetActiveRole(id string) error {
	for _, role := range u.roles {
		if role.Id == id {
			u.activeRole = id
			u.activeRoleName = role.Name
			return nil
		}
	}

	return ErrUserDoesNotHaveRole
}

func (u *User) HasPermission(permission string) bool {
	for _, role := range u.roles {
		if role.Id == u.activeRole {
			for _, perm := range role.Permissions {
				if perm == permission {
					return true
				}
			}
		}
	}

	return false
}

func (u *User) HashedPassword() string {
	return u.hashedPassword
}

func (u *User) SetHashedPassword(hashedPassword string) {
	u.hashedPassword = hashedPassword
}
