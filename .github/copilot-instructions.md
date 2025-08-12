# Raft Library - Copilot Instructions

**ALWAYS follow these instructions first and only fallback to additional search and context gathering if the information here is incomplete or found to be in error.**

This is a production-grade Go library implementing the Raft consensus algorithm. It's used by etcd, Kubernetes, Docker Swarm, CockroachDB, TiDB, and other major distributed systems. The library follows a minimalistic design philosophy, implementing only the core Raft algorithm while leaving network and disk I/O to users.

## Key Repository Structure

- **Root package (`go.etcd.io/raft/v3`)**: Core Raft algorithm implementation
- **`raftpb/`**: Protocol Buffer definitions for Raft messages and state
- **`rafttest/`**: Testing utilities and interaction-based test framework
- **`tracker/`**: Progress tracking for followers and flow control
- **`quorum/`**: Quorum calculations and voting logic  
- **`confchange/`**: Configuration change handling (add/remove nodes)
- **`tools/mod/`**: Development tools and dependencies

## Build and Test Commands

### CRITICAL: Set Long Timeouts for Build Commands
**NEVER CANCEL** builds or tests early. Always set timeouts of 90+ minutes for any build operation.

### Environment Setup
```bash
# Ensure Go 1.24.5+ is installed (check .go-version file)
go version  # Must be go1.24.5 or later

# Required for verification steps:
export PATH=/usr/local/bin:$PATH:$(go env GOPATH)/bin

# Install protobuf compiler 3.20.3 (EXACT version required):
cd /tmp
wget https://github.com/protocolbuffers/protobuf/releases/download/v3.20.3/protoc-3.20.3-linux-x86_64.zip
unzip protoc-3.20.3-linux-x86_64.zip
sudo cp bin/protoc /usr/local/bin/
sudo cp -r include/* /usr/local/include/

# Install golangci-lint v2+ (config requires v2):
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

### Build and Test
```bash
# Build all packages (very fast - under 1 second):
go build ./...

# Run full test suite - NEVER CANCEL, takes ~25 seconds:
make test
# OR directly: PASSES="unit" ./scripts/test.sh

# Run quick tests with race detection - takes ~4 seconds:
go test -race -short ./...

# Run verification suite - NEVER CANCEL, takes ~15 seconds total:
make verify
# Includes: gofmt, dependency check, linting, mod tidy, protobuf generation
```

### Individual Verification Steps
```bash
# Format check (instant):
make verify-gofmt

# Dependency consistency check (~8 seconds):
make verify-dep

# Lint check (~5 seconds):
make verify-lint

# Go mod tidy check (~1 second):  
make verify-mod-tidy

# Protobuf generation check (~5 seconds):
make verify-genproto
```

## Manual Validation Scenarios

After making changes to the core Raft logic, ALWAYS run these validation tests:

### 1. Leader Election Validation
```bash
# Test leader election in various scenarios:
go test -v -run "TestLeaderElection" -timeout 60s

# Expected: Detailed logs showing candidates becoming leaders, vote exchanges
```

### 2. Log Replication Validation  
```bash
# Test log replication and consistency:
go test ./tracker -v -run "TestProgressBecomeProbe|TestProgressBecomeReplicate" -timeout 60s

# Expected: Progress state transitions, message flows
```

### 3. Core Algorithm Validation
```bash
# Test fundamental Raft operations:
go test -v -run "TestRaftNodes\\|TestRaftFreesReadOnlyMem" -timeout 60s

# Expected: Node initialization, configuration switches, state transitions
```

### 4. Benchmark Validation
```bash
# Ensure performance hasn't regressed:
go test -bench=BenchmarkOneNode -run=^$ -count=3 -timeout 120s

# Expected: Consistent throughput measurements
```

## Working with Protocol Buffers

The `raftpb/` package contains Protocol Buffer definitions. When modifying `.proto` files:

1. **NEVER modify generated `.pb.go` files directly**
2. Edit only `.proto` files in `raftpb/`
3. Run `make verify-genproto` to regenerate (requires protoc 3.20.3 EXACTLY)
4. Generated files should be committed along with proto changes

## Common Development Workflows

### Adding New Raft Features
1. Modify core algorithm in root package
2. Add/update protobuf messages in `raftpb/` if needed
3. Update tracking logic in `tracker/` if needed
4. Add tests using `rafttest/` framework
5. Run full verification: `make verify && make test`

### Debugging Raft Behavior
- Use `rafttest/` interaction framework for controlled scenarios
- Enable verbose test logging: `go test -v -run TestName`
- Check logs for state transitions: `became leader`, `became follower`, `became candidate`

### Performance Testing
- Use existing benchmarks in `*_bench_test.go` files
- Focus on message processing and state machine transitions
- Monitor memory allocation patterns

## CI/CD Integration

The repository uses GitHub Actions with specific test configurations:
- **Linux AMD64**: 4-CPU race tests (longest running)
- **Linux 386**: Single CPU tests  
- **ARM64**: 4-CPU race tests (on etcd-io/raft repo only)

Always run `make verify && make test` locally before committing.

## Key Configuration Files

- **`go.mod`**: Go module definition, requires Go 1.24+ 
- **`.golangci.yaml`**: Linter configuration (requires golangci-lint v2+)
- **`.go-version`**: Specifies exact Go version (1.24.5)
- **`Makefile`**: Standard build targets for verification and testing

## Important Notes

- **This is a library, not an application** - no main packages or executables
- **Network and storage are user responsibility** - library only implements algorithm
- **Deterministic behavior** - same inputs always produce same outputs
- **Production-grade** - powers major distributed systems in production
- **Minimalistic design** - focused solely on core Raft consensus algorithm

## Example Usage Pattern

See `example_test.go` for basic integration pattern:
```go
c := &Config{}
n := StartNode(c, nil)
defer n.Stop()

for {
    rd := <-n.Ready()
    saveStateToDisk(rd.HardState)
    saveToDisk(rd.Entries)
    sendMessages(rd.Messages)
    go applyToStore(rd.CommittedEntries)
    n.Advance()
}
```

Always reference the README.md for detailed usage examples and API documentation.