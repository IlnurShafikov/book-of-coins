package events

import (
	"cmp"
	"slices"

	"github.com/IlnurShafikov/slot-engine/pkg/calculators"
	"github.com/IlnurShafikov/slot-engine/pkg/model"
)

// WinsEvent represents an event containing all winning lines and scatter wins from a spin.
// It includes all winning positions, line details with payouts, and scatter wins.
type WinsEvent struct {
	// Event is the event type identifier
	Event Type `json:"type"`

	// Positions contains all unique positions involved in wins (both line and scatter)
	Positions model.Positions `json:"positions"`
	// AnyWhereWins contains all scatter (pay anywhere) wins
	AnyWhereWins []AnywhereWin `json:"anyWhere,omitempty"`
	// LineWins contains all line-based wins
	LineWins []LineWin `json:"lines,omitempty"`
}

// NewWinsEvent creates a new wins event from calculated line wins and scatter wins.
// It extracts unique positions from both win types and constructs details with position indices.
//
// Parameters:
//   - lineWins: calculated line wins from the game engine
//   - scatterWins: calculated scatter (pay anywhere) wins from the game engine
//
// Returns:
//   - pointer to a new WinEvent containing both line and scatter wins
func NewWinsEvent(
	lineWins calculators.LineWins,
	scatterWins calculators.PayAnyWhereWins,
) *WinsEvent {
	scatterUniquePos := collectUniqueAnywherePositions(scatterWins)
	lineUniquePos := collectUniqueLinePositions(lineWins)
	allPositions := mergeUniquePositions(scatterUniquePos, lineUniquePos)

	winAnywhere := transferAnyWhere(scatterWins, allPositions)
	lines := transferLines(lineWins, allPositions)

	return &WinsEvent{
		Event:        WinsType,
		Positions:    allPositions,
		AnyWhereWins: winAnywhere,
		LineWins:     lines,
	}
}

// EventType returns the type identifier for this event.
func (c *WinsEvent) EventType() Type {
	return c.Event
}

// LineWin represents a single winning line with its positions and payout.
// Position indices refer to entries in the LineWinsEvent.Positions array.
type LineWin struct {
	// PosIndices contains indices into the Positions array indicating winning positions
	PosIndices []int `json:"posIndices"`

	// Line is the line number (pay line index)
	Line int `json:"line"`

	// Symbol is the symbol ID that formed the winning combination
	Symbol int `json:"symbol"`

	// Payout is the win amount for this line (before bet multiplication)
	Payout uint64 `json:"payout"`
}

// AnywhereWin represents a single scatter (pay anywhere) win with its positions and payout.
// Position indices refer to entries in the WinEvent.Positions array.
type AnywhereWin struct {
	// PosIndices contains indices into the Positions array indicating winning positions
	PosIndices []int `json:"posIndices"`

	// Payout is the win amount for this scatter win
	Payout uint64 `json:"payout"`

	// Symbol is the symbol ID that formed the scatter win
	Symbol model.Symbol `json:"symbol"`
}

// transferAnyWhere converts scatter wins into the event format with position indices.
// It collects unique positions and maps them to indices for compact representation.
//
// Parameters:
//   - wins: scatter wins from the game engine
//
// Returns:
//   - slice of AnywhereWin with mapped position indices
func transferAnyWhere(
	wins calculators.PayAnyWhereWins,
	uniquePositions model.Positions,
) []AnywhereWin {
	if len(wins) == 0 {
		return nil
	}

	positionIndices := buildPositionIndices(uniquePositions)

	result := make([]AnywhereWin, len(wins))
	for i, win := range wins {
		result[i] = AnywhereWin{
			Payout:     win.Payout,
			Symbol:     win.Symbol,
			PosIndices: mapToIndices(win.Positions, positionIndices),
		}
	}

	return result
}

// transferLines converts line wins into the event format with position indices.
// It collects unique positions and maps them to indices for compact representation.
//
// Parameters:
//   - wins: line wins from the game engine
//
// Returns:
//   - slice of LineWin with mapped position indices
func transferLines(
	wins calculators.LineWins,
	uniquePositions model.Positions,
) []LineWin {
	if len(wins) == 0 {
		return nil
	}

	positionIndices := buildPositionIndices(uniquePositions)

	result := make([]LineWin, len(wins))
	for i, win := range wins {
		result[i] = LineWin{
			Payout:     win.Payout,
			Symbol:     win.Symbol,
			PosIndices: mapToIndices(win.Positions, positionIndices),
			Line:       win.Number,
		}
	}

	return result
}

// collectUniqueAnywherePositions extracts all unique positions from scatter wins.
// Positions are deduplicated and sorted by reel and row indices.
//
// Parameters:
//   - wins: scatter wins to extract positions from
//
// Returns:
//   - sorted slice of unique positions
func collectUniqueAnywherePositions(wins calculators.PayAnyWhereWins) model.Positions {
	return collectPositions(func(yield func(model.Position)) {
		for _, win := range wins {
			for _, pos := range win.Positions {
				yield(pos)
			}
		}
	})
}

// collectUniqueLinePositions extracts all unique positions from line wins.
// Positions are deduplicated and sorted by reel and row indices.
//
// Parameters:
//   - wins: line wins to extract positions from
//
// Returns:
//   - sorted slice of unique positions
func collectUniqueLinePositions(wins calculators.LineWins) model.Positions {
	return collectPositions(func(yield func(model.Position)) {
		for _, win := range wins {
			for _, pos := range win.Positions {
				yield(pos)
			}
		}
	})
}

// collectPositions is a helper function that collects unique positions using an iterator.
// It deduplicates positions and returns them sorted by reel and row indices.
//
// Parameters:
//   - iter: iterator function that yields positions
//
// Returns:
//   - sorted slice of unique positions
func collectPositions(iter func(yield func(model.Position))) model.Positions {
	positionMap := make(map[model.Position]struct{})

	iter(func(pos model.Position) {
		positionMap[pos] = struct{}{}
	})

	positions := make(model.Positions, 0, len(positionMap))
	for pos := range positionMap {
		positions = append(positions, pos)
	}

	slices.SortFunc(positions, comparePositions)

	return positions
}

// buildPositionIndices creates a mapping from positions to their array indices.
// This is used to convert positions to compact index references.
//
// Parameters:
//   - positions: slice of positions to index
//
// Returns:
//   - map from position to its index in the slice
func buildPositionIndices(positions model.Positions) map[model.Position]int {
	indices := make(map[model.Position]int, len(positions))
	for i, pos := range positions {
		indices[pos] = i
	}

	return indices
}

// mapToIndices converts a slice of positions to their corresponding indices.
// Uses the position-to-index mapping to create a compact array of indices.
//
// Parameters:
//   - positions: positions to convert to indices
//   - positionIndices: mapping from position to index
//
// Returns:
//   - slice of indices corresponding to the input positions
func mapToIndices(positions model.Positions, positionIndices map[model.Position]int) []int {
	indices := make([]int, 0, len(positions))

	for _, pos := range positions {
		if index, ok := positionIndices[pos]; ok {
			indices = append(indices, index)
		}
	}

	return indices
}

// comparePositions is a comparison function for sorting positions.
// Sorts first by reel index, then by row index.
//
// Parameters:
//   - src: first position to compare
//   - target: second position to compare
//
// Returns:
//   - negative if src < target, 0 if equal, positive if src > target
func comparePositions(src, target model.Position) int {
	if src.ReelIndex != target.ReelIndex {
		return cmp.Compare(src.ReelIndex, target.ReelIndex)
	}

	return cmp.Compare(src.RowIndex, target.RowIndex)
}

// mergeUniquePositions combines positions from two slices, removing duplicates.
// The result is sorted by reel and row indices.
//
// Parameters:
//   - posSlice1: first slice of positions
//   - posSlice2: second slice of positions
//
// Returns:
//   - sorted slice containing unique positions from both inputs
func mergeUniquePositions(posSlice1, posSlice2 model.Positions) model.Positions {
	positionMap := make(map[model.Position]struct{}, len(posSlice1)+len(posSlice2))

	for _, pos := range posSlice1 {
		positionMap[pos] = struct{}{}
	}

	for _, pos := range posSlice2 {
		positionMap[pos] = struct{}{}
	}

	result := make(model.Positions, 0, len(positionMap))
	for pos := range positionMap {
		result = append(result, pos)
	}

	slices.SortFunc(result, comparePositions)

	return result
}
