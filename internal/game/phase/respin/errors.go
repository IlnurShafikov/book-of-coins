// Package free implements the free spins game phase logic.
// The free phase represents free spin rounds triggered during base or other free spins.
package respin

import (
	"errors"
)

// ErrInvalidIntervalCoinCount is returned when the number of coins matches no coin chance group.
var ErrInvalidIntervalCoinCount = errors.New("invalid interval of coin count")
