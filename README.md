Go-Sync
------
[![Build Status](https://travis-ci.org/Redundancy/go-sync.svg?branch=master)](https://travis-ci.org/Redundancy/go-sync)
[![GoDoc](https://godoc.org/github.com/Redundancy/go-sync?status.svg)](https://godoc.org/github.com/Redundancy/go-sync)

gosync is a library inspired by zsync and rsync. The intent is that it's easier to build upon than the zsync/rsync codebases. By writing it in Go, it's easier to create in a way that's cross-platform, can take advantage of multiple CPUs with built in benchmarks, code documentation and unit tests.

There are many areas that benefit from the use of multiple threads & connections:
* Making use of multiple http connections, to avoid the bandwidth latency product limiting transfer rates to remote servers
* Making use of multiple CPUs while doing the comparisons

gosync includes a commandline tool that's intended as a starting point for building the functionality that you want.

Zsync modified rsync to allow it to be used against a "dumb" http server, but we can go further:
* Arrange files by hashed path, and checksum: if a file hasn't changed, you can serve the existing version (NB: this works well with s3 sync)
* Split the checksum blocks from the data: serve checksum blocks securely over https, while allowing caching on the data over http
* Split up the data files: improve the fallback when there's a non HTTP 1.1 proxy between the client and server 

### Current State

The commandline tools are fleshed out enough for testing comparison behaviour and profiling it against real input files.
There isn't yet an http block source implemented, but that should probably be the next thing.

### Performance
On an 8 MB file with few matches, I'm hitting about 16 MB/s with 4 threads. I think that we're mostly CPU bound, and should scale reasonably well with more processors.

When there are very similar files, the speed is far higher.

#### Some numbers:
Budget for 8 MB/s byte by byte comparison on single thread: 120ns

Current Benchmark State:
- Checksum: 62.7 ns
- Comparison (No index lookup)
  - Weak Rejection: 85.2 ns
  - Strong Rejection: 458 ns (MD5: 391 ns)
- Index Lookup: 70 ns + (for reject)

The 32 bit Rollsum hash produces far fewer false positives than the 16 bit one, with the same 4 byte hash overhead.

Index generation:
- Easily hits 100 MB/s on my workstation, satisfying the idea that you should be able to build 12 GB payloads in ~1m 
- More likely to be bottlenecked by the disk throughput / seek time than the CPU

### TODO
- [ ] Provide more helpers to make common usage easier (multi-threading etc)
- [ ] Clean up naming consistency and clarity: Block / Chunk etc
- [ ] Flesh out full directory build / sync
- [ ] Implement 'patch' payloads from a known start point to a desired endstate
- [ ] Validate full file checksum after patching
- [ ] Provide bandwidth limiting / monitoring as part of http blocksource
- [ ] Think about turning the filechecksum into an interface
- [ ] Avoid marshalling / unmarshalling blocks during checksum generation
- [ ] Sequential patcher to resume after error?

### Testing

#### Unit tests

    go test github.com/Redundancy/go-sync/...

#### Commandline & files

	go build github.com/Redundancy/go-sync/gosync
	gosync build filenameToPatchTo
	gosync patch filenameToPatchFrom filenameToPatchTo.gosync