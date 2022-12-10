package error_play

import (
	"errors"
	"fmt"
)

var (
	errorInternal = errors.New("file_not_found")
)

func genError(level int) error {
	level1Err := fmt.Errorf("read_saved_file error: %w", errorInternal)
	if level == 1 {
		return level1Err
	}
	if level == 2 {
		return fmt.Errorf("auto_load error: %w", level1Err)
	}

	return errorInternal
}

func Play_Error_Wrap_UnWrap() {
	err := genError(1)
	if errors.Is(err, errorInternal) {
		fmt.Printf("is errorInternal: Yes - full detail: %v\n", err)
	}
	fmt.Printf("unwrapped error: %v\n", errors.Unwrap(err))

	fmt.Printf("---\n")

	err = genError(2)
	if errors.Is(err, errorInternal) {
		fmt.Printf("is errorInternal:  Yes - full detail: %v\n", err)
	}
	fmt.Printf("error:\n  %v\n", err)
	fmt.Printf("unwrapped error:\n  %v\n", errors.Unwrap(err))
	fmt.Printf("unwrapped unwrapped error:\n  %v\n", errors.Unwrap(errors.Unwrap(err)))
}
