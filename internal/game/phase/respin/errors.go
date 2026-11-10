// Package free implements the free spins game phase logic.
// The free phase represents free spin rounds triggered during base or other free spins.
package respin

import (
	"errors"
)

var ErrInvalidIntervalCoinCount = errors.New("invalid interval of coin count")
