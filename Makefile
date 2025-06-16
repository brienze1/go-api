OS=linux
ARCH=amd64
GO111MODULE=on
GOENV=GOOS=${OS} GOARCH=${ARCH}
TEST_THRESHOLD_COVERAGE := 40
MUTATION_THRESHOLD_COVERAGE := 45
MUTATION_THRESHOLD_EFFICACY=100

scenarios=all

install:
	go mod download -x
	go mod tidy

build-local:
	go mod tidy
	go build -o go-api -ldflags="-s -w" ./cmd/server/main.go

docker-up:
	./scripts/docker-up.sh

docker-down:
	./scripts/docker-down.sh

run:
	go run ./cmd/server/main.go

mocks:
	go install go.uber.org/mock/mockgen
	go get go.uber.org/mock/mockgen
	mockgen -source=internal/domain/adapters/create_note_repository_adapter.go -destination=test/mocks/create_note_repository_adapter_mock.go -package=mocks
	mockgen -source=internal/integration/queues/notes_queue.go -destination=test/mocks/notes_queue_mock.go -package=mocks

unit-test:
	go test -v ./test/unit/...

unit-test-coverage:
	go test ./test/unit/... -covermode=count -coverpkg=./internal/...,./cmd/...,./pkg/... -coverprofile ./coverage.out
	go tool cover -func coverage.out
	go tool cover -html ./coverage.out

bdd-test:
	go test -v ./test/integration/... --scenarios=$(scenarios)

bdd-test-coverage:
	go test ./test/integration/... -covermode=count -coverpkg=./internal/...,./cmd/...,./pkg/... -coverprofile ./coverage.out
	go tool cover -func coverage.out
	go tool cover -html ./coverage.out

test-coverage:
	go test -coverpkg=./internal... -coverprofile=coverage.out ./test/integration... ./test/unit...
	go tool cover -func coverage.out
	go tool cover -html ./coverage.out

validate-coverage:
	@echo "Validating coverage..."
	coverage=$$(go tool cover -func=coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	threshold=$(TEST_THRESHOLD_COVERAGE); \
	coverage_num=$${coverage%.*}; \
	if [ "$$coverage_num" -lt "$$threshold" ]; then \
		echo "❌ Test coverage ($$coverage_num%) is below the required threshold ($$threshold%)."; \
		exit 1; \
	else \
		echo "✅ Test coverage meets the required threshold ($$coverage_num%)/($$threshold%)."; \
	fi
	@echo "Coverage validation completed!"

bench-test:
	go test -v ./test/benchmark/... -bench .  -benchmem -count=10 | tee benchmark.txt

fuzzy-test:
	go test -v ./test/fuzzy/...  -fuzz=. -fuzztime 10s

queues-benchmark-profile:
	go test -v ./test/benchmark/internal/integration/queues/... -bench .  -benchmem -count=10 -memprofile queues-benchmark-mem.out -cpuprofile queues-benchmark-cpu.out -o queues-benchmark-profile.out  

queues-benchmark-mem-pprof:
	go tool pprof -http=:8080 queues-benchmark-mem.out 

queues-benchmark-cpu-pprof:
	go tool pprof -http=:8080 queues-benchmark-cpu.out

mutation-test:
	@echo "Installing Gremlins..."
	go get github.com/go-gremlins/gremlins/cmd/gremlins
	go install github.com/go-gremlins/gremlins/cmd/gremlins
	go mod tidy
	@echo "Gremlins installed successfully!"
	@echo "Running mutation tests..."
	gremlins unleash \
              --silent=false \
              --integration=true \
              --dry-run=false \
              --tags="" \
              --output="" \
              --workers=4 \
              --test-cpu=2 \
              --timeout-coefficient=0 \
              --threshold-efficacy=90 \
              --threshold-mcover=50 \
              --coverpkg "./internal/...,./pkg/..." \
		   	  --exclude-files "test/mock/..." --exclude-files "test/mocks/..." \
              --arithmetic-base=true \
              --conditionals-boundary=true \
              --conditionals-negation=true \
              --increment-decrement=true \
              --invert-assignments=true \
              --invert-bitwise=true \
              --invert-bwassign=true \
              --invert-negatives=true \
              --invert-logical=true \
              --invert-loopctrl=true \
              --remove-self-assignments=true | tee mutation_results.log
	@echo "Mutation tests completed!"

validate-mutation-coverage:
	@echo "Validating mutation coverage..."
	coverage=$$(grep -E 'Mutator coverage: [0-9]+(\.[0-9]+)?%' mutation_results.log | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	threshold=$(MUTATION_THRESHOLD_COVERAGE); \
	result=$$(echo "$$coverage >= $$threshold" | bc); \
	if [ "$$result" -eq 0 ]; then \
		echo "❌ Mutation test coverage ($$coverage%) is below the required threshold ($$threshold%)."; \
		exit 1; \
	else \
		echo "✅ Mutation test coverage meets the required threshold ($$coverage%)/($$threshold%)."; \
	fi
	@echo "Mutation coverage validation completed!"

validate-test-efficacy:
	@echo "Validating test efficacy..."
	efficacy=$$(grep -E 'Test efficacy: [0-9]+(\.[0-9]+)?%' mutation_results.log | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	threshold=$(MUTATION_THRESHOLD_EFFICACY); \
	result=$$(echo "$$efficacy >= $$threshold" | bc); \
	if [ "$$result" -eq 0 ]; then \
		echo "❌ Mutation test coverage ($$efficacy%) is below the required threshold ($$threshold%)."; \
		exit 1; \
	else \
		echo "✅ Mutation test efficacy meets the required threshold ($$efficacy%)/($$threshold%)."; \
	fi
	@echo "Mutation coverage validation completed!"

dependency-check:
	@echo "Installing Nancy..."
	go install github.com/sonatype-nexus-community/nancy@latest
	go mod tidy
	@echo "Nancy installed successfully!"
	@echo "Running dependency check..."
	go list -json -deps ./... | nancy sleuth | tee nancy-report.log
	@echo "Dependency check completed!"

autofix-dependency-check:
	@echo "Running dependency check fixes..."
	findings_file="nancy-report.log"; \
	grep -Eo 'pkg:golang/[a-zA-Z0-9._/-]+@[v0-9]+\.[0-9]+\.[0-9]+' $$findings_file | while read -r line; do \
		module=$$(echo $$line | sed 's/pkg:golang\///' | awk -F'@' '{print $$1}'); \
		echo "Updating $$module to the latest version..."; \
		go get -u "$$module"; \
	done; \
	go mod tidy
	@echo "Dependency check fix completed!"

# Example of vulnerable dependency: github.com/redis/go-redis/v9 v9.6.1

