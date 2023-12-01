package adapters

import "bitbucket.org/dptsi/base-go-libraries/contracts"

type SliceUser struct {
	Id             string
	Username       string
	HashedPassword string
}

type SliceUserRepository struct {
	Users []SliceUser
}

func NewSliceUserRepository(users []SliceUser) *SliceUserRepository {
	return &SliceUserRepository{
		Users: users,
	}
}

func (r SliceUserRepository) FindByUsername(username string) (*contracts.User, error) {
	for _, user := range r.Users {
		if user.Username == username {
			u := contracts.NewUser(user.Id)
			u.SetHashedPassword(user.HashedPassword)
			return u, nil
		}
	}

	return nil, nil
}
