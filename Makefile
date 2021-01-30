
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
	@find . -name '*.go' ! -path './server/*' ! -name '*_test.go' | xargs wc -l | sort -nr

server.build:
	@cd server && docker build -t gg-server .

server.run:
	@cd server && docker run -p 3000:80 -it --rm gg-server
