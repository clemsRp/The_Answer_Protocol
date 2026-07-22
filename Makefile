



client:
	go run ./client/tui

server:
	go run ./server

test-fragmented-packets:
	(echo -n "say hello " ; sleep 2 ; echo -n "to eve" ; \
	sleep 2 ; echo "ryone") | nc localhost 8080

test-multiple-packets-in-one-send:
	printf "move north\nattack goblin\nlook around\n" | nc localhost 8080

.PHONY: client server test-fragmented-packets test-multiple-packets
