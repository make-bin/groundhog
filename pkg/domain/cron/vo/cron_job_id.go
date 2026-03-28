package vo

import (
	"crypto/rand"
	"fmt"
	"regexp"
)

// uuidRegexp validates UUID v4 format (lowercase hex with dashes).
var uuidRegexp = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

// CronJobID represents a unique identifier for a cron job.
// It is immutable after creation.
type CronJobID struct {
	value string
}

// NewCronJobID creates a new CronJobID after validating that the value is
// non-empty and conforms to UUID format.
func NewCronJobID(v string) (CronJobID, error) {
	if v == "" {
		return CronJobID{}, fmt.Errorf("cron job ID must not be empty")
	}
	if !uuidRegexp.MatchString(v) {
		return CronJobID{}, fmt.Errorf("cron job ID must be a valid UUID, got: %s", v)
	}
	return CronJobID{value: v}, nil
}

// GenerateCronJobID generates a new random UUID v4 and returns it as a CronJobID.
func GenerateCronJobID() CronJobID {
	var b [16]byte
	_, err := rand.Read(b[:])
	if err != nil {
		// crypto/rand failure is unrecoverable in practice; panic to surface it.
		panic(fmt.Sprintf("failed to generate cron job ID: %v", err))
	}
	// Set version 4 bits (bits 12-15 of byte 6).
	b[6] = (b[6] & 0x0f) | 0x40
	// Set variant bits (bits 6-7 of byte 8).
	b[8] = (b[8] & 0x3f) | 0x80

	id := fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
	return CronJobID{value: id}
}

// Value returns the cron job ID string.
func (c CronJobID) Value() string { return c.value }

// Equals returns true if c and other represent the same cron job ID.
func (c CronJobID) Equals(other CronJobID) bool {
	return c.value == other.value
}

// String returns the string representation of the cron job ID.
func (c CronJobID) String() string { return c.value }
