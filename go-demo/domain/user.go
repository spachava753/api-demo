package domain

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
)

type User struct {
	Id   uuid.UUID
	Name string
	Role string
}

type UserProperties struct {
	Name *string
	Role *string
}

type UpdateUser struct {
	Name *string
	Role *string
}

func NewUser(name string, role string) User {
	return User{Name: name, Role: role, Id: uuid.New()}
}

var ErrBadUserId = errors.New("invalid user id")

type UserService interface {
	CreateUser(ctx context.Context, user User) (User, error)
	GetUserById(ctx context.Context, id uuid.UUID) (User, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateUser(ctx context.Context, id uuid.UUID, updateUser UpdateUser) (User, error)
	GetByProperty(ctx context.Context, up *UserProperties) ([]User, error)
}

// UserServiceImpl is an implementation of UserService
type UserServiceImpl struct {
	repo UserRepo
}

// NewUserServiceImpl returns a new instance of UserServiceImpl
func NewUserServiceImpl(repo UserRepo) (UserServiceImpl, error) {
	if repo == nil {
		return UserServiceImpl{},
			fmt.Errorf("cannot create service, missing repo")
	}
	return UserServiceImpl{repo: repo}, nil
}

func (u *UserServiceImpl) CreateUser(ctx context.Context, user User) (User, error) {
	log.Println("creating user")

	if user.Id == uuid.Nil {
		return User{}, ErrBadUserId
	}

	err := u.repo.SaveUser(ctx, user)
	if err != nil {
		return User{}, fmt.Errorf("could not create user: %w", err)
	}

	return user, nil
}

func (u *UserServiceImpl) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	log.Println("fetching user by id")

	if id == uuid.Nil {
		return User{}, ErrBadUserId
	}

	user, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("could not fetch user by id: %w", err)
	}

	return user, nil
}

func (u *UserServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	log.Println("deleting user by id")

	if id == uuid.Nil {
		return ErrBadUserId
	}

	err := u.repo.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("could not delete user by id: %w", err)
	}

	return nil
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, id uuid.UUID, updateUser UpdateUser) (User, error) {
	log.Println("updating user")

	if id == uuid.Nil {
		return User{}, ErrBadUserId
	}

	user, err := u.repo.GetUserById(ctx, id)
	if err != nil {
		return User{}, fmt.Errorf("could fetch user by id: %w", err)
	}

	if updateUser.Name != nil {
		user.Name = *updateUser.Name
	}

	if updateUser.Role != nil {
		user.Role = *updateUser.Role
	}

	err = u.repo.SaveUser(ctx, user)
	if err != nil {
		return User{}, fmt.Errorf("could not update user: %w", err)
	}

	return user, nil
}

func (u *UserServiceImpl) GetByProperty(ctx context.Context, up *UserProperties) ([]User, error) {
	log.Printf("fetching users by property %#v", up)

	users, err := u.repo.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list stored users: %w", err)
	}

	if up != nil {
		users = filterUsers(users, *up)
	}

	return users, nil
}

// filterUsers filters users without preserving order
func filterUsers(users []User, properties UserProperties) []User {
	i := 0
	for i < len(users) {
		if matchesUser(users[i], properties) {
			i++
			continue
		}

		users[i] = users[len(users)-1]
		users = users[:len(users)-1]
	}

	return users
}

// matchesUser checks if a single user matches against UserProperties
func matchesUser(user User, properties UserProperties) bool {
	if properties.Role != nil && *properties.Role != user.Role {
		return false
	}

	if properties.Name != nil && *properties.Name != user.Name {
		return false
	}

	return true
}

type ErrUserIdNotFound struct {
	Id uuid.UUID
}

func (e ErrUserIdNotFound) Error() string {
	return fmt.Sprintf("could not find user with id %s", e.Id.String())
}

type UserRepo interface {
	SaveUser(ctx context.Context, user User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserById(ctx context.Context, id uuid.UUID) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
}
