package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// --- NewCronScheduleAt ---

func TestNewCronScheduleAt_Valid(t *testing.T) {
	s, err := NewCronScheduleAt("2025-01-01T00:00:00Z")
	require.NoError(t, err)
	assert.Equal(t, ScheduleKindAt, s.Kind())
	assert.Equal(t, "2025-01-01T00:00:00Z", s.At())
}

func TestNewCronScheduleAt_WithOffset(t *testing.T) {
	s, err := NewCronScheduleAt("2025-06-15T12:30:00+08:00")
	require.NoError(t, err)
	assert.Equal(t, ScheduleKindAt, s.Kind())
}

func TestNewCronScheduleAt_EmptyString(t *testing.T) {
	_, err := NewCronScheduleAt("")
	assert.Error(t, err)
}

func TestNewCronScheduleAt_InvalidFormat(t *testing.T) {
	_, err := NewCronScheduleAt("not-a-date")
	assert.Error(t, err)
}

func TestNewCronScheduleAt_DateOnlyNotValid(t *testing.T) {
	// RFC3339 requires time component; date-only strings should fail
	_, err := NewCronScheduleAt("2025-01-01")
	assert.Error(t, err)
}

// --- NewCronScheduleEvery ---

func TestNewCronScheduleEvery_Valid(t *testing.T) {
	s, err := NewCronScheduleEvery(60000, nil)
	require.NoError(t, err)
	assert.Equal(t, ScheduleKindEvery, s.Kind())
	assert.Equal(t, int64(60000), s.EveryMs())
	assert.Nil(t, s.AnchorMs())
}

func TestNewCronScheduleEvery_WithAnchor(t *testing.T) {
	anchor := int64(1700000000000)
	s, err := NewCronScheduleEvery(5000, &anchor)
	require.NoError(t, err)
	assert.Equal(t, &anchor, s.AnchorMs())
}

func TestNewCronScheduleEvery_ZeroInterval(t *testing.T) {
	_, err := NewCronScheduleEvery(0, nil)
	assert.Error(t, err)
}

func TestNewCronScheduleEvery_NegativeInterval(t *testing.T) {
	_, err := NewCronScheduleEvery(-1000, nil)
	assert.Error(t, err)
}

// --- NewCronScheduleCron ---

func TestNewCronScheduleCron_FiveField(t *testing.T) {
	s, err := NewCronScheduleCron("0 * * * *", "", 0)
	require.NoError(t, err)
	assert.Equal(t, ScheduleKindCron, s.Kind())
	assert.Equal(t, "0 * * * *", s.Expr())
	assert.Equal(t, "", s.Tz())
	assert.Equal(t, int64(0), s.StaggerMs())
}

func TestNewCronScheduleCron_SixField(t *testing.T) {
	s, err := NewCronScheduleCron("0 0 * * * *", "UTC", 0)
	require.NoError(t, err)
	assert.Equal(t, ScheduleKindCron, s.Kind())
}

func TestNewCronScheduleCron_Descriptor(t *testing.T) {
	s, err := NewCronScheduleCron("@hourly", "", 0)
	require.NoError(t, err)
	assert.Equal(t, "@hourly", s.Expr())
}

func TestNewCronScheduleCron_WithTimezone(t *testing.T) {
	s, err := NewCronScheduleCron("0 9 * * *", "Asia/Shanghai", 0)
	require.NoError(t, err)
	assert.Equal(t, "Asia/Shanghai", s.Tz())
}

func TestNewCronScheduleCron_WithStagger(t *testing.T) {
	s, err := NewCronScheduleCron("0 * * * *", "", 30000)
	require.NoError(t, err)
	assert.Equal(t, int64(30000), s.StaggerMs())
}

func TestNewCronScheduleCron_EmptyExpr(t *testing.T) {
	_, err := NewCronScheduleCron("", "", 0)
	assert.Error(t, err)
}

func TestNewCronScheduleCron_InvalidExpr(t *testing.T) {
	_, err := NewCronScheduleCron("not a cron", "", 0)
	assert.Error(t, err)
}

func TestNewCronScheduleCron_InvalidTimezone(t *testing.T) {
	_, err := NewCronScheduleCron("0 * * * *", "Invalid/Zone", 0)
	assert.Error(t, err)
}

func TestNewCronScheduleCron_NegativeStagger(t *testing.T) {
	_, err := NewCronScheduleCron("0 * * * *", "", -1)
	assert.Error(t, err)
}
