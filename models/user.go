package models

import "errors"

var (
	ErrUserDoesNotHaveRole = errors.New("user_does_not_have_role")
)

type UnitOrganization struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	NameEn string `json:"name_en"`
}

type Role struct {
	Id                string             `json:"id"`
	Name              string             `json:"name"`
	Permissions       []string           `json:"permissions"`
	UnitOrganizations []UnitOrganization `json:"unit_organizations"`
	IsDefault         bool               `json:"is_default"`
}

type User struct {
	id                 string
	regId              string
	name               string
	preferred_username string
	email              string
	phone              string
	picture            string
	roles              []Role
	hashedPassword     string
	impersonatorId     *string
}

func NewUser(id string) *User {
	return &User{
		id:    id,
		roles: []Role{},
	}
}

func (u *User) Id() string {
	return u.id
}

func (u *User) RegId() string {
	return u.regId
}

func (u *User) SetRegId(regId string) {
	u.regId = regId
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

func (u *User) Phone() string {
	return u.phone
}

func (u *User) SetPhone(phone string) {
	u.phone = phone
}

func (u *User) Picture() string {
	return u.picture
}

func (u *User) SetPicture(picture string) {
	u.picture = picture
}

func (u *User) Roles() []Role {
	return u.roles
}

func (u *User) SetRoles(roles []Role) {
	u.roles = make([]Role, 0)
	var isRoleDuplicate = make(map[string]bool)
	for _, role := range roles {
		if isRoleDuplicate[role.Id] {
			continue
		}
		u.roles = append(u.roles, role)
		isRoleDuplicate[role.Id] = true
	}
}

func (u *User) AddRole(id string, name string, permissions []string, isDefault bool, unitOrganizations []UnitOrganization) {
	u.roles = append(u.roles, Role{
		Id:                id,
		Name:              name,
		Permissions:       permissions,
		IsDefault:         isDefault,
		UnitOrganizations: unitOrganizations,
	})
}

func (u *User) RoleUnitOrganizations(roleId string) []UnitOrganization {
	for _, role := range u.roles {
		if role.Id == roleId {
			return role.UnitOrganizations
		}
	}

	return make([]UnitOrganization, 0)
}

func (u *User) HasUnitOrganizationPermission(unitOrgId, permission string) bool {
	for _, role := range u.roles {
		for _, unitOrg := range role.UnitOrganizations {
			if unitOrg.Id == unitOrgId {
				for _, perm := range role.Permissions {
					if perm == permission {
						return true
					}
				}
			}
		}
	}

	return false
}

func (u *User) HasPermission(permission string) bool {
	for _, role := range u.roles {
		for _, perm := range role.Permissions {
			if perm == permission {
				return true
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

func (u *User) ImpersonatorId() *string {
	return u.impersonatorId
}

func (u *User) SetImpersonatorId(impersonatorId *string) {
	u.impersonatorId = impersonatorId
}
