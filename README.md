# Batch Insert Benchmark

This project benchmarks **single-row inserts**, **batch inserts**, and **concurrent single-row inserts using goroutines** in Go to measure performance differences in terms of time, memory, and allocations.

---

## 🧪 Benchmark Setup

### Commands

```bash
# Benchmark single inserts (each operation inserts one record)
go test -benchmem -run=^$ -bench ^BenchmarkInsertSingle$ batchinsert -benchtime=60s

# Benchmark batch inserts (each operation inserts multiple records per transaction)
go test -benchmem -run=^$ -bench ^BenchmarkInsertBatch$ batchinsert -benchtime=60s

# Benchmark single inserts using multiple goroutines
go test -benchmem -run=^$ -bench ^BenchmarkInsertGoroutine$ batchinsert -count=5
````

### Notes

* `-run=^$` skips all regular tests (runs only benchmarks).
* `-benchmem` includes memory allocation metrics (`B/op` and `allocs/op`).
* `-benchtime=60s` runs each benchmark for at least 60 seconds to collect stable results.
* `-bench ^BenchmarkInsertSingle$`, `^BenchmarkInsertBatch$`, and `^BenchmarkInsertGoroutine$` use regex to match the benchmark function names.

---

## 🧠 Environment

| Field            | Value                                     |
| :--------------- | :---------------------------------------- |
| **OS**           | Linux                                     |
| **Architecture** | amd64                                     |
| **CPU**          | Intel(R) Core(TM) i5-10310U CPU @ 1.70GHz |

---

## 📊 Benchmark Results

### Single Insert

```
BenchmarkInsertSingle-8                6    11524933202 ns/op    113411536 B/op    1325805 allocs/op
PASS
ok   batchinsert   80.179s
```

**Interpretation:**

* Each single insert operation took ~11.52 seconds on average.
* Allocated ~108 MB per operation.
* Performed ~1.32 million allocations per operation.
* Only 6 iterations were completed in 60 seconds (since each op is heavy).

---

### Batch Insert

```
BenchmarkInsertBatch-8                93    868565534 ns/op    102641359 B/op    760770 allocs/op
PASS
ok   batchinsert   81.692s
```

**Interpretation:**

* Each batch insert operation took ~0.87 seconds on average.
* Allocated ~97 MB per operation.
* Performed ~760k allocations per operation.
* 93 iterations were completed — much faster throughput than single inserts.

---

### Single Insert with Goroutines

```
BenchmarkInsertGoroutine-8
    1    1263208642 ns/op    116398528 B/op    1332710 allocs/op
    1    1386003834 ns/op    116216832 B/op    1332056 allocs/op
    1    1249484892 ns/op    116228352 B/op    1332137 allocs/op
    1    1649092937 ns/op    116224696 B/op    1332106 allocs/op
    1    1386319101 ns/op    116220736 B/op    1332135 allocs/op
PASS
ok   batchinsert   7.122s
```

**Interpretation:**

* Each concurrent insert operation took ~1.26–1.65 seconds per iteration.
* Memory usage per operation is ~111 MB with ~1.33 million allocations.
* Using multiple goroutines improved throughput compared to single-row sequential inserts (~9× faster per wall-clock time) while still maintaining per-operation allocation overhead.
* Variability is observed due to scheduling and concurrent DB access.

---

## 📈 Summary

| Mode                           | Avg Time / op | B/op   | Allocs/op | Iterations    | Relative Speed |
| ------------------------------ | ------------- | ------ | --------- | ------------- | -------------- |
| **Single Insert**              | ~11.52 s      | 113 MB | 1.32 M    | 6             | 1×             |
| **Batch Insert**               | ~0.87 s       | 103 MB | 760 K     | 93            | ~13× faster    |
| **Single Insert + Goroutines** | ~1.38 s       | 111 MB | 1.33 M    | 5 (per count) | ~8–9× faster   |

---

## 🧩 Conclusions

* **Batch inserts** remain the most efficient for throughput and memory usage.
* **Concurrent single-row inserts** significantly reduce wall-clock time but allocations per operation remain high.
* Using goroutines can help improve performance without changing the database batch logic, but batching + concurrency may offer the best performance.

---

## 🔍 Future Improvements

* Profile using CPU and memory profiles:

  ```bash
  go test -run ^$ -bench ^BenchmarkInsertBatch$ -cpuprofile cpu.prof -memprofile mem.prof
  go tool pprof -http=:8080 ./batchinsert.test cpu.prof
  ```
* Experiment with varying batch sizes and concurrent workers.
* Optimize allocations by reusing buffers and prepared statements.
* Add more realistic dataset and transaction handling.

---

## 📚 References

* [Go testing package (Benchmarks)](https://pkg.go.dev/testing#hdr-Benchmarks)
* [Go pprof tool](https://pkg.go.dev/runtime/pprof)
* [Effective Go: Benchmarking](https://go.dev/doc/effective_go#benchmarking)

---

> 🧭 **Tip:** Keep long-running benchmarks isolated and avoid running them during normal CI test runs. Use `-run ^$ -bench .` explicitly when you want to measure performance.

