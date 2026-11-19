package common

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/getter"
	"github.com/IlnurShafikov/book-of-coins/internal/configs"
	"github.com/IlnurShafikov/book-of-coins/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

// ReadValueGetter reads the coin values of a phase and builds the coin value getter.
//
// Parameters:
//   - csvReader: CSV file adapter for reading configuration files
//   - phasePath: path builder pointing to the phase directory
//   - general: general configuration providing the jackpot types
//
// Returns:
//   - pointer to CoinValue
//   - error if the coin values cannot be read
func ReadValueGetter(
	csvReader *reader.Adapter,
	phasePath *tools.PathBuilder,
	general *configs.GeneralConfig,
) (*getter.CoinValue, error) {
	const fileName = "coin_values"

	coinValues, err := ReadWeight(csvReader, phasePath.File(fileName))
	if err != nil {
		return nil, fmt.Errorf("read coin values: %w", err)
	}

	valueGetter := getter.NewCoinValue(
		coinValues,
		general.JpTypes)

	return valueGetter, nil
}
