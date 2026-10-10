# Engine

The **Engine** component is the central authority managing game state, world navigation, turn-based combat, quest validation, inventory, and player grouping for *The Answer Protocol* (TAP).

---

## Concurrency & Event Loop Model

The Engine uses a single-threaded broadcaster event-loop model to eliminate race conditions and avoid global mutex locking across shared game entities:

- **Broadcaster Loop (`broadcaster`)**: Processes incoming requests sequentially from the `ServerInput` channel of `pr.Exchanger`.
- **Session Registry (`sessions` & `players`)**: Maps TCP connection UUIDs to player identities and active `Player` structures.
- **Panic Recovery**: Encapsulates the main loop in a defer/recover mechanism to log unexpected panics without shutting down the server process.

---

## Core Systems

### 1. World & Navigation
- Loads world definitions from `src/world.json` into room, item, NPC, and quest lookup tables (`Map`).
- Validates directional movements (`MOVE north`, `MOVE east`, etc.) and enforces room availability.

### 2. Turn-Based Combat Engine (`combat_system.go`)
- Initiates combat sessions (`CombatSession`) when a player attacks a hostile NPC or when allies in the same room participate.
- Calculates turn order by sorting fighters' initiative attributes (`sortTurnsOrderByInitiative`).
- Executes combat turns (attacks, defend stance, flee attempts) and processes turn timeout callbacks via system channels.
- Resolves victory (rewards & quest updates; players who fled receive no rewards) or defeat (respawning at entrance with 50% max HP).

### 3. Quest Engine (`handle_commands.go`)
- Tracks active player quests and objectives (NPC target defeat or item possession).
- Auto-evaluates completion criteria via `refreshQuestProgress` upon quest commands or combat victory.
- Grants rewards (items, key items) directly to player inventory.

### 4. Group Engine (`group.go`)
- Allows players to form adventuring parties (`GROUP CREATE`, `INVITE`, `JOIN`, `KICK`, `PROMOTE`, `LEAVE`).
- Shares combat entry for party members present in the same room.
- Facilitates group-wide notifications and dedicated group chat streams.
