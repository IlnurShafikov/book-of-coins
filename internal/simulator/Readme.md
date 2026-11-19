# Math simulator

A Monte Carlo simulator that plays the game through the same `game.Mod` as the production code and
collects the statistics of the math model: RTP split by feature, standard deviation, volatility index,
hit rate, max win and win-range histograms.

## Modes

### 1. Simple mode: statistical analysis

Verifies RTP, volatility and the other metrics of an RTP version.

```bash
go run ./internal/simulator/cmd -spins 10000000 -rtp 96 -workers 4
```

The report is written to `out/v96/cost10/simple/seed<seed>/<date>.txt`. Rounds are split between the
workers by spin index, so a run is reproducible for the same `-seed` and `-workers`.

### 2. Interactive mode: step-by-step debugging

Plays one spin per Enter press and prints the result of every phase.

```bash
go run ./internal/simulator/cmd -play -spins 10
```

### 3. Player mode: balance simulation

Plays players with the given start balances (in bets) until they run out of money or the spins end
and reports the median number of spins they survive.

```bash
go run ./internal/simulator/cmd -playerMod -balance "10,20,30,40" -playersCount 100
```

The report is written to `out/player_simulation/v96/cost10/simple/players100/seed<seed>/<date>.txt`.

## Parameters

| Parameter | Default | Description |
|---|---|---|
| `-spins` | 100 000 000 | Number of rounds |
| `-rtp` | 96 | RTP version (90, 92, 94, 96) |
| `-workers` | 1 | Parallel workers; `1` runs sequentially |
| `-seed` | time based | Main seed for reproducible results |
| `-play` | false | Interactive mode |
| `-playerMod` | false | Player balance simulation |
| `-balance` | empty | Comma-separated start balances in bets |
| `-playersCount` | 100 | Players per balance group |
| `-bet` | 1 | Bet value |
| `-costMultiplier` | 10 | Cost multiplier of the simple mode, scales the win range buckets |
| `-mod` / `-feature` | `simple` / empty | Game mode and optional feature |
| `-rng` | `mt` | RNG backend |
| `-description` | empty | Free-form label stored in the report |

## Examples

```bash
# check the RTP of the version 96
go run ./internal/simulator/cmd -spins 10000000 -rtp 96 -workers 4

# a reproducible run
go run ./internal/simulator/cmd -spins 1000000 -seed 42

# session length with different start balances
go run ./internal/simulator/cmd -playerMod -balance "50,100,200" -playersCount 500
```
