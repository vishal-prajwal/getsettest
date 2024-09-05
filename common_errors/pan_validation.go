package common_errors

import "fmt"

var (
	ErrorInvalidPanNumber = fmt.Errorf("Invalid Pan Number")
	ErrNotIndividualPan   = fmt.Errorf("Not a individual pan number")
)
