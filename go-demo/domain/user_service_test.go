package domain_test

import (
	"api-demo/domain"
	mockDomain "api-demo/mock/domain"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewUserServiceImpl(t *testing.T) {
	_, err := domain.NewUserServiceImpl(nil)
	assert.Equal(t, fmt.Errorf("cannot create service, missing repo"), err, "expected error bad user id")
}

func TestUserServiceImpl_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)

	userRepo := mockDomain.NewMockUserRepo(ctrl)

	user := domain.NewUser("Shashank Pachava", "admin")

	userRepo.EXPECT().SaveUser(gomock.Any(), gomock.Any()).Return(nil)

	userService, err := domain.NewUserServiceImpl(userRepo)
	assert.NoError(t, err, "expected no error")

	savedUser, err := userService.CreateUser(context.Background(), user)
	assert.NoError(t, err, "expected no error")

	assert.Equal(t, savedUser, user, "expected saved user to be the same")
}

func TestUserServiceImpl_GetUserById(t *testing.T) {
	ctrl := gomock.NewController(t)

	userRepo := mockDomain.NewMockUserRepo(ctrl)

	user := domain.NewUser("Shashank Pachava", "admin")

	userRepo.EXPECT().GetUserById(context.Background(), user.Id).Return(user, nil)

	userService, err := domain.NewUserServiceImpl(userRepo)
	assert.NoError(t, err, "expected no error")

	savedUser, err := userService.GetUserById(context.Background(), user.Id)
	assert.NoError(t, err, "expected no error")

	assert.Equal(t, savedUser, user, "expected saved user to be the same")
}

func TestUserServiceImpl_GetUserById_BadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)

	userRepo := mockDomain.NewMockUserRepo(ctrl)

	userService, err := domain.NewUserServiceImpl(userRepo)
	assert.NoError(t, err, "expected no error")

	_, err = userService.GetUserById(context.Background(), uuid.Nil)
	assert.Equal(t, domain.ErrBadUserId, err, "expected error")
}

func TestUserServiceImpl_UpdateUser_BadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)

	userRepo := mockDomain.NewMockUserRepo(ctrl)

	userService, err := domain.NewUserServiceImpl(userRepo)
	assert.NoError(t, err, "expected no error")

	_, err = userService.UpdateUser(context.Background(), uuid.Nil, domain.UpdateUser{})
	assert.Equal(t, domain.ErrBadUserId, err, "expected error")
}

func TestUserServiceImpl_UpdateUser(t *testing.T) {
	t.Run("Update name", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userRepo := mockDomain.NewMockUserRepo(ctrl)

		user := domain.NewUser("Shashank Pachava", "admin")
		modifiedUser := domain.User{Id: user.Id, Name: "Shank", Role: user.Role}

		userRepo.EXPECT().GetUserById(context.Background(), user.Id).Return(user, nil)
		userRepo.EXPECT().
			SaveUser(
				context.Background(),
				modifiedUser,
			).
			Return(nil)

		userService, err := domain.NewUserServiceImpl(userRepo)
		assert.NoError(t, err, "expected no error")

		savedModifiedUser, err := userService.UpdateUser(context.Background(), user.Id, domain.UpdateUser{
			Name: &modifiedUser.Name,
		})
		assert.Equal(t, modifiedUser, savedModifiedUser, "different user found than expected")
	})

	t.Run("Update role", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userRepo := mockDomain.NewMockUserRepo(ctrl)

		user := domain.NewUser("Shashank Pachava", "admin")
		modifiedUser := domain.User{Id: user.Id, Name: user.Name, Role: "role"}

		userRepo.EXPECT().GetUserById(context.Background(), user.Id).Return(user, nil)
		userRepo.EXPECT().
			SaveUser(
				context.Background(),
				modifiedUser,
			).
			Return(nil)

		userService, err := domain.NewUserServiceImpl(userRepo)
		assert.NoError(t, err, "expected no error")

		savedModifiedUser, err := userService.UpdateUser(context.Background(), user.Id, domain.UpdateUser{
			Role: &modifiedUser.Role,
		})
		assert.Equal(t, modifiedUser, savedModifiedUser, "different user found than expected")
	})

	t.Run("Update name and role", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userRepo := mockDomain.NewMockUserRepo(ctrl)

		user := domain.NewUser("Shashank Pachava", "admin")
		modifiedUser := domain.User{Id: user.Id, Name: "Shank", Role: "role"}

		userRepo.EXPECT().GetUserById(context.Background(), user.Id).Return(user, nil)
		userRepo.EXPECT().
			SaveUser(
				context.Background(),
				modifiedUser,
			).
			Return(nil)

		userService, err := domain.NewUserServiceImpl(userRepo)
		assert.NoError(t, err, "expected no error")

		savedModifiedUser, err := userService.UpdateUser(context.Background(), user.Id, domain.UpdateUser{
			Name: &modifiedUser.Name,
			Role: &modifiedUser.Role,
		})
		assert.Equal(t, modifiedUser, savedModifiedUser, "different user found than expected")
	})

	t.Run("Update nothing", func(t *testing.T) {
		ctrl := gomock.NewController(t)

		userRepo := mockDomain.NewMockUserRepo(ctrl)

		user := domain.NewUser("Shashank Pachava", "admin")
		modifiedUser := domain.User{Id: user.Id, Name: user.Name, Role: user.Role}

		userRepo.EXPECT().GetUserById(context.Background(), user.Id).Return(user, nil)
		userRepo.EXPECT().
			SaveUser(
				context.Background(),
				modifiedUser,
			).
			Return(nil)

		userService, err := domain.NewUserServiceImpl(userRepo)
		assert.NoError(t, err, "expected no error")

		savedModifiedUser, err := userService.UpdateUser(context.Background(), user.Id, domain.UpdateUser{})
		assert.Equal(t, modifiedUser, savedModifiedUser, "different user found than expected")
	})
}

func TestUserServiceImpl_GetByProperty(t *testing.T) {
	ctrl := gomock.NewController(t)

	userRepo := mockDomain.NewMockUserRepo(ctrl)

	user1 := domain.NewUser("Shashank Pachava", "admin")
	user2 := domain.NewUser("Shank", "user")
	user3 := domain.NewUser("Sasi", "user")
	user4 := domain.NewUser("Sridhar", "user")

	userRepo.EXPECT().ListUsers(context.Background()).Return([]domain.User{
		user1,
		user2,
		user3,
		user4,
	}, nil)

	userService, err := domain.NewUserServiceImpl(userRepo)
	assert.NoError(t, err, "expected no error")

	role := "admin"

	filteredUsers, err := userService.GetByProperty(context.Background(), &domain.UserProperties{
		Role: &role,
	})
	assert.Equal(t, []domain.User{
		user1,
	}, filteredUsers, "different user found than expected")
}

func TestUserServiceImpl_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)

	userRepo := mockDomain.NewMockUserRepo(ctrl)

	user := domain.NewUser("Shashank Pachava", "admin")

	userRepo.EXPECT().DeleteUser(context.Background(), user.Id).Return(nil)

	userService, err := domain.NewUserServiceImpl(userRepo)
	assert.NoError(t, err, "expected no error")

	err = userService.Delete(context.Background(), user.Id)
	assert.NoError(t, err, "expected no error")
}

func TestUserServiceImpl_Delete_BadUserId(t *testing.T) {
	ctrl := gomock.NewController(t)

	userRepo := mockDomain.NewMockUserRepo(ctrl)

	userService, err := domain.NewUserServiceImpl(userRepo)
	assert.NoError(t, err, "expected no error")

	err = userService.Delete(context.Background(), uuid.Nil)
	assert.Equal(t, domain.ErrBadUserId, err, "expected error")
}
