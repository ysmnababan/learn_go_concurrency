# Batch Insert Benchmark

This project benchmarks **single-row inserts** vs **batch inserts** in Go to measure performance differences in terms of time, memory, and allocations.

---

## 🧪 Benchmark Setup

### Commands

```bash
# Benchmark single inserts (each operation inserts one record)
go test -benchmem -run=^$ -bench ^BenchmarkInsertSingle$ batchinsert -benchtime=60s

# Benchmark batch inserts (each operation inserts multiple records per transaction)
go test -benchmem -run=^$ -bench ^BenchmarkInsertBatch$ batchinsert -benchtime=60s
````

### Notes

* `-run=^$` skips all regular tests (runs only benchmarks).
* `-benchmem` includes memory allocation metrics (`B/op` and `allocs/op`).
* `-benchtime=60s` runs each benchmark for at least 60 seconds to collect stable results.
* `-bench ^BenchmarkInsertSingle$` and `-bench ^BenchmarkInsertBatch$` use regex to match the benchmark function names.

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

## 📈 Summary

| Mode              | Avg Time / op | B/op   | Allocs/op | Iterations | Relative Speed |
| ----------------- | ------------- | ------ | --------- | ---------- | -------------- |
| **Single Insert** | ~11.52 s      | 113 MB | 1.32 M    | 6          | 1×             |
| **Batch Insert**  | ~0.87 s       | 103 MB | 760 K     | 93         | ~13× faster    |

---

## 🧩 Conclusions

* **Batch inserts** significantly reduce execution time compared to inserting rows individually.
* Memory allocations are also lower per operation, although still substantial.
* This confirms that batching is far more efficient for database-heavy workloads.

---

## 🔍 Future Improvements

* Profile using CPU and memory profiles:

  ```bash
  go test -run ^$ -bench ^BenchmarkInsertBatch$ -cpuprofile cpu.prof -memprofile mem.prof
  go tool pprof -http=:8080 ./batchinsert.test cpu.prof
  ```
* Experiment with varying batch sizes (e.g. 10, 100, 1000 rows per batch).
* Optimize allocations by reusing buffers and prepared statements.
* Add more realistic dataset and transaction handling.

---

## 📚 References

* [Go testing package (Benchmarks)](https://pkg.go.dev/testing#hdr-Benchmarks)
* [Go pprof tool](https://pkg.go.dev/runtime/pprof)
* [Effective Go: Benchmarking](https://go.dev/doc/effective_go#benchmarking)

---

> 🧭 **Tip:** Keep long-running benchmarks isolated and avoid running them during normal CI test runs. Use `-run ^$ -bench .` explicitly when you want to measure performance.

