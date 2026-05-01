package ranking

import (
	"strings"
	"testing"

	"github.com/josimar/ranking/backend/pkg/apperror"
	"github.com/josimar/ranking/backend/pkg/uuid"
	"github.com/stretchr/testify/require"
)

const (
	testAttrGraphics = "Graphics"
	testAttrSound    = "Sound"
	testTagGames     = "games"
)

// --- Visibility tests ---

func TestValidateVisibility_Public(t *testing.T) {
	t.Parallel()
	require.NoError(t, ValidateVisibility(Public))
}

func TestValidateVisibility_Private(t *testing.T) {
	t.Parallel()
	require.NoError(t, ValidateVisibility(Private))
}

func TestValidateVisibility_Invalid(t *testing.T) {
	t.Parallel()
	err := ValidateVisibility(Visibility("unknown"))
	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ValidationError, appErr.Code)
}

// --- NormalizeName tests ---

func TestNormalizeName_Lowercase(t *testing.T) {
	t.Parallel()
	require.Equal(t, "video games", NormalizeName("Video Games"))
}

func TestNormalizeName_StripAccents(t *testing.T) {
	t.Parallel()
	require.Equal(t, "cafe", NormalizeName("Café"))
}

func TestNormalizeName_ComplexAccents(t *testing.T) {
	t.Parallel()
	require.Equal(t, "acoes e operacoes", NormalizeName("Ações e Operações"))
}

func TestNormalizeName_AlreadyNormalized(t *testing.T) {
	t.Parallel()
	require.Equal(t, "simple", NormalizeName("simple"))
}

// --- NewRanking constructor tests ---

func TestNewRanking_Success(t *testing.T) {
	t.Parallel()
	attrs := []AttributeInput{{Name: testAttrGraphics, Description: "Visual quality"}}
	r, err := NewRanking("Video Games", "A ranking", Public, []string{testTagGames}, attrs, "owner-1")
	require.NoError(t, err)
	require.True(t, uuid.IsValid(r.ID))
	require.Equal(t, "Video Games", r.Name)
	require.Equal(t, "video games", r.NameLower)
	require.Equal(t, "A ranking", r.Description)
	require.Equal(t, Public, r.Visibility)
	require.Equal(t, []string{testTagGames}, r.Tags)
	require.Equal(t, "owner-1", r.OwnerUserID)
	require.Len(t, r.Attributes, 1)
	require.True(t, uuid.IsValid(r.Attributes[0].ID))
	require.Equal(t, testAttrGraphics, r.Attributes[0].Name)
	require.Equal(t, "Visual quality", r.Attributes[0].Description)
	require.True(t, r.Attributes[0].Active)
	require.False(t, r.CreatedAt.IsZero())
	require.False(t, r.UpdatedAt.IsZero())
}

func TestNewRanking_EmptyName(t *testing.T) {
	t.Parallel()
	_, err := NewRanking("", "desc", Public, nil, []AttributeInput{{Name: "A"}}, "owner")
	require.Error(t, err)
	var appErr *apperror.AppError
	require.ErrorAs(t, err, &appErr)
	require.Equal(t, apperror.ValidationError, appErr.Code)
}

func TestNewRanking_NameTooLong(t *testing.T) {
	t.Parallel()
	name := strings.Repeat("a", 101)
	_, err := NewRanking(name, "", Public, nil, []AttributeInput{{Name: "A"}}, "owner")
	require.Error(t, err)
}

func TestNewRanking_DescriptionTooLong(t *testing.T) {
	t.Parallel()
	desc := strings.Repeat("a", 501)
	_, err := NewRanking("Name", desc, Public, nil, []AttributeInput{{Name: "A"}}, "owner")
	require.Error(t, err)
}

func TestNewRanking_InvalidVisibility(t *testing.T) {
	t.Parallel()
	_, err := NewRanking("Name", "", Visibility("bad"), nil, []AttributeInput{{Name: "A"}}, "owner")
	require.Error(t, err)
}

func TestNewRanking_TooManyTags(t *testing.T) {
	t.Parallel()
	tags := make([]string, 11)
	for i := range tags {
		tags[i] = "tag"
	}
	_, err := NewRanking("Name", "", Public, tags, []AttributeInput{{Name: "A"}}, "owner")
	require.Error(t, err)
}

func TestNewRanking_NoAttributes(t *testing.T) {
	t.Parallel()
	_, err := NewRanking("Name", "", Public, nil, nil, "owner")
	require.Error(t, err)
}

func TestNewRanking_TooManyAttributes(t *testing.T) {
	t.Parallel()
	attrs := make([]AttributeInput, 7)
	for i := range attrs {
		attrs[i] = AttributeInput{Name: "A"}
	}
	_, err := NewRanking("Name", "", Public, nil, attrs, "owner")
	require.Error(t, err)
}

func TestNewRanking_AttributeNameEmpty(t *testing.T) {
	t.Parallel()
	_, err := NewRanking("Name", "", Public, nil, []AttributeInput{{Name: ""}}, "owner")
	require.Error(t, err)
}

func TestNewRanking_AttributeNameTooLong(t *testing.T) {
	t.Parallel()
	name := strings.Repeat("a", 51)
	_, err := NewRanking("Name", "", Public, nil, []AttributeInput{{Name: name}}, "owner")
	require.Error(t, err)
}

func TestNewRanking_AttributeDescriptionTooLong(t *testing.T) {
	t.Parallel()
	desc := strings.Repeat("a", 201)
	_, err := NewRanking("Name", "", Public, nil, []AttributeInput{{Name: "A", Description: desc}}, "owner")
	require.Error(t, err)
}

func TestNewRanking_EmptyOwnerUserID(t *testing.T) {
	t.Parallel()
	_, err := NewRanking("Name", "", Public, nil, []AttributeInput{{Name: "A"}}, "")
	require.Error(t, err)
}

// --- Update tests ---

func TestUpdate_AddNewAttribute(t *testing.T) {
	t.Parallel()
	r := validRanking(t)
	originalAttrID := r.Attributes[0].ID

	err := r.Update("New Name", "New Desc", Private, []string{"new"}, []AttributeInput{
		{Name: testAttrGraphics},
		{Name: testAttrSound, Description: "Audio quality"},
	})
	require.NoError(t, err)
	require.Equal(t, "New Name", r.Name)
	require.Equal(t, "new name", r.NameLower)
	require.Equal(t, "New Desc", r.Description)
	require.Equal(t, Private, r.Visibility)
	require.Equal(t, []string{"new"}, r.Tags)
	require.Len(t, r.Attributes, 2)
	// Existing attribute kept with same ID
	require.Equal(t, originalAttrID, r.Attributes[0].ID)
	require.True(t, r.Attributes[0].Active)
	// New attribute has new ID
	require.True(t, uuid.IsValid(r.Attributes[1].ID))
	require.Equal(t, testAttrSound, r.Attributes[1].Name)
	require.True(t, r.Attributes[1].Active)
}

func TestUpdate_SoftDeleteAttribute(t *testing.T) {
	t.Parallel()
	attrs := []AttributeInput{
		{Name: testAttrGraphics},
		{Name: testAttrSound},
	}
	r, err := NewRanking("Test", "", Public, nil, attrs, "owner")
	require.NoError(t, err)
	require.Len(t, r.Attributes, 2)

	// Update removing "Sound"
	err = r.Update("Test", "", Public, nil, []AttributeInput{{Name: testAttrGraphics}})
	require.NoError(t, err)
	// Both attributes still exist
	require.Len(t, r.Attributes, 2)
	// Graphics is active
	require.True(t, r.Attributes[0].Active)
	// Sound is soft-deleted
	require.False(t, r.Attributes[1].Active)
}

func TestUpdate_ValidationErrors(t *testing.T) {
	t.Parallel()
	r := validRanking(t)

	err := r.Update("", "", Public, nil, []AttributeInput{{Name: "A"}})
	require.Error(t, err)
}

// --- IsOwner tests ---

func TestIsOwner_True(t *testing.T) {
	t.Parallel()
	r := validRanking(t)
	require.True(t, r.IsOwner("owner-1"))
}

func TestIsOwner_False(t *testing.T) {
	t.Parallel()
	r := validRanking(t)
	require.False(t, r.IsOwner("other"))
}

// --- ActiveAttributes tests ---

func TestActiveAttributes_FiltersInactive(t *testing.T) {
	t.Parallel()
	attrs := []AttributeInput{
		{Name: testAttrGraphics},
		{Name: testAttrSound},
	}
	r, err := NewRanking("Test", "", Public, nil, attrs, "owner")
	require.NoError(t, err)

	// Soft-delete Sound
	err = r.Update("Test", "", Public, nil, []AttributeInput{{Name: testAttrGraphics}})
	require.NoError(t, err)

	active := r.ActiveAttributes()
	require.Len(t, active, 1)
	require.Equal(t, testAttrGraphics, active[0].Name)
}

func TestActiveAttributes_AllActive(t *testing.T) {
	t.Parallel()
	r := validRanking(t)
	active := r.ActiveAttributes()
	require.Len(t, active, 1)
	require.Equal(t, testAttrGraphics, active[0].Name)
}

// --- Helper ---

func validRanking(t *testing.T) *Ranking {
	t.Helper()
	r, err := NewRanking("Video Games", "A ranking", Public, []string{testTagGames}, []AttributeInput{{Name: testAttrGraphics}}, "owner-1")
	require.NoError(t, err)
	return r
}
