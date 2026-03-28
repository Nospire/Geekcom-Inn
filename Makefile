.PHONY: build server client dev lint fmt vet test test-race ci clean

# ── Build ──────────────────────────────────────────
build: server client

server:
	go build -o bin/tavrn-admin ./cmd/tavrn-admin

client:
	go build -o bin/tavrn ./cmd/tavrn

# ── Dev ────────────────────────────────────────────
# Terminal 1: make run
# Terminal 2: make connect (or: ssh localhost -p 2222)
run: server
	./bin/tavrn-admin

connect: client
	./bin/tavrn --dev

dev: server
	@echo "Starting server on :2222 — connect with: make connect"
	./bin/tavrn-admin

# ── Lint (mirrors CI) ─────────────────────────────
fmt:
	gofmt -w .

lint: fmt
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Files not formatted:"; echo "$$unformatted"; exit 1; \
	fi
	go vet ./...

# ── Test ───────────────────────────────────────────
test:
	go test ./internal/... ./ui/...

test-race:
	go test -race ./internal/... ./ui/...

# ── CI (run before push) ──────────────────────────
ci: lint build test-race
	@echo "All checks passed."

# ── Clean ──────────────────────────────────────────
clean:
	rm -rf bin/ tavrn.db
