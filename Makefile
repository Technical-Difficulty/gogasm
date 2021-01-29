
help:
	@echo ""
	@echo "usage: make COMMAND"
	@echo ""
	@echo "Commands:"
	@echo "  lines"
	@echo "  server.run"
	@echo ""

# never ever go over 1000 lines!
lines:
	@find . -name '*.go' | xargs wc -l | sort -nr

server.run:
	@go run server/server.go
