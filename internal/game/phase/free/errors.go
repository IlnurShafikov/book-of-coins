// Package free implements the free spins game phase logic.
// The free phase represents free spin rounds triggered during base or other free spins.
package free

import (
	"errors"
)

var ErrInvalidExpandSymbolID = errors.New("invalid expand symbol id")
