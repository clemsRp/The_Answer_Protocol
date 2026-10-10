# GUI Client

The Graphical User Interface (GUI) component provides a 2D visual visual client for The Answer Protocol (TAP). It is developed in Go using Raylib bindings via `github.com/gen2brain/raylib-go/raylib`.

---

## Architectural Structure

The GUI codebase is organized into modular packages under `src/client/gui/src/`:

- **Core Package (`core/`)**: Coordinates application initialization (`app.go`), state updates (`update.go`), combat flow management (`combat.go`), item/inventory data synchronization (`items.go`), NPC interactions (`npcs.go`), and configuration options (`options.go`).
- **Drawer Package (`drawer/`)**: Implements the rendering pipeline, handling tilemap rendering, sprite animations, HUD elements, combat overlays, chat boxes, and dialogue menus.
- **Updater Package (`updater/`)**: Handles frame-by-frame state mutation, keyboard/mouse input processing, collision detection, and coordinate interpolation.
- **UI Package (`ui/`)**: Reusable UI widget system managing interactive buttons (`button.go`), selection dropdowns (`select.go`), dialogue overlays (`interaction.go`), options toggles (`option.option.go`), and emote indicators (`emote.go`).
- **Variables Package (`variables/`)**: Maintains central layout dimensions, font assets, color schemes, texture references, and panel state flags.

---

## Rendering Pipeline

The main draw loop in `drawer/draw.go` executes render passes in a strict sequence:

1. **Clear Screen & Camera Matrix**: Clears the frame buffer and applies the 2D world camera transform.
2. **Tilemap & Environment Layer (`drawer/map.go`)**: Renders ground tiles, background decorations, and wall structures loaded from JSON map assets (`src/client/gui/maps/`).
3. **Entities & Sprites Layer (`drawer/player.go` & `drawer/selects.go`)**:
   - Renders player avatars with animated walking cycles and directional sprite orientations.
   - Renders nearby NPCs, enemies, and ground items.
   - Displays player names, health bars, and active status indicators above avatars.
4. **Collision & Trigger Layer**: Evaluates room exit portals and interactive entity collision boxes (`updater/collisions.go`).
5. **HUD & User Interface Overlay**:
   - **Game HUD (`drawer/game.go`)**: Displays room titles, player health represented by dynamic heart icons, mini-map indicators, and quick action bars.
   - **Chat Window (`drawer/chat.go`)**: Renders scrollable chat overlays with channel selection tabs (Global, Room, Group, Combat).
   - **Group Panel (`drawer/group.go`)**: Party member list and management overlay.
   - **NPC Dialogue Menu (`drawer/talk.go`)**: Dialogue box with selectable conversation options.
   - **Inspector Overlay (`drawer/inspect.go`)**: Detailed modal displaying entity attributes and item descriptions.
6. **Combat Overlay (`drawer/combat*.go`)**: When entering combat, renders a specialized battle interface:
   - Fighter status cards for team members and opponents (`combat_fighters.go`).
   - Turn order carousel and current turn timer indicator (`combat.go`).
   - Action command buttons: Attack, Defend, Flee, Use Item (`combat_actions.go`).
   - Dedicated combat event chat log (`combat_chat.go`).

---

## Real-Time Position Synchronization

To support smooth 2D movement in a multiplayer environment, the GUI implements a position synchronization mechanism:

- **Movement Input & Collision (`updater/player.go` & `updater/collisions.go`)**: Captures WASD or Arrow key inputs, checks tile collision masks, and updates local player coordinates.
- **Position Broadcast**: Periodically sends `NOTIFY_PLAYER_POSITION` commands to the server containing grid coordinates and movement direction.
- **Position Polling & Interpolation**: Requests position updates for other players in the same room using `GET_PLAYER_POSITIONS`. Smoothly interpolates foreign player sprite positions across rendered frames to eliminate jitter.

---

## Resource & Asset Management

All assets were find on itch.io. It's the **Sprout Lands** packages:
- [Basic assets](https://cupnooble.itch.io/sprout-lands-asset-pack)
- [UI assets](https://cupnooble.itch.io/sprout-lands-ui-pack)
- [Addiational assets](https://cupnooble.itch.io/sorry-package)

All maps were manuelly made on the Tiled app using those assets.

Assets are loaded during startup and are now embedded directly in the binary using `//go:embed`. This means you can distribute the standalone executable without the asset folders!
- **Tilesets & Maps**: Embedded map JSONs and tile textures defining room geometries.
- **Sprite Sheets**: Embedded character movement frames, NPC artwork, item icons, and action button textures.
- **Fonts & Shaders (`drawer/shader/`)**: Custom TrueType fonts and optional post-processing GLSL shaders for visual effects.
