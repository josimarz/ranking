package item

import (
	"testing"

	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/josimar/ranking/backend/pkg/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewItem_Success(t *testing.T) {
	t.Parallel()

	item, err := NewItem("PlayStation 5", "ranking-123", "user-456")

	require.NoError(t, err)
	require.NotNil(t, item)
	require.True(t, uuid.IsValid(item.ID))
	require.Equal(t, "PlayStation 5", item.Name)
	require.Equal(t, "", item.ImageKey)
	require.Equal(t, "ranking-123", item.RankingID)
	require.Equal(t, "user-456", item.CreatedBy)
	require.False(t, item.CreatedAt.IsZero())
	require.False(t, item.UpdatedAt.IsZero())
	require.Equal(t, item.CreatedAt, item.UpdatedAt)
}

func TestNewItem_EmptyName(t *testing.T) {
	t.Parallel()

	item, err := NewItem("", "ranking-123", "user-456")

	require.Nil(t, item)
	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ValidationError, appErr.Code)
}

func TestNewItem_NameTooLong(t *testing.T) {
	t.Parallel()

	longName := make([]byte, 101)
	for i := range longName {
		longName[i] = 'a'
	}

	item, err := NewItem(string(longName), "ranking-123", "user-456")

	require.Nil(t, item)
	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ValidationError, appErr.Code)
}

func TestNewItem_NameExactly100Chars(t *testing.T) {
	t.Parallel()

	name := make([]byte, 100)
	for i := range name {
		name[i] = 'a'
	}

	item, err := NewItem(string(name), "ranking-123", "user-456")

	require.NoError(t, err)
	require.NotNil(t, item)
	require.Len(t, item.Name, 100)
}

func TestItem_Update_Success(t *testing.T) {
	t.Parallel()

	item, err := NewItem("Old Name", "ranking-123", "user-456")
	require.NoError(t, err)

	originalUpdatedAt := item.UpdatedAt

	err = item.Update("New Name")

	require.NoError(t, err)
	require.Equal(t, "New Name", item.Name)
	require.True(t, item.UpdatedAt.After(originalUpdatedAt) || item.UpdatedAt.Equal(originalUpdatedAt))
}

func TestItem_Update_EmptyName(t *testing.T) {
	t.Parallel()

	item, err := NewItem("Old Name", "ranking-123", "user-456")
	require.NoError(t, err)

	err = item.Update("")

	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ValidationError, appErr.Code)
	require.Equal(t, "Old Name", item.Name)
}

func TestItem_Update_NameTooLong(t *testing.T) {
	t.Parallel()

	item, err := NewItem("Old Name", "ranking-123", "user-456")
	require.NoError(t, err)

	longName := make([]byte, 101)
	for i := range longName {
		longName[i] = 'a'
	}

	err = item.Update(string(longName))

	require.Error(t, err)
	require.Equal(t, "Old Name", item.Name)
}

func TestItem_CanBeModifiedBy_Creator(t *testing.T) {
	t.Parallel()

	item, err := NewItem("Test", "ranking-123", "creator-id")
	require.NoError(t, err)

	require.True(t, item.CanBeModifiedBy("creator-id", "other-owner"))
}

func TestItem_CanBeModifiedBy_RankingOwner(t *testing.T) {
	t.Parallel()

	item, err := NewItem("Test", "ranking-123", "creator-id")
	require.NoError(t, err)

	require.True(t, item.CanBeModifiedBy("some-user", "some-user"))
}

func TestItem_CanBeModifiedBy_UnauthorizedUser(t *testing.T) {
	t.Parallel()

	item, err := NewItem("Test", "ranking-123", "creator-id")
	require.NoError(t, err)

	require.False(t, item.CanBeModifiedBy("random-user", "owner-id"))
}
