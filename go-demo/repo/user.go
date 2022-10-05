package repo

import (
	"api-demo/domain"
	"context"
	"errors"
	"github.com/google/uuid"
	"sync"
)

type InMemUserRepo struct {
	*sync.Map
}

func NewInMemUserRepo() InMemUserRepo {
	return InMemUserRepo{Map: new(sync.Map)}
}

func (i *InMemUserRepo) SaveUser(ctx context.Context, user domain.User) error {
	i.Store(user.Id, user)
	return nil
}

func (i *InMemUserRepo) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, found := i.LoadAndDelete(id)
	if !found {
		return domain.ErrUserIdNotFound{Id: id}
	}
	return nil
}

func (i *InMemUserRepo) GetUserById(ctx context.Context, id uuid.UUID) (domain.User, error) {
	val, ok := i.Load(id)
	if !ok {
		return domain.User{}, domain.ErrUserIdNotFound{Id: id}
	}
	user, ok := val.(domain.User)
	if !ok {
		return user, errors.New("could not cast map val to user type")
	}
	return user, nil
}

func (i *InMemUserRepo) ListUsers(ctx context.Context) ([]domain.User, error) {
	var resp []domain.User
	var err error
	i.Range(func(_, val any) bool {
		user, ok := val.(domain.User)
		if !ok {
			err = errors.New("could not cast map val to user type")
			return false
		}
		resp = append(resp, user)
		return true
	})
	return resp, err
}
