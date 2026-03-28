// @AI_GENERATED
package entity

import "time"

// Device is an entity representing a device associated with a principal.
type Device struct {
	id         string
	name       string
	platform   string
	lastSeenAt time.Time
}

// NewDevice creates a new Device with the given id, name, and platform.
// The lastSeenAt field is set to the current time.
func NewDevice(id, name, platform string) *Device {
	return &Device{
		id:         id,
		name:       name,
		platform:   platform,
		lastSeenAt: time.Now(),
	}
}

// ID returns the device identifier.
func (d *Device) ID() string { return d.id }

// Name returns the device name.
func (d *Device) Name() string { return d.name }

// Platform returns the device platform.
func (d *Device) Platform() string { return d.platform }

// LastSeenAt returns the time the device was last seen.
func (d *Device) LastSeenAt() time.Time { return d.lastSeenAt }

// @AI_GENERATED: end
