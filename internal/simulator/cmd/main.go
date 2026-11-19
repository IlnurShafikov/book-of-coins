// Command simulator runs Monte Carlo simulations of the game to verify RTP and volatility
// and writes the reports into the out directory.
package main

import (
	"log"

	"github.com/IlnurShafikov/book-of-coins/internal/simulator"
	"github.com/IlnurShafikov/book-of-coins/internal/simulator/runners"
)

// main parses the CLI flags and runs the simulation.
// A failed simulation terminates the program with a non-zero exit code.
func main() {
	cfg := simulator.ParseConfigs()

	simManager := runners.NewSimulationManager()
	if err := simManager.RunSimulation(cfg); err != nil {
		//nolint:forbidigo // it's ok for simulator
		log.Fatal(err)
	}
}
