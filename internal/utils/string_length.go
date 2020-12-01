package utils

import (
	"fmt"
)

func StringLength(name, v string, min, max int) error {
	if len(v) < min || len(v) > max {
		return fmt.Errorf("%s need to be between %d and %d characters", name, min, max)
	}
	return nil
}
