default:
  @just --list

test:
  gotestsum --format testname -- -v -p 10 -tags=integration ./...
