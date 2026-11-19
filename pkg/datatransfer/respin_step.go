package datatransfer

import (
	"fmt"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"

	"github.com/IlnurShafikov/book-of-coins/internal/entities/result"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/result/events"
	"github.com/IlnurShafikov/book-of-coins/internal/entities/rules"
	"github.com/IlnurShafikov/book-of-coins/pkg/datatransfer/dto/common"
	respinDto "github.com/IlnurShafikov/book-of-coins/pkg/datatransfer/dto/respin"
	"github.com/IlnurShafikov/book-of-coins/pkg/datatransfer/dto/step"
)

// TransformReSpinSteps transforms all respins into client-facing steps.
// Handles respins from either base game or free spins phase.
//
// Parameters:
//   - respinSteps: slice of respin results to transform
//   - allCliSteps: accumulated list of client steps
//   - parentStep: board and counter of the phase that triggered the respins
//   - baseParent: base game state to return to when max win is reached
//   - initRsCount: number of respins initially granted
//   - totalWin: cumulative win before respins
//   - betValue: the player's bet amount
//
// Returns:
//   - updated slice of client steps
//   - total win amount (including respins)
//   - error if transformation fails
func (t *Transfer) TransformReSpinSteps(
	respinSteps []*result.ReSpin, allCliSteps []*step.Step,
	parentStep, baseParent *common.ParentStep,
	initRsCount uint32, totalWin, betValue uint64,
) ([]*step.Step, uint64, error) {
	if initRsCount == 0 {
		return allCliSteps, totalWin, nil
	}

	counter := initCounter(initRsCount)

	for stepNum := range respinSteps {
		isLastStep := stepNum == len(respinSteps)-1

		cliStep, stepWin, err := t.processRespinStep(
			respinSteps[stepNum], isLastStep, &counter,
			parentStep, baseParent, betValue, totalWin)
		if err != nil {
			return nil, 0, err
		}

		totalWin += stepWin

		allCliSteps = append(allCliSteps, cliStep)
	}

	return allCliSteps, totalWin, nil
}

// transformParentStep transforms the parent base phase board into a ParentStep DTO.
// This is used to send the parent board state to the client when respins complete.
//
// Parameters:
//   - parentBoardResult: the parent phase board result to transform
//   - parentCounter: counter of the parent phase (nil if the parent is the base game)
//   - betValue: the player's bet amount for symbol property scaling
//
// Returns:
//   - pointer to ParentStep DTO containing parent base phase data
//   - error if symbol transformation fails
func transformParentStep(
	parentBoardResult *result.BoardResult,
	parentCounter *common.Counter,
	betValue uint64,
) (*common.ParentStep, error) {
	parentSymbols, err := TransformSymbols(parentBoardResult, betValue)
	if err != nil {
		return nil, fmt.Errorf("transform parent symbols: %w", err)
	}

	parentStep := &common.ParentStep{
		Payload: common.ParentPayload{
			Symbols: parentSymbols,
			Counter: parentCounter,
		},
	}

	return parentStep, nil
}

// processRespinStep transforms a single respin into a client step.
// Handles symbol transformation, event processing, counter updates, and parent phase context.
//
// Parameters:
//   - rsRes: the respin result to transform
//   - isLastStep: whether this is the final respin in the sequence
//   - counter: the respins counter to update
//   - parentStep: symbol board and counter from parent phase
//   - baseParent: base game state to return to when max win is reached
//   - betValue: the player's bet amount
//   - totalWin: cumulative win before this respin
//
// Returns:
//   - pointer to protobuf Step message
//   - win amount for this respin
//   - error if transformation fails
func (t *Transfer) processRespinStep(
	rsRes *result.ReSpin, isLastStep bool, counter *common.Counter,
	parentStep, baseParent *common.ParentStep, betValue, totalWin uint64,
) (*step.Step, uint64, error) {
	stepWin := rsRes.StepWin() * betValue

	symbols, err := TransformSymbols(rsRes.BoardResult(), betValue)
	if err != nil {
		return nil, 0, fmt.Errorf("transform symbols: %w", err)
	}

	eventList, err := ExtractAndTransformEvents(rsRes, betValue)
	if err != nil {
		return nil, 0, fmt.Errorf("transform events: %w", err)
	}

	updateRsCounter(counter, rsRes, isLastStep)

	payload := buildRespinPayload(
		*counter, eventList, rsRes, symbols, parentStep, baseParent, isLastStep)

	currStep := buildRespinStep(payload, totalWin+stepWin, stepWin, rsRes)

	data, err := jsoniter.Marshal(currStep)
	if err != nil {
		return nil, 0, fmt.Errorf("encode free step: %w", err)
	}

	id, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, 0, fmt.Errorf("generate uuid: %w", err)
	}

	cliStep := makeRsClientStep(id, data, stepWin, totalWin+stepWin, rsRes)

	finalizeRSCounter(counter)

	return cliStep, stepWin, nil
}

// updateRsCounter updates the respins counter based on the step result.
// Handles retriggers, max win and last step scenarios.
//
// Parameters:
//   - counter: the counter to update
//   - rsRes: the respin result
func updateRsCounter(counter *common.Counter, rsRes *result.ReSpin, isLastStep bool) {
	if rsRes.ReTriggers() > 0 {
		counter.Added = rsRes.ReTriggers()
		counter.Left += rsRes.ReTriggers()
	}

	if rsRes.MaxWinReached() || isLastStep {
		counter.Left = 0
		counter.Added = 0
	}
}

// finalizeRSCounter updates the counter after processing a respin.
// Increments played count, decrements remaining, and resets added.
//
// Parameters:
//   - counter: the counter to finalize
func finalizeRSCounter(counter *common.Counter) {
	counter.Played++
	counter.Left--
	counter.Added = 0
}

// buildRespinPayload constructs the payload for a respin step DTO.
// Handles the complex logic of determining what information to include based on
// whether this is the last respin and whether we're returning to free spins.
//
// Parameters:
//   - counter: current state of the respins counter
//   - eventList: game events that occurred
//   - rsRes: the respin result
//   - symbols: current symbol board
//   - parentStep: board and counter of the phase that triggered the respins
//   - baseParent: base game state to return to when max win is reached
//   - isLastStep: whether this is the final respin in the sequence
//
// Returns:
//   - pointer to respin Payload DTO
func buildRespinPayload(
	counter common.Counter,
	eventList []events.Event,
	rsRes *result.ReSpin,
	symbols []*common.Symbol,
	parentStep, baseParent *common.ParentStep,
	isLastStep bool,
) *respinDto.Payload {
	var (
		curParentStep *common.ParentStep
		stepsLeftInFS = 0
	)

	if parentStep != nil && parentStep.Payload.Counter != nil {
		stepsLeftInFS = int(parentStep.Payload.Counter.Left)
	}

	// Show parent and after board on the last respin
	if isLastStep {
		curParentStep = parentStep
	}

	// If max win was reached, always return to the base game board.
	if rsRes.MaxWinReached() || (stepsLeftInFS == 0 && isLastStep) {
		curParentStep = baseParent
	}

	return &respinDto.Payload{
		Symbols:    symbols,
		ParentStep: curParentStep,
		Events:     eventList,
		Counter:    counter,
		NextStepType: getNextStepTypeFromRs(
			stepsLeftInFS, isLastStep, rsRes.MaxWinReached()),
	}
}

// getNextStepTypeFromRs determines what type of step comes after this respin.
// Logic accounts for whether respins are complete, max win reached, and parent phase.
//
// Parameters:
//   - stepsLeftInFs: number of free spins remaining (0 if from base game)
//   - maxWinReached: whether max win was reached
//
// Returns:
//   - TYPE_BASE_SPIN if max win reached or returning to base
//   - TYPE_FREE_SPIN if returning to free spins
//   - TYPE_RE_SPIN if more respins remain
func getNextStepTypeFromRs(
	stepsLeftInFs int, isLastStep, maxWinReached bool,
) int {
	if maxWinReached {
		return rules.Step_TYPE_BASE_SPIN
	}

	if isLastStep {
		if stepsLeftInFs > 0 {
			return rules.Step_TYPE_FREE_SPIN
		}

		return rules.Step_TYPE_BASE_SPIN
	}

	return rules.Step_TYPE_RE_SPIN
}

// buildRespinStep constructs a complete respin Step DTO.
//
// Parameters:
//   - payload: the step payload with game data
//   - totalWin: cumulative win amount
//   - stepWin: win amount for this respin
//   - rsRes: the respin result
//
// Returns:
//   - pointer to respin Step DTO
func buildRespinStep(
	payload *respinDto.Payload, totalWin, stepWin uint64, rsRes *result.ReSpin,
) *respinDto.Step {
	return &respinDto.Step{
		Payload:       *payload,
		TotalWin:      totalWin,
		StepWin:       stepWin,
		Type:          rules.Step_TYPE_RE_SPIN,
		MaxWinReached: rsRes.MaxWinReached(),
	}
}

// makeRsClientStep constructs the protobuf Step message for a respin.
//
// Parameters:
//   - id: unique identifier for the step
//   - data: JSON-serialized step data
//   - stepWin: win amount for this respin
//   - totalWin: cumulative win amount
//   - rsRes: the respin result
//
// Returns:
//   - pointer to protobuf Step message
func makeRsClientStep(
	id, data []byte, stepWin, totalWin uint64, stepRs *result.ReSpin,
) *step.Step {
	return &step.Step{
		Id:             id,
		ClientResponse: data,
		Cost:           0,
		StepWin:        stepWin,
		TotalWin:       totalWin,
		Type:           rules.Step_TYPE_RE_SPIN,
		StepsAdded:     stepRs.ReTriggers(),
		MaxWinReached:  stepRs.MaxWinReached(),
	}
}
