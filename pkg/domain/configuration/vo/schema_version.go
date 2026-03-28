// @AI_GENERATED
package vo

import (
	"fmt"
)

// SchemaVersion represents a semantic version for configuration schema.
// It is immutable after creation.
type SchemaVersion struct {
	major int
	minor int
	patch int
}

// NewSchemaVersion creates a new SchemaVersion after validating that all components are non-negative.
func NewSchemaVersion(major, minor, patch int) (SchemaVersion, error) {
	if major < 0 {
		return SchemaVersion{}, fmt.Errorf("major version must be non-negative, got %d", major)
	}
	if minor < 0 {
		return SchemaVersion{}, fmt.Errorf("minor version must be non-negative, got %d", minor)
	}
	if patch < 0 {
		return SchemaVersion{}, fmt.Errorf("patch version must be non-negative, got %d", patch)
	}
	return SchemaVersion{major: major, minor: minor, patch: patch}, nil
}

// Major returns the major version component.
func (v SchemaVersion) Major() int { return v.major }

// Minor returns the minor version component.
func (v SchemaVersion) Minor() int { return v.minor }

// Patch returns the patch version component.
func (v SchemaVersion) Patch() int { return v.patch }

// LessThan returns true if v is less than other using semantic version ordering.
func (v SchemaVersion) LessThan(other SchemaVersion) bool {
	if v.major != other.major {
		return v.major < other.major
	}
	if v.minor != other.minor {
		return v.minor < other.minor
	}
	return v.patch < other.patch
}

// String returns the version in "major.minor.patch" format.
func (v SchemaVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
}

// Equals returns true if v and other represent the same version.
func (v SchemaVersion) Equals(other SchemaVersion) bool {
	return v.major == other.major && v.minor == other.minor && v.patch == other.patch
}

// @AI_GENERATED: end
