# TUI Client

The Terminal User Interface (TUI) component provides a terminal-based interface for The Answer Protocol (TAP). It is written in Go and uses the `github.com/rivo/tview` framework along with `github.com/gdamore/tcell/v2` for low-level terminal event handling and cell drawing.

---

## Architectural Overview

The TUI architecture separates visual presentation from game logic:

- **Application (`app.go`)**: Manages the main event definition, layout initialization, screen redraws, and global key interceptors.
- **View Factory (`view_factory.go`)**: Instantiates widgets, configures grid/flex layout hierarchies, and binds UI panels to the central layout manager.
- **Focus Manager (`move_focus.go`)**: Handles panel focus transitions using keyboard controls (`Tab`, `Shift+Tab`, directional keys).
- **Theme Manager (`theme.go` & `style_utils.go`)**: Defines color palettes, box borders, highlight styles, and typography settings across all terminal widgets.
- **Controller Binding (`tui_client.go`)**: Connects the TUI front-end to `client/controller`, registering asynchronous event callbacks that trigger UI redraws (`app.Draw()`).

---

## Comprehensive Panel Reference

The TUI divides the screen into specialized interactive panels managed by a `tview.Pages` and `tview.Flex` layout structure:

### Connection Panel (`connect_panel.go`)
Modal prompt displayed on launch. Collects server address, port, and player pseudonym. Validates input fields and initiates the TCP connection through the controller before revealing the main application grid.

### Navigation Panel (`navigation_panel.go`)
Renders compass exit indicators (North, South, East, West). Indicates available directions based on current room data and allows quick movement via directional keyboard shortcuts or button clicks.

### Room View & Interaction Panel (`interaction_panel.go`)
Displays room title, detailed description, visible items on the ground, present NPCs, and active players in the same location. Allows selecting NPCs to initiate conversation or selecting ground items to pick up. NPC dialogue responses are dynamically displayed underneath their respective names in gray.

### Combat Panel (`combat_panel.go`)
Dedicated dashboard active during combat sessions:
- **Turn Order Header**: Highlights current active turn holder and remaining turn time.
- **Fighter Status Lists**: Displays HP/MaxHP progress bars, status effects, and shield values for both allies and enemy NPCs.
- **Action Controls**: Interactive buttons for `ATTACK`, `FLEE`, and `USE` item.
- **Combat Log**: Real-time scrolling combat activity log.

### Quest Panel (`quest_panel.go`)
Lists player active and completed quests. Displays quest titles, descriptions, target NPCs to defeat, required items to collect, and expected rewards.

### Items & Inventory Panel (`items_panel.go`)
Renders player inventory. Supports item inspection, equipping (automatically filtering for weapon-type items), dropping items into the current room, and consuming usable items.

### Group Panel (`group_panel.go`)
Displays party information:
- List of group members and their current room locations.
- Group leader status indicator.
- Action controls to create a group, invite nearby players, kick members, promote a new leader, or leave the group.

### Chat & Event Log Panel (`chat_panel.go` & `server_response_panel.go`)
Multi-tabbed communication output:
- **Global Channel**: Public chat across all connected players.
- **Room Channel**: Chat restricted to players in the current room.
- **Group Channel**: Private party messaging.
- **Combat Channel**: Messages restricted to participants in an active combat session.
- **Server Event Log**: Low-level stream of server notifications and system messages.

### Inspector Panel (`inspector_panel.go`)
Modal overlay used to inspect entities (self stats, target players, NPCs, items). Displays raw stats, descriptions, and metadata.

### Choice List & Popup Panels (`choice_list_panel.go` & `popup_panel.go`)
Provides interactive selection lists for dialogue choices with NPCs and error alert overlays for invalid server operations.

### Command Line Panel (`command_line_panel.go`)
Bottom input field allowing direct execution of manual text commands with history navigation (`Up`/`Down` arrows).

---

## State Synchronization & Event Dispatching

When the server sends responses or asynchronous events (`EVT`), the controller updates its local state cache and triggers event listeners registered in `tui_client.go`. These callbacks run on the application thread using `app.QueueUpdateDraw()`, ensuring thread-safe terminal UI re-rendering without flickering or data corruption.
