migrate-up:
	@echo "🚀 Running migration UP..."
	go run ./migration/run.go up

migrate-down:
	@echo "↩️ Running migration DOWN..."
	go run ./migration/run.go down