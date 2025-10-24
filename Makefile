build:
	go build -o todo-mcp main.go

run: build
	./todo-mcp
