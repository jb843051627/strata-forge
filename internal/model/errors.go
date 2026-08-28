package model

import "errors"

var (
	ErrNotFound          = errors.New("strata forge: not found")
	ErrInvalidInput      = errors.New("strata forge: invalid input")
	ErrInvalidTransition = errors.New("strata forge: invalid state transition")
	ErrConflict          = errors.New("strata forge: state conflict")
	ErrCancelled         = errors.New("strata forge: operation cancelled")
	ErrQualityHold       = errors.New("strata forge: quality hold")
	ErrAlreadyArchived   = errors.New("strata forge: already archived")
)
