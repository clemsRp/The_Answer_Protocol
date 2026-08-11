package engine

import "tap/protocol/parser"

func isNpcInRoom(room parser.Room, npc_name string) bool {
	for _, npc := range room.Npcs {
		if npc_name == npc {
			return true
		}
	}
	return false
}
