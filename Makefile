OS=linux
ARCH=amd64
GO111MODULE=on
GOENV=GOOS=${OS} GOARCH=${ARCH}

scenarios=all

install:
	go mod download -x
	go mod tidy

build-local:
	go mod tidy
	go build -o go-api -ldflags="-s -w" ./cmd/server/main.go

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
	go tool cover -html ./coverage.out

bdd-test:
	go test -v ./test/integration/... --scenarios=$(scenarios)

bdd-test-coverage:
	go test ./test/integration/... -covermode=count -coverpkg=./internal/...,./cmd/...,./pkg/... -coverprofile ./coverage.out
	go tool cover -html ./coverage.out

test-coverage:
	go test ./test/integration... ./test/unit... -covermode=count -coverpkg=./internal/...,./cmd/...,./pkg/... -coverprofile ./coverage.out
	go tool cover -html ./coverage.out

docker-up:
	./scripts/docker-up.sh

docker-down:
	./scripts/docker-down.sh

mutant-test:
	go get github.com/go-gremlins/gremlins/cmd/gremlins
	go install github.com/go-gremlins/gremlins/cmd/gremlins
	gremlins unleash --config=gremlins.yaml --exclude-files "test/mock/..." --exclude-files "test/mocks/..."

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