package errors

import (
	stderrors "errors"
)

// Wrap standard errors functions
var (
	New    = stderrors.New
	Is     = stderrors.Is
	As     = stderrors.As
	Unwrap = stderrors.Unwrap
)
