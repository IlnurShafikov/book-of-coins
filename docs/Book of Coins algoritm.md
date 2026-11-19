# Game algorithm: Book 5x3, 10 lines, expanding special symbol, Coin + Hold & Win (working title: Book of Coins v2.5)

Symbol IDs start at 1. 1 High1, 2 High2, 3 High3, 4 High4 (high-pay, 1 is the most valuable), 5 Low1, 6 Low2, 7 Low3, 8 Low4, 9 Low5 (low-pay), 10 Wild/Scatter (one symbol that is both Wild and Scatter), 11 Coin. ID 0 = EMPTY is used only for an empty Hold & Win cell and never appears on strips. `-1` is not a symbol: it only pads shorter strip columns in `group_N.csv` to equal length.

A **round** is one paid base spin plus every feature it triggers. All payouts of a round are summed and capped together (step 1.9).

## Configuration files

All configs are CSV files in `internal/configs/csv/`, embedded into the binary.

- `general/` is shared by all RTP versions: `paytable.csv`, `scatter_paytable.csv`, `win_lines.csv`, `symbols_on_reel.csv`, `simple_cost_multiplier.csv`, `max_win_limit.csv`, `fs_init_counter.csv`, `expand_min_count.csv`, `rs_coin_intervals.csv`, `jp_types.csv`, `grand_win.csv`, `grand_jp_type.csv`, `coin_count_to_grand.csv`, `volatility.csv`.
- `<rtp>/exp_rtp.csv` holds the expected RTP of the version (`<rtp>` is 90, 92, 94 or 96).
- `<rtp>/simple/base/` holds `distributions.csv`, `expand_symbol.csv`, `coin_values.csv` and `reels/group_N.csv` of the base game.
- `<rtp>/simple/free/` holds `distributions.csv`, `coin_values.csv` and `reels/group_N.csv` of the free spins.
- `<rtp>/simple/respin/` holds `coin_intervals_group.csv`, `coin_chance_group.csv` and `coin_values.csv` of Hold & Win (the respin phase).

**Units.** 1 unit = the bet of one line. A spin costs 10 units (1 total bet = 10 units). Every payout and value in the CSV files (`paytable.csv`, `scatter_paytable.csv`, all `coin_values.csv`, `jp_types.csv`, `grand_win.csv`) is stored in units. Only `max_win_limit.csv` is in total bets. The simulator and the client DTOs divide or multiply by the cost multiplier, the engine itself works in units.

1. **Spin.**

1.1. **Bet.** One spin costs `simple_cost_multiplier.csv` = 10 units (10 lines x 1 unit; 1 unit = line bet).

1.2. **Reel set selection.** One draw per spin, no draw per reel. All 5 reels use the strips of the chosen `group_N.csv` (column = reel 1..5, row = strip position).

- Base game: group by weights from `base/distributions.csv` (group 1: 82, group 2: 18 in all versions).
  - base group_1 has Wild/Scatter on all reels and no Coins, so Hold & Win cannot trigger on it.
  - base group_2 has Coins on all reels and Wild/Scatter only on reels 1-2, so at most 2 Wild/Scatters can land and free spins cannot trigger on it.
  - Therefore free spins and Hold & Win never trigger on the same base spin (in the code free spins have priority anyway).
- Free spins: the group is **not** drawn. The free spin uses `free/reels/group_N.csv` where N is the ID of the session's special symbol: symbol 1-9 -> group 1-9, Coin 11 -> group 11. `free/distributions.csv` (weight 1 for each of the groups 1-9 and 11) only lists the groups whose strips must be loaded. Groups 1-9 have no Coins; group 11 has Coins.
- Hold & Win uses no strips (step 1.7).

  1.3. **Reel stops.** For each reel independently, the stop is chosen uniformly among the non-`-1` positions of that reel's column. Strips are circular. The window is `symbols_on_reel.csv` = 3,3,3,3,3 consecutive symbols at positions stop, stop+1, stop+2 (mod strip length); the stop position is row 0 (top).
  Strip guarantees (by construction, verified):

- Any two Wild/Scatters on one strip are at least 3 positions apart, so a reel window holds at most one Wild/Scatter.
- In free groups 1-9, any two copies of the group's special symbol (group N belongs to symbol N) are at least 3 positions apart; in free group 11, any two Coins are at least 3 positions apart. In free spins every reel window therefore holds at most one special symbol. The code relies on this: it counts the special symbols on the board and takes the reel of each of them as an expanded reel.
- In base group_2, Coins stand in stacks of at most 2. No other symbol has identical neighbours (circularly).

  1.4. **Board modifications: Coin values.** No symbol on the board is changed. Every Coin that lands receives one value at once, drawn independently by weights from the `coin_values.csv` (`value in units,weight`) of the current phase:

- base game Coins: `base/coin_values.csv`;
- free spin Coins (group 11): `free/coin_values.csv` (it has no jackpot values, so an expansion never copies a jackpot);
- new Coins in Hold & Win respins: `respin/coin_values.csv`.
  Cash values are 10, 20, 30, 40, 50, 80, 100, 150 units (1, 2, 3, 4, 5, 8, 10, 15 total bets). A value that equals a value of `jp_types.csv` (250 / 500 / 2500 units = 25 / 50 / 250 total bets) is a jackpot Coin with the jackpot type 1 MINI, 2 MINOR, 3 MAJOR; the type is sent to the client. GRAND (type 4, `jp_types.csv` = 50000 units = 5000 total bets) is never a Coin value (step 1.7).

  1.5. **Win evaluation.**

  1.5.1. **Lines.** For each of the 10 lines in `win_lines.csv` (row index per reel, 0 = top), read reels 1 to 5 and compute two candidate wins:

- Wild/Scatter combination: k = number of consecutive Wild/Scatters (ID 10) from reel 1. Win = `paytable[10][k]`, which equals `paytable[1][k]` (2 -> 10, 3 -> 100, 4 -> 1000, 5 -> 5000 units); 0 for k = 0 or 1.
- Symbol combination: if k = 5 there is none. Otherwise X = the symbol on reel k+1 of the line (the first symbol that is not Wild/Scatter). If X = Coin (11), it pays 0. Otherwise count = consecutive positions from reel 1 whose symbol is X or Wild/Scatter (stops at the first other symbol, including a Coin); win = `paytable[X][count]` (High1-High4 pay from 2, Low1-Low5 from 3).
- The line pays the larger of the two. One payout per line; all lines pay independently. Wild/Scatter substitutes for 1-9 only, never for the Coin: in the code the Coin is the symbol that the line calculator skips.
  Examples: WS, WS, Low2, Low2, Low3 -> max(10, `paytable[6][4]` = 30) = 30. WS, WS, Coin -> 10. WS x5 -> 5000.

  1.5.2. **Scatter.** n = Wild/Scatters anywhere on the board (0-5, one per reel). If n >= 3, pay `scatter_paytable.csv[n]` units (3 -> 20, 4 -> 200, 5 -> 2000, i.e. 2 / 20 / 200 total bets), in addition to line wins, in the base game and in free spins.

  1.6. **Summation.** Spin win = sum of the 10 line wins + scatter win. There are no win multipliers.

  1.7. **Hold & Win.** In the code this feature is the respin phase.

- **Trigger in the base game:** `rs_coin_intervals.csv` maps the number of Coins on the board to the number of respins: 0-5 Coins -> 0, 6-15 Coins -> 4. So 6 or more Coins trigger it (possible only in base group_2). It is played after the spin win of step 1.6 is paid. Start board = the landed Coins with their values from step 1.4.
- **Trigger in free spins:** only through the Coin expansion of step 1.8. Start board = the expanded Coins (the same table `rs_coin_intervals.csv` is used).
- All other cells are EMPTY. Held Coins never change value or position.
- Respin counter = 4 at the start (`rules.InitRespinsCount` in the code; the count in `rs_coin_intervals.csv` is the number of respins awarded by the trigger).
- Each respin:
  1. m = number of Coins on the board at the start of this respin, recounted on every respin. `coin_intervals_group.csv` maps m to a coin chance group: 6-8 -> group 1, 9-10 -> group 2, 11-12 -> group 3, 13-15 -> group 4.
  2. `coin_chance_group.csv` has the columns `group, empty weight, coin weight` (group 1: 945/55, group 2: 960/40, group 3: 975/25, group 4: 994/6). Every EMPTY cell independently becomes a Coin with the coin weight or stays EMPTY with the empty weight. All cells of one respin use the same group; Coins landing in this respin change m only from the next respin.
  3. Each new Coin gets a value from `respin/coin_values.csv` (step 1.4) and is held.
  4. If at least one new Coin landed the counter is reset to 4, otherwise it decreases by 1.
- End: the counter is 0 or all 15 cells hold Coins.
- Payout: the sum of all Coin values is paid **once, on the last respin**. A full board (`coin_count_to_grand.csv` = 15 Coins) adds GRAND = `grand_win.csv` = 50000 units (5000 total bets, jackpot type `grand_jp_type.csv` = 4) on top of the Coin sum. Both go through the win cap of step 1.9. The client receives a collect-coin event and, for a full board, a grand-win event.
- No line or scatter wins during Hold & Win. Hold & Win does not use free spins.

  1.8. **Free spins.**

- **Trigger:** `fs_init_counter.csv` maps the number of Wild/Scatters to the number of free spins: 0-2 -> 0, 3, 4 and 5 -> 10 spins. It is evaluated on a base spin (possible only in base group_1). The trigger spin win is paid first.
- **Special symbol:** chosen once per session by weights from `base/expand_symbol.csv` (version-dependent, see below; v96: 1:3, 2:5, 3:7, 4:9, 5:11, 6:12, 7:14, 8:16, 9:18, Coin 11:5). It never changes, also after a retrigger.
- **Each free spin, in order:**
  1. Reel set for the special symbol (step 1.2), stops (step 1.3), Coin values from `free/coin_values.csv` (step 1.4).
  2. Line and scatter wins (steps 1.5-1.6).
  3. Expansion:
     - Regular special symbol S (1-9): n = number of S on the board (one per reel, see the strip guarantees). If n >= `expand_min_count.csv[S]` (2 for 1-4, 3 for 5-9), those reels are covered by S and the expansion pays `paytable[S][n]` units on each of the 10 lines, i.e. 10 x `paytable[S][n]` units (every line crosses all n reels, adjacent or not). It is paid in addition to sub-step 2 and does not carry over to the next spin. **A regular special symbol never starts Hold & Win.**
     - Coin special symbol (11): n = number of Coins on the board (at most one per reel). If n >= `expand_min_count.csv[11]` = 2, each such reel becomes 3 Coins with the value of the Coin that landed on that reel, and Hold & Win (step 1.7) starts from these 3n Coins. After Hold & Win the session continues; Hold & Win can start again on any later free spin.
  4. Retrigger: the Wild/Scatters of the spin add spins by the same table `fs_init_counter.csv` (3+ Wild/Scatters add 10 spins). No limit on retriggers or on the total number of free spins.
- **End:** no spins remain, or the cap of step 1.9 is reached (the cap is checked before every free spin).

  1.9. **Win cap.** `max_win_limit.csv` = 5000 total bets (50000 units) per round. Every award (spin win, expansion win, Hold & Win payout) goes through one win limiter that is shared by the whole round: the part of the award that would go over the cap is cut, the cumulative round win is set to exactly the cap and the round ends. Free spins that are not played yet are forfeited; a Hold & Win session that is already running finishes without paying more, and the client sees the round end at the step that reached the cap.

  1.10. **Spin result.** Round payout in total bets = round win in units (after the cap) / 10.

## RTP versions

Versions 96, 94, 92 and 90 use the same algorithm, strips, paytable and all other files. They differ only in `base/expand_symbol.csv`, the weights of `base/coin_values.csv` and `exp_rtp.csv` (expected RTP 96.6%, 94.4%, 92.4%, 90.14%).

| Special symbol | 1   | 2   | 3   | 4   | 5   | 6   | 7   | 8   | 9   | 11  |
| -------------- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| v96            | 3   | 5   | 7   | 9   | 11  | 12  | 14  | 16  | 18  | 5   |
| v94            | 2   | 4   | 6   | 8   | 11  | 13  | 16  | 19  | 22  | 5   |
| v92            | 2   | 3   | 5   | 7   | 10  | 13  | 17  | 22  | 27  | 5   |
| v90            | 1   | 2   | 3   | 6   | 9   | 13  | 19  | 24  | 31  | 5   |

| base coin weight, value (total bets) | 1   | 2   | 3   | 4   | 5   | 8   | 10  | 15  | MINI | MINOR | MAJOR |
| ------------------------------------ | --- | --- | --- | --- | --- | --- | --- | --- | ---- | ----- | ----- |
| v96                                  | 580 | 490 | 300 | 200 | 165 | 95  | 62  | 38  | 12   | 4     | 1     |
| v94                                  | 585 | 495 | 300 | 200 | 163 | 92  | 58  | 35  | 12   | 4     | 1     |
| v92                                  | 665 | 525 | 300 | 185 | 148 | 78  | 43  | 24  | 11   | 3     | 1     |
| v90                                  | 700 | 550 | 300 | 175 | 135 | 65  | 38  | 20  | 10   | 3     | 1     |

## Notes for the implementer

- Hold & Win runs from two places (base trigger and Coin expansion in free spins); both use the same respin routine, and the coin chance group is chosen again on every respin.
- `expand_min_count.csv` is read by symbol ID (the Coin row is `11,2`), not by row position.
- In free spins only the Coin special symbol starts Hold & Win.
- The respin counter start (4) is a constant in the code (`rules.InitRespinsCount`), the number of respins awarded by the trigger comes from `rs_coin_intervals.csv`; keep the two values equal.
- The first prototype run had exactly these three errors; see the simulator results doc.
