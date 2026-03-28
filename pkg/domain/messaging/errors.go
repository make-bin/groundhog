package messaging

import "errors"

var (
	ErrMessageNotFound     = errors.New("message not found")
	ErrChannelNotFound     = errors.New("channel not found")
	ErrInvalidMessageState = errors.New("invalid message state transition")
	ErrRoutingFailed       = errors.New("routing failed")
	ErrChannelInactive     = errors.New("channel is inactive")
)
