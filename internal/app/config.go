package app

import (
	"fmt"
	"strings"
)

type Config struct {
	DatabasePath string
	Address      string
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.DatabasePath) == "" {
		return fmt.Errorf("database path is required")
	}
	if strings.TrimSpace(c.Address) == "" {
		return fmt.Errorf("address is required")
	}
	return nil
}
