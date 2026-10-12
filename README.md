*This project has been created as part of the 42 curriculum by crappo, aluslu.*

# TAP - The Answer Protocol

A multiplayer retro text and graphical adventure game (MUD/RPG) built in **Go**. *The Answer Protocol (TAP)* features a high-performance concurrent TCP server, a decoupled thread-safe game engine, a customizable terminal interface (TUI), and an interactive 2D graphical visual client (GUI).

---

## Sub-Component Documentation

For in-depth documentation on individual modules, please consult the dedicated README files:
- [Server](./src/server/README.md)
- [Engine](./src/engine/README.md)
- [Client Logic](./src/client/controller/README.md)
- [TUI Client](./src/client/tui/README.md)
- [GUI Client](./src/client/gui/README.md)

---

## Description

**The Answer Protocol (TAP)** brings retro text-adventure mechanics into a modern multiplayer architecture. Players connect to a high-performance TCP server hosting a shared, in-memory world (state resets cleanly upon server restart, strictly following the RFC specification without disk persistence) where they can:
- Explore interconnected rooms and regions forming loops and branches.
- Interact and converse with non-player characters (NPCs) across multiple distinct roles.
- Accept, track, and explicitly complete quests for unique item rewards.
- Manage inventories, inspect items, and consume usable items for health and stat boosts.
- Form adventuring parties with other players, with full party management and shared combat.
- Engage in initiative-driven, turn-based combat against hostile creatures with flee penalties and turn timers.

The project is designed with strict separation of concerns, decoupling connection handling from game state mutation and presentation layers.

---

## Instructions

### Prerequisites
- **Go**: Version 1.26 or higher.
- **C Compiler / X11 / OpenGL headers**: Required for Raylib GUI compilation on Linux (e.g., `libgl1-mesa-dev`, `xorg-dev`).

### Installation & Compilation
1. **Download and tidy dependencies**:
   ```bash
   make install
   ```

2. **Build all binaries**:
   ```bash
   make build
   ```
   *Binaries are compiled to the `exec/` directory (`exec/server`, `exec/tui`, `exec/gui`).*
   *Note: All assets, maps, and configuration files (`world.json`) are embedded within the binaries using Go's `//go:embed`. Executables can be distributed and run standalone.*

3. **Rebuild or Clean**:
   ```bash
   make clean   # Remove compiled binaries in exec/
   make re      # Clean and rebuild all binaries
   ```

### Running the Server & Clients
All executables support command-line arguments to specify host and port. You can pass them directly or via the Makefile `ARGS` variable.

#### 1. Server
- **Makefile**:
  ```bash
  make run-server
  # With custom port:
  make run-server ARGS="--port 8080"
  ```
- **Direct binary**:
  ```bash
  ./exec/server
  ./exec/server --port 8080
  ```
  *(Default port: `8080`)*

#### 2. TUI Client (Terminal Interface)
- **Makefile**:
  ```bash
  make run-client
  # With custom host and port:
  make run-client ARGS="--ip 127.0.0.1 --port 8080"
  ```
- **Direct binary**:
  ```bash
  ./exec/tui
  ./exec/tui --ip 127.0.0.1 --port 8080
  ```
  *(Default host: `localhost`, default port: `8080`)*

#### 3. GUI Client (Graphical Interface)
- **Makefile**:
  ```bash
  make run-client-gui
  # With custom host and port:
  make run-client-gui ARGS="--ip 127.0.0.1 --port 8080"
  ```
- **Direct binary**:
  ```bash
  ./exec/gui
  ./exec/gui --ip 127.0.0.1 --port 8080
  ```
  *(Default host: `localhost`, default port: `8080`)*

---

## CLI Interface Choice

In accordance with **Chapter V.3 (CLI Client)** of the subject, our group implemented a hybrid approach combining both design choices:

1. **Direct RFC Protocol Syntax (CLI Command Line)**:
   The TUI client embeds a dedicated bottom command-line console (`tap-cli> ` in `command_line_panel.go`). Any command entered by the player (e.g., `LOOK`, `MOVE north`, `TAKE mais_sucre`, `ATTACK bob`, `COMPLETE_QUEST effraye_maelis`, `USE marshmallow`) is sent directly to the server using standard RFC protocol syntax over the TCP stream.
2. **Interactive UI Translation (TUI Panels)**:
   Alongside manual RFC command input, the TUI provides user-friendly interactive widgets (Navigation compass, Room interaction list, Inventory manager, Combat dashboard, and Quest tracker). When a user selects a dialogue option, moves with arrow keys, or clicks an action button, the client controller translates the interaction into RFC-compliant packets behind the scenes before dispatching them to the server.

Both approaches operate simultaneously, giving advanced users direct terminal control while offering new players an accessible interface.

---

## Architecture

### Server Design & Dispatcher Model
The server architecture decouples network connection handling from game state processing using Go's concurrent channels:
- **Network Gateway (Server)**: Listens for incoming TCP connections on a dedicated listener goroutine. Each connected client gets lightweight reader and writer goroutines to manage non-blocking socket stream I/O.
- **Thread-Safe Exchanger (`pr.Exchanger`)**: Acts as a thread-safe message bridge between the network layer and game engine, exposing `ServerInput`, `ServerOutput`, `JoinChan`, and `LeaveChan` channels.
- **Single-Threaded Engine Loop**: The game `Engine` runs a single-threaded event loop (`broadcaster`). By processing requests sequentially from `ServerInput`, state updates (combat turns, room movement, inventory modifications, quest validations) occur sequentially without requiring complex mutex locking on game entities.

### Client Architecture
Clients use a decoupled Controller pattern (`src/client/controller`):
- **Controller**: Maintains client-side state caches, parses incoming server responses and asynchronous events (`EVT`), and serializes outgoing protocol commands.
- **Presenters (TUI & GUI)**: Interchangeable frontends (`tview`/`tcell` for TUI, `raylib` for GUI) bind to controller events to render terminal panels or 2D tilemaps.

### System Architecture Diagram

```
┌──────────────┐
│  GUI Client  │
│  (Raylib)    │──┐
└──────────────┘  │
                  │   ┌────────────────┐      TCP       ┌──────────────┐    Go ch.    ┌──────────────┐
                  ├──>│   Controller   │<──────────────>│    Server    │<────────────>│    Engine    │
                  │   │ (Client Logic) │   Connection   │ (TCP server) │              │ (Game Logic) │
┌──────────────┐  │   └────────────────┘                └──────────────┘              └──────────────┘
│  TUI Client  │──┘
│  (Tview)     │
└──────────────┘
```

---

## Protocol Implementation

TAP implements and extends the **RFC 42TAP** specification.

### Framing & Wire Format
- Over the TCP wire, all messages are UTF-8 encoded and terminated with a newline character (`\n`).
- **Standard Status Prefixes**:
  - `OK [json_payload]`: On success, the JSON data (e.g., room information, item arrays, combat status) immediately follows the `OK ` prefix. *Note: There is no intermediate `"datas"` envelope on the network wire; JSON payloads are placed directly after `OK `.*
  - `ERR <code> <MESSAGE>`: Error responses follow standard numeric codes and error identifiers.
  - `EVT <CATEGORY> <TYPE> [data]`: Asynchronous event notifications pushed to relevant clients.

### Added & Extended Commands
To support full RPG mechanics, graphical synchronization, and group play, the following commands were added beyond base RFC 42TAP:

| Command | Arguments | Description & Justification |
| :--- | :--- | :--- |
| `USE` | `<item_id \| item_name>` | Consumes an item from player inventory (`consumable` type), applying health recovery or stat modifications and decrementing charges. |
| `FLEE` | *none* | Attempts to escape the active combat turn. The player leaves the combat session and forfeits all battle loot, even if their group wins. |
| `COMPLETE_QUEST` | `<quest_id>` | Explicitly validates quest prerequisites (target NPC defeat in `DefeatedNpcs` or required item in inventory), consumes target items, grants rewards, and marks the quest as completed (`unavailable`). |
| `INSPECT` | `PLAYER <name>`<br>`NPC <name>`<br>`ITEM <id>`<br>`INVENTORY_ITEM <id>`<br>`SELF` | Returns detailed attributes, descriptions, and statistics for entities to populate GUI and TUI inspection panels. |
| `GROUPED` | *none* | Returns a JSON array of all player names currently in the caller's party. |
| `UNGROUPED` | *none* | Returns a JSON array of all connected players on the server who are not currently members of any group. |
| `COMBAT STATS` | *none* | Queries real-time combat session status, including current turn holder, fighter HP/Shield/Initiative, and round state. |
| `CHAT COMBAT` | `<message>` | Scoped chat channel restricted exclusively to fighters participating in the active combat session. |
| `GROUP` Subcommands | `CREATE`<br>`INVITE <player>`<br>`KICK <player>`<br>`JOIN <group_id>`<br>`LEAVE`<br>`PROMOTE <player>`<br>`ACCEPT`<br>`DECLINE` | Extended group management enabling invitations, party kicks, voluntary leaving, and leadership handovers with accept/decline workflows. |
| `NOTIFY_PLAYER_POSITION`<br>`GET_PLAYER_POSITIONS` | `<x> <y> <direction>` | Synchronizes 2D coordinates and facing directions for smooth real-time sprite rendering in GUI/TUI. Exempted from rate limiting. |
| `NOTIFY_ITEM_POSITION`<br>`GET_ITEM_POSITIONS` | *none* | Synchronizes dynamic 2D ground coordinates of dropped items in rooms. Exempted from rate limiting. |
| `NOTIFY_PLAYER_EMOTE`<br>`GET_PLAYER_EMOTES` | `<emote_id>` | Broadcasts and queries visual emote reactions above player avatars. |

### Added Protocol Events (`EVT`)

| Event Category | Event Header | Description |
| :--- | :--- | :--- |
| **ITEM** | `EVT ITEM DROPPED <id>` | Broadcast to room occupants when a player drops an item. |
| | `EVT ITEM TOOK <id>` | Broadcast to room occupants when an item is picked up. |
| | `EVT ITEM REMOVED <id>` | Broadcast when an item is consumed or removed from the world. |
| | `EVT ITEM_POSITION <json>` | Broadcasts updated ground item positions for graphical clients. |
| **QUEST** | `EVT QUEST COMPLETED <id>` | Global broadcast informing all players that a quest was completed. |
| | `EVT QUEST AVAILABLE <id>` | Notification that a new quest is available. |
| **COMBAT** | `EVT COMBAT FIGHT_STARTED` | Emitted to combatants when a battle commences. |
| | `EVT COMBAT ALLY_TURN <name>` | Broadcast at the start of each turn indicating which fighter must act. |
| | `EVT COMBAT UPDATE <json>` | Broadcast after actions with updated attacker/target HP and damage. |
| | `EVT COMBAT VICTORY <json>` | Emitted on defeating all hostile NPCs with loot distribution details. |
| | `EVT COMBAT DEFEAT` | Emitted when all players fall in battle, triggering respawn. |
| | `EVT GROUP DISTANT_ALLIES_COMBAT_START` | Notifies remote party members that allies entered combat. |
| **GROUP** | `EVT GROUP INVITE <from>`<br>`EVT GROUP JOIN <name>`<br>`EVT GROUP LEAVE <name>`<br>`EVT GROUP KICK <name>`<br>`EVT GROUP PROMOTE <from>`<br>`EVT GROUP PROMOTE ACCEPTED <name>`<br>`EVT GROUP PROMOTE DECLINED <name>`<br>`EVT new_leader=<name>` | Asynchronous party lifecycle notifications pushed to group members. |
| **PLAYER** | `EVT PLAYER_POSITION <json>`<br>`EVT PLAYER_EMOTE <json>` | Broadcasts player movement vectors and emote animations. |

### Added Error Codes

| Error Code | Identifier | Trigger Condition |
| :--- | :--- | :--- |
| `ERR 202` | `NAME_TOO_LONG` | Nickname exceeds maximum allowed character limit during `CONNECT`. |
| `ERR 302` | `INVALID_COORDINATE_FORMAT` | Coordinates passed to position synchronization are malformed. |
| `ERR 302` | `INVALID_EMOTE_TYPE` | Unknown emote identifier specified. |
| `ERR 406` | `QUEST_NOT_COMPLETED` | Prerequisites for `COMPLETE_QUEST` (target NPC or item) are not met. |
| `ERR 408` | `NO_ALLY_ATTACK` | Friendly fire prevented; cannot target party members or allies. |
| `ERR 409` | `IN_COMBAT` | Disallowed action attempted while engaged in combat (e.g., standard movement). |
| `ERR 409` | `NOT_IN_COMBAT` | Combat command (`FLEE`, `COMBAT STATS`, `CHAT COMBAT`) sent while outside combat. |
| `ERR 411` | `CANT_KICK_LEADER` | Attempt to kick the party leader. |
| `ERR 412` | `ALREADY_LEADER` | Attempt to promote a player who is already the group leader. |
| `ERR 413` | `NOT_PROMOTED` | Responding to a leadership promotion without an active offer. |
| `ERR 414` | `ITEM_NOT_USABLE` | Attempt to execute `USE` on a non-consumable item (weapon or resource). |

### Design Choices Justifications (P1 to P3)

- **P1: Turn-Based Combat Mechanics**:
  - *Initiative-Driven Order*: Turn order is computed deterministically via `sortTurnsOrderByInitiative`, sorting fighters by their `Initiative` stat descending. This rewards agility and ensures fair, tactical turn order.
  - *100-Second Turn Timeout*: `MaxTurnSeconds = 100` is enforced by a Go background timer (`time.AfterFunc`). If an active player stalls, the turn automatically advances to keep the game flowing.
  - *Group Ally Auto-Joining*: Party members in the same room enter combat together, reinforcing cooperative play.
  - *Flee Forfeit Rule*: Players fleeing with `FLEE` exit the combat session and forfeit all battle loot/rewards, even if the remaining group members achieve victory, preventing exploit strategies.
  - *Single-Leader Loot Distribution*: Defeat rewards are deeply cloned (`.Clone()`) and distributed once to the party leader, preventing item duplication exploits.
- **P2: Dynamic Item Management & Lifecycle**:
  - *Unique Instance Cloning*: Items are instantiated as deep copies (`.Clone()`) rather than shared pointer references, eliminating data races and unintended side-effects across inventories and rooms.
  - *Coordinate Tracking*: Items dropped in a room retain 2D coordinates for graphical clients while remaining accessible via standard text commands (`TAKE <id|name>`).
  - *Consumable Mechanics (`USE`)*: Usable items modify stats (`hp`, `shield`, etc.) and consume charges, enabling survival and preparation between encounters.
- **P3: Quest Progression & Completion System**:
  - *Explicit Lifecycle States*: Quests use three explicit states: `available` (world registry template), `progress` (in player quest log with progress metric `0/1`), and `unavailable` (completed with progress `1/1`).
  - *Automated Objective Refresh*: `refreshQuestProgress` automatically monitors defeated target NPCs in `player.DefeatedNpcs` and required items in `player.inventory`.
  - *Explicit Claiming via `COMPLETE_QUEST`*: Players must explicitly complete quests with `COMPLETE_QUEST <quest_id>`, which safely removes required items, awards cloned rewards once, and prevents reward duplication.

---

## Combat System

### Turn Mechanics & Initiative
- Combat starts when a player attacks a hostile NPC using `ATTACK [target]`.
- All group members present in the same room automatically join the active `CombatSession`.
- **Initiative Order**: Fighters are sorted in descending order of their `Initiative` attribute (`sortTurnsOrderByInitiative`).
- **Turn Timeout**: Players have 100 seconds (`MaxTurnSeconds`) per turn. A background timer advances the turn if no action is submitted, preventing griefing and stalling.

### Combat Actions
- `ATTACK [target]`: Deals base damage to the specified enemy NPC (or default enemy if omitted).
- `FLEE`: Attempts to escape the battle. The fleeing player exits combat and forfeits all battle loot and quest credit for this fight. If all players flee, the combat session is cancelled.
- `USE [item]`: Consumes a usable item during combat (e.g., healing food/potions) to restore HP.
- `COMBAT STATS`: Displays real-time team HP, opponent HP, shield values, and active turn holder.
- `CHAT COMBAT [msg]`: Scoped chat for party members in the combat session.

### Damage Formula & Shield Mechanics
- **Player Damage Dealt**: `equippedWeapon.Damage + stats.Damage`.
- **NPC Damage Dealt**: `Stats.Damage`.
- **Target Player Defense**:
  - If a player possesses an active `Shield`, incoming damage is absorbed by the shield first:
    `absorbed = min(stats.Shield, incoming_damage)`
    `stats.Shield -= absorbed`
    `remaining_damage = incoming_damage - absorbed`
  - The remaining damage reduces `stats.Hp`.
- **Target NPC Defense**: NPCs do not have active shields in `takeDamage`; incoming damage directly reduces `Stats.Hp`.

### Victory & Defeat
- **Victory**: Grants cloned enemy loot items to the party leader, registers the enemy NPC ID into each participant's `player.DefeatedNpcs`, and updates quest objectives.
- **Defeat**: If all players fall to 0 HP, combat ends in defeat. Players respawn at the central village square (`place_du_village`, "La place de l'An 2000") with 50% max HP (`stats.HpMax / 2`).

---

## Quest System

### Mechanics & Validation
- Quests are defined in `src/world.json` with target requirements, rewards, and descriptions.
- **Quest Lifecycle States**:
  - `available`: Initial template status in world configuration.
  - `progress`: Accepted by player via `QUEST <npc>` (progress initialized to `"0/1"`).
  - `unavailable`: Completed and claimed via `COMPLETE_QUEST <quest_id>` (progress set to `"1/1"`).
- **Progress Tracking (`refreshQuestProgress`)**:
  - *Target Defeat*: Verifies if the target NPC ID is registered in `player.DefeatedNpcs`.
  - *Target Item*: Verifies if the required item ID is present in the player's inventory.
- **Explicit Completion & Rewards**:
  - Quests are **not** automatically completed; players must run `COMPLETE_QUEST <quest_id>`.
  - Upon validation, target items are removed from inventory, the quest is marked `unavailable`, and a cloned reward item is placed into the player's inventory.

### World Quests Reference
1. **`effraye_maelis`** (Quest-giver: `thomas_charbonneau`):
   - *Target*: Defeat hostile NPC `maelis` in `ferme_charbonneau`.
   - *Reward*: `graines` (Seeds).
2. **`trouve_graines`** (Quest-giver: `maya`):
   - *Target*: Bring `graines` to Maya at `verger_de_la_savane`.
   - *Reward*: `carotte_salee` (Salted Carrot).
3. **`bas_bob`** (Quest-giver: `franck`):
   - *Target*: Defeat `bob` in `foret_malicieuse` to obtain `cancoillotte`.
   - *Reward*: `flocons_d_avoine` (Oat Flakes).
4. **`repare_barrage`** (Quest-giver: `roger`):
   - *Target*: Collect `acier` from `barrage` to repair the hydroelectric dam.
   - *Reward*: `fricadelle` (Fricadelle).

---

## World Design

The game world is configured via `src/world.json` and structured as an interconnected room graph forming loops and exploratory branches:

```
                            [fromagerie] <-------> [fontaneilles] <-> [lac_du_barrage]
                      		      ^                       ^                 ^
                      		      |                       |                 |
   [ferme_uzoloise] <-> [verger_de_la_savane] <-> [foret_malicieuse] <-> [barrage]
                                  ^                       ^                 ^
                                  |                       |                 |
                         [ferme_charbonneau] <-> [place_du_village] <-> [riviere]
                                                          ^                 ^
                                                          |                 |
                                                       [jardin]          [camping]
```

### Key Areas & NPCs
- **La place de l'An 2000 (`place_du_village`)**: Central village square connecting key regions and central respawn point. Home to **`guide`** ("Guide du Village", dialogue NPC explaining game mechanics).
- **La ferme des Légumes Charbonneau (`ferme_charbonneau`)**: Home to `thomas_charbonneau` (quest giver) and `maelis` (hostile enemy).
- **La ferme Uzoloise (`ferme_uzoloise`)**: Royal farm of Uzols; home to hostile enemy `le_prince`.
- **Fontaneilles (`fontaneilles`)**: Cheese-scented town; home to hostile enemy `gabinap` (`HardcoreLevelingWarrior`).
- **Le verger de la Savane (`verger_de_la_savane`)**: Fruit orchard; home to `bernard` (dialogue) and `maya` (quest giver).
- **La fromagerie de St Julien (`fromagerie`)**: Cancoillotte shop managed by `franck` (quest giver).
- **Le lac de St Etienne de Cantales (`lac_du_barrage`)**: Artificial lake connecting Fontaneilles and the dam.
- **Le barrage de St Etienne de Cantales (`barrage`)**: Hydroelectric dam maintained by `roger` (quest giver).
- **La forêt malicieuse (`foret_malicieuse`)**: Forest housing creature `bob` (hostile Bob le blob).
- **Le camping des 3 chênes (`camping`)**: Peaceful campsite south of the river.
- **La rivière (`riviere`)**: River location home to `salomon` (dialogue salmon).
- **Le jardin (`jardin`)**: Community garden south of the village square.

### Item Distribution
- **Weapons & Tools**: `pioche`, `branche_solide`, `canne_a_peche`.
- **Consumables**: `mais_sucre`, `carotte_salee`, `cancoillotte`, `flocons_d_avoine`, `fricadelle`, `marshmallow`.
- **Resources**: `graines`, `acier`.

---

## Server Logging

Server logging utilizes Go's structured `log/slog` library with `slog.NewJSONHandler(os.Stdout, ...)` directly emitting flat JSON records to standard output.

### Log Format
Logs use standard `slog` fields (`time`, `level`, `msg`) with flat contextual key-value pairs (no nested `"datas"` wrapper):

**Connection Log:**
```json
{
  "time": "2026-09-27T16:00:00.123456789Z",
  "level": "INFO",
  "msg": "CLIENT_CONNECTED",
  "client_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
  "ip": "127.0.0.1:54321"
}
```

**Command Execution Log:**
```json
{
  "time": "2026-09-27T16:00:01.123456789Z",
  "level": "INFO",
  "msg": "Server response",
  "name": "alice",
  "response": "OK",
  "command": "MOVE north"
}
```

### Events & Abuse Detection
- **Monitored Events**: Connection lifecycle (`CLIENT_CONNECTED`, `CLIENT_DISCONNECTED`), authentication, command dispatches, world mutations, and panic recoveries.
- **Abuse Prevention**: A Token Bucket rate-limiter monitors command frequency per socket. Flooding triggers `WARN_SPAM` (`ERR 902`) warnings, and repeated spam triggers socket teardown (`ERR 902 CONNECTION_CLOSED_DUE_TO_SPAM`).
- **Zombie Sockets**: Idle connections without authentication or activity are automatically timed out and pruned.

---

## Testing

The codebase includes comprehensive integration, network, and scenario test suites:

### Running Tests
To run all tests:
```bash
make test
```
To run tests with the Go race detector enabled:
```bash
make test-race
```
To run full strict verification (lint + tests + race detection):
```bash
make check
```

### Test Suite Coverage
- **Network Resilience (`src/tests/network/`)**:
  - Connection acceptance & max connection limits (`max_player_limit_test.go`).
  - Rate limiting & anti-flood detection (`anti_flood_test.go`).
  - TCP payload fragmentation & coalescing (`tcp_fragment_coalescing_test.go`).
  - Huge payload rejection (>1024 bytes with `ERR 400 INVALID_COMMAND`) (`tcp_huge_payload_test.go`).
  - Zombie client timeouts and graceful server shutdowns.
- **Scenario Integration (`src/tests/scenarios/`)**:
  - Combat initiation, turn order, and victory/defeat resolution (`attack_scenarios.go`).
  - Multi-user concurrency and movement race condition prevention (`concurrency_scenarios.go`).
  - Group creation, invitations, kicks, and leadership promotions (`group_scenarios.go`).
  - Item pick-up, dropping, using, and quest completion flows (`item_scenarios.go`).
- **Memory & Resource Leaks (`src/tests/leaks/`)**:
  - Clean client disconnection and TUI graceful shutdown without goroutine leaks (`normal_disconnection_test.go`).

---

## Group Contributions

- **`crappo` (Clément Rappo)**:
  - Server architecture Implementation
  - Adaptive TUI interface
  - GUI interface system using Tiled application
  - Embed system to add utils/assets files to binary files

- **`aluslu` (Ali Uslu)**:
  - Complete organization of the project test
  - Rework of the Client logic allowing multiple interfaces for a unique client logic
  - Implementation of many systems (rate limiter, GUI notify system...)
  - Protocol package implementation defining common datas structures, commands...

---

## Resources

### References
- [Go Documentation](https://golang.org/doc/): Standard library reference for `net`, `sync`, channels, and `log/slog`.
- [Raylib Go Bindings](https://github.com/gen2brain/raylib-go): Documentation for 2D graphics rendering.
- [tview & tcell](https://github.com/rivo/tview): Terminal UI widget framework for Go.
- [RFC 42TAP Specification](./utils/protocol-rfc.html): Standard guidelines for text-based multiplayer adventure protocols.

### AI Usage
During this project, AI was used in order to understand key concepts like:
- Channels/goroutines concurrency patterns in Go.
- `tview`/`tcell` event models and widget hierarchies.
- Raylib 2D rendering and camera helpers.

It was also used to help us:
- Brainstorm and refine decoupled software architectures (exchanger bridge).
- Perform code reviews and identify race conditions.
- Structure and format technical documentation.
