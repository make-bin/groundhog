// @AI_GENERATED
package service

import (
	"github.com/make-bin/groundhog/pkg/domain/configuration/aggregate/config_root"
	"github.com/make-bin/groundhog/pkg/domain/configuration/vo"
)

// ConfigValidationService validates a ConfigRoot and returns any validation errors.
type ConfigValidationService interface {
	Validate(root *config_root.ConfigRoot) []vo.ValidationError
}

// @AI_GENERATED: end
