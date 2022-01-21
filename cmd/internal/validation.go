package internal

import "fmt"

func ValidatePixelValue(v int) error {
	if v < 0 || v > 65535 {
		return fmt.Errorf("value must be in range [0, 65535]")
	}
	return nil
}

func ValidatePositivePixelValue(v int) error {
	if v < 1 || v > 65535 {
		return fmt.Errorf("value must be in range [1, 65535]")
	}
	return nil
}
