#!/usr/bin/env bash

go test -bench=BenchmarkPDQSortInt1K -cpuprofile cpu.prof -memprofile mem.prof
go tool pprof -svg cpu.prof > cpu.svg
go tool pprof -svg mem.prof > mem.svg
