# Controller

The **Client Controller** forms the intermediate logic and state synchronization layer of *The Answer Protocol* (TAP) client architecture. It decouples network socket communication from user interface frontends (both TUI and GUI).

---

## Architectural Responsibilities

- **Network Manager**: Manages TCP socket connections, outgoing command formatting, and line-delimited message streaming to the server.
- **State Store (`state`)**: Maintains local cached copies of player stats, active inventory, room descriptions, party status, quest objectives, and active combat state.
- **Event Dispatcher (`events.go`)**: Parses incoming `EVT` protocol notifications (room entry/leave, chat broadcasts, combat turns, position updates) and notifies attached UI listeners asynchronously.
- **Response Parser (`responses.go`)**: Unmarshals JSON response structures (`pr.ServerResponse`) returned by the server and updates UI models accordingly.

---

## Interaction Pattern

```
	  [ UI (TUI / GUI) ]
               │
               │  1. Invokes Controller Action (e.g. Move("north"))
               ▼
     [ Client Controller ]
               │
               │  2. Formats & sends raw command over TCP ("MOVE north\n")
               ▼
           [ SERVER ]
               │
               │  3. Parses request & forwards to Engine via Go Channel (ServerInput)
               ▼
        [ GAME ENGINE ]
               │
               │  4. Processes game logic & replies via Go Channel (ServerOutput)
               ▼
           [ SERVER ]
               │
               │  5. Sends TCP response back to client ("OK room=north_bridge ...")
               ▼
     [ Client Controller ]
               │
               │  6. Updates State Store & fires Event Listener
               ▼
      [ UI Redraw Trigger ]
```
