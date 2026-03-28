package cron

import "errors"

var (
	ErrCronJobNotFound        = errors.New("cron job not found")
	ErrCronJobAlreadyExists   = errors.New("cron job already exists")
	ErrCronJobAlreadyRunning  = errors.New("cron job is already running")
	ErrInvalidSchedule        = errors.New("invalid cron schedule")
	ErrInvalidSessionTarget   = errors.New("invalid session target")
	ErrPayloadSessionMismatch = errors.New("payload type does not match session target")
)
