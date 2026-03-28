// @AI_GENERATED
package principal

import (
	"fmt"
	"time"

	"github.com/make-bin/groundhog/pkg/domain/identity/entity"
	"github.com/make-bin/groundhog/pkg/domain/identity/vo"
)

// PrincipalType represents the type of a principal.
type PrincipalType int

const (
	PrincipalUser PrincipalType = iota
	PrincipalDevice
	PrincipalService
	PrincipalAPIKey
)

// Principal is the aggregate root for the Identity domain.
type Principal struct {
	id            vo.PrincipalID
	principalType PrincipalType
	credentials   []vo.Credential
	devices       []entity.Device
	auditLog      []entity.AuditEntry
	rateLimit     vo.RateLimitPolicy
	createdAt     time.Time
	lastSeenAt    time.Time
}

// NewPrincipal creates a new Principal with the given id and type.
func NewPrincipal(id vo.PrincipalID, pType PrincipalType) *Principal {
	now := time.Now()
	return &Principal{
		id:            id,
		principalType: pType,
		credentials:   []vo.Credential{},
		devices:       []entity.Device{},
		auditLog:      []entity.AuditEntry{},
		createdAt:     now,
		lastSeenAt:    now,
	}
}

// ID returns the principal's identifier.
func (p *Principal) ID() vo.PrincipalID { return p.id }

// Type returns the principal's type.
func (p *Principal) Type() PrincipalType { return p.principalType }

// Devices returns a copy of the principal's devices.
func (p *Principal) Devices() []entity.Device {
	result := make([]entity.Device, len(p.devices))
	copy(result, p.devices)
	return result
}

// RateLimit returns the principal's rate limit policy.
func (p *Principal) RateLimit() vo.RateLimitPolicy { return p.rateLimit }

// CreatedAt returns the time the principal was created.
func (p *Principal) CreatedAt() time.Time { return p.createdAt }

// LastSeenAt returns the time the principal was last seen.
func (p *Principal) LastSeenAt() time.Time { return p.lastSeenAt }

// AddDevice adds a device to the principal. Returns an error if a device
// with the same ID already exists.
func (p *Principal) AddDevice(device entity.Device) error {
	for _, d := range p.devices {
		if d.ID() == device.ID() {
			return fmt.Errorf("device with id %q already exists", device.ID())
		}
	}
	p.devices = append(p.devices, device)
	return nil
}

// RemoveDevice removes a device by its ID. Returns an error if the device
// is not found.
func (p *Principal) RemoveDevice(deviceID string) error {
	for i, d := range p.devices {
		if d.ID() == deviceID {
			p.devices = append(p.devices[:i], p.devices[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("device with id %q not found", deviceID)
}

// AddCredential adds a credential to the principal.
func (p *Principal) AddCredential(cred vo.Credential) error {
	p.credentials = append(p.credentials, cred)
	return nil
}

// RecordAudit appends an audit entry to the principal's audit log.
func (p *Principal) RecordAudit(entry entity.AuditEntry) {
	p.auditLog = append(p.auditLog, entry)
}

// UpdateLastSeen updates the lastSeenAt timestamp to the current time.
func (p *Principal) UpdateLastSeen() {
	p.lastSeenAt = time.Now()
}

// @AI_GENERATED: end
