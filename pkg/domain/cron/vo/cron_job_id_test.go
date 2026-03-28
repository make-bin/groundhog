package vo

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCronJobID_EmptyString_ReturnsError(t *testing.T) {
	_, err := NewCronJobID("")
	assert.Error(t, err)
}

func TestNewCronJobID_InvalidFormat_ReturnsError(t *testing.T) {
	cases := []string{
		"not-a-uuid",
		"12345678-1234-1234-1234-12345678901",   // too short
		"12345678-1234-1234-1234-1234567890123", // too long
		"XXXXXXXX-1234-1234-1234-123456789012",  // uppercase / non-hex
	}
	for _, tc := range cases {
		_, err := NewCronJobID(tc)
		assert.Errorf(t, err, "expected error for input %q", tc)
	}
}

func TestNewCronJobID_ValidUUID_Success(t *testing.T) {
	id, err := NewCronJobID("550e8400-e29b-41d4-a716-446655440000")
	require.NoError(t, err)
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.Value())
}

func TestGenerateCronJobID_ProducesValidUUID(t *testing.T) {
	id := GenerateCronJobID()
	assert.NotEmpty(t, id.Value())

	// Must be re-parseable by NewCronJobID.
	_, err := NewCronJobID(id.Value())
	require.NoError(t, err)
}

func TestGenerateCronJobID_ProducesUniqueIDs(t *testing.T) {
	seen := make(map[string]struct{}, 100)
	for i := 0; i < 100; i++ {
		id := GenerateCronJobID()
		_, dup := seen[id.Value()]
		assert.False(t, dup, "duplicate ID generated: %s", id.Value())
		seen[id.Value()] = struct{}{}
	}
}

func TestGenerateCronJobID_HasCorrectFormat(t *testing.T) {
	id := GenerateCronJobID()
	parts := strings.Split(id.Value(), "-")
	require.Len(t, parts, 5)
	assert.Len(t, parts[0], 8)
	assert.Len(t, parts[1], 4)
	assert.Len(t, parts[2], 4)
	assert.Len(t, parts[3], 4)
	assert.Len(t, parts[4], 12)
}

func TestCronJobID_Equals(t *testing.T) {
	a, _ := NewCronJobID("550e8400-e29b-41d4-a716-446655440000")
	b, _ := NewCronJobID("550e8400-e29b-41d4-a716-446655440000")
	c, _ := NewCronJobID("550e8400-e29b-41d4-a716-446655440001")

	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}

func TestCronJobID_String(t *testing.T) {
	id, _ := NewCronJobID("550e8400-e29b-41d4-a716-446655440000")
	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", id.String())
}
