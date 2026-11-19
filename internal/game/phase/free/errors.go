// Package free implements the free spins game phase logic.
// The free phase represents free spin rounds triggered during base or other free spins.
package free

import (
	"errors"
)

// ErrInvalidExpandSymbolID is returned when the special symbol has no minimum expand count.
var ErrInvalidExpandSymbolID = errors.New("invalid expand symbol id")
