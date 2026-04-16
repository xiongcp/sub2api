//go:build unit

package service

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type userRepoStubForUpdateUser struct {
	user    *User
	updated *User
	updateErr error
}

func (s *userRepoStubForUpdateUser) Create(context.Context, *User) error { panic("unexpected Create call") }

func (s *userRepoStubForUpdateUser) GetByID(context.Context, int64) (*User, error) {
	if s.user == nil {
		return nil, ErrUserNotFound
	}
	return s.user, nil
}

func (s *userRepoStubForUpdateUser) GetByEmail(context.Context, string) (*User, error) {
	panic("unexpected GetByEmail call")
}

func (s *userRepoStubForUpdateUser) GetFirstAdmin(context.Context) (*User, error) {
	panic("unexpected GetFirstAdmin call")
}

func (s *userRepoStubForUpdateUser) Update(_ context.Context, user *User) error {
	if s.updateErr != nil {
		return s.updateErr
	}
	s.updated = user
	return nil
}

func (s *userRepoStubForUpdateUser) Delete(context.Context, int64) error { panic("unexpected Delete call") }

func (s *userRepoStubForUpdateUser) List(context.Context, pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected List call")
}

func (s *userRepoStubForUpdateUser) ListWithFilters(context.Context, pagination.PaginationParams, UserListFilters) ([]User, *pagination.PaginationResult, error) {
	panic("unexpected ListWithFilters call")
}

func (s *userRepoStubForUpdateUser) UpdateBalance(context.Context, int64, float64) error {
	panic("unexpected UpdateBalance call")
}

func (s *userRepoStubForUpdateUser) DeductBalance(context.Context, int64, float64) error {
	panic("unexpected DeductBalance call")
}

func (s *userRepoStubForUpdateUser) UpdateConcurrency(context.Context, int64, int) error {
	panic("unexpected UpdateConcurrency call")
}

func (s *userRepoStubForUpdateUser) ExistsByEmail(context.Context, string) (bool, error) {
	panic("unexpected ExistsByEmail call")
}

func (s *userRepoStubForUpdateUser) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error) {
	panic("unexpected RemoveGroupFromAllowedGroups call")
}

func (s *userRepoStubForUpdateUser) AddGroupToAllowedGroups(context.Context, int64, int64) error {
	panic("unexpected AddGroupToAllowedGroups call")
}

func (s *userRepoStubForUpdateUser) RemoveGroupFromUserAllowedGroups(context.Context, int64, int64) error {
	panic("unexpected RemoveGroupFromUserAllowedGroups call")
}

func (s *userRepoStubForUpdateUser) UpdateTotpSecret(context.Context, int64, *string) error {
	panic("unexpected UpdateTotpSecret call")
}

func (s *userRepoStubForUpdateUser) EnableTotp(context.Context, int64) error {
	panic("unexpected EnableTotp call")
}

func (s *userRepoStubForUpdateUser) DisableTotp(context.Context, int64) error {
	panic("unexpected DisableTotp call")
}

func TestAdminService_UpdateUser_WithPasswordIncrementsTokenVersion(t *testing.T) {
	repo := &userRepoStubForUpdateUser{
		user: &User{
			ID:           7,
			Email:        "user@test.com",
			PasswordHash: mustHashForAdminUpdateTest(t, "old-password"),
			Status:       StatusActive,
			Role:         RoleUser,
			TokenVersion: 3,
		},
	}
	svc := &adminServiceImpl{userRepo: repo}

	user, err := svc.UpdateUser(context.Background(), 7, &UpdateUserInput{
		Password: "new-password",
	})
	require.NoError(t, err)
	require.NotNil(t, user)
	require.NotNil(t, repo.updated)
	require.True(t, user.CheckPassword("new-password"))
	require.False(t, user.CheckPassword("old-password"))
	require.Equal(t, int64(4), user.TokenVersion)
	require.Equal(t, int64(4), repo.updated.TokenVersion)
}

func TestAdminService_UpdateUser_WithoutPasswordKeepsTokenVersion(t *testing.T) {
	repo := &userRepoStubForUpdateUser{
		user: &User{
			ID:           8,
			Email:        "user@test.com",
			PasswordHash: mustHashForAdminUpdateTest(t, "same-password"),
			Status:       StatusActive,
			Role:         RoleUser,
			TokenVersion: 5,
		},
	}
	svc := &adminServiceImpl{userRepo: repo}

	username := "updated-name"
	user, err := svc.UpdateUser(context.Background(), 8, &UpdateUserInput{
		Username: &username,
	})
	require.NoError(t, err)
	require.NotNil(t, user)
	require.Equal(t, int64(5), user.TokenVersion)
	require.True(t, user.CheckPassword("same-password"))
}

func mustHashForAdminUpdateTest(t *testing.T, password string) string {
	t.Helper()

	user := &User{}
	require.NoError(t, user.SetPassword(password))
	return user.PasswordHash
}
