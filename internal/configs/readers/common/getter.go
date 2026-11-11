package common

import (
	"fmt"

	"github.com/IlnurShafikov/book-of-coins/internal/components/getter"
	"github.com/IlnurShafikov/book-of-coins/internal/configs"
	"github.com/IlnurShafikov/book-of-coins/internal/tools"
	"github.com/IlnurShafikov/slot-engine/pkg/tools/reader"
)

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
