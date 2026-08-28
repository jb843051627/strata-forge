package service

import (
	"errors"
	"fmt"

	"github.com/jb843051627/strata-forge/internal/model"
)

func wrap(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", operation, err)
}

func isNotFound(err error) bool {
	return errors.Is(err, model.ErrNotFound)
}

func requirePositiveID(id int64, label string) error {
	if id <= 0 {
		return fmt.Errorf("%w: %s id must be positive", model.ErrInvalidInput, label)
	}
	return nil
}
