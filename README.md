# ⚡ Go HTTP Server

A high-performance, multi-threaded HTTP web server built entirely from scratch in Go. 

This project intentionally bypasses the standard `net/http` library to demonstrate core backend engineering principles. It implements layer-7 HTTP parsing, concurrent request routing, and thread-safe data handling directly on top of raw TCP sockets (`net` package).

---

## 🧠 Architecture Overview

This server demonstrates a deep understanding of system design, concurrency, and memory management:

*   **Custom M:N Worker Pool:** Multiplexes connections across a fixed pool of lightweight Goroutines, minimizing context-switching overhead compared to OS threads.
*   **Buffered Channel Queues:** Incoming TCP connections are pushed into a thread-safe `chan net.Conn` job queue, providing built-in backpressure.
*   **Layer-7 Protocol Parsing:** Custom `bufio` parsers extract HTTP methods and intelligently read `Content-Length` for dynamic POST payloads.
*   **Zero-Allocation Router:** An O(1) hash map routing engine optimized with native string concatenation to completely bypass Go's reflection engine (`fmt`) and eliminate heap allocations on the hot path.

## 📁 Project Structure
```text
micro-http-server/
├── go.mod                  # Go module definition
├── main.go                 # Entry point, initializes router and worker pool
├── server/
│   ├── dispatcher.go       # TCP Listener and Goroutine pool management
│   └── router.go           # O(1) Route matching engine
├── protocol/
│   └── http.go             # Raw byte parsing and HTTP response formatting
└── handlers/
    └── api.go              # Business logic for specific endpoints
```

## 🛠️ Getting Started

### Prerequisites
*   [Go](https://go.dev/dl/) 1.20 or higher installed on your machine.

### Running the Server
Clone the repository and run the entry point. The server will initialize the worker pool and listen on port `8080`.
```bash
git clone [https://github.com/PriyanshuJajal/GO-HTTP-SERVER](https://github.com/PriyanshuJajal/GO-HTTP-SERVER)
cd micro-http-server
go run main.go
```

---

## 📡 API Endpoints

The router strictly enforces HTTP methods and URIs.

### 1. Check Server Status
*   **Method:** `GET`
*   **Path:** `/api/status`
*   **Response:**
    ```json
    {
      "status": "Server is running smoothly!"
    }
    ```

### 2. Submit Data
*   **Method:** `POST`
*   **Path:** `/api/data`
*   **Headers:** `Content-Length` is strictly required for parsing.
*   **Body:** Any valid JSON payload (e.g., `{"sensor_id": 104, "status": "active"}`)
*   **Response:** The server dynamically echoes back the exact payload you sent.
    ```json
    {
      "message": "Data received!",
      "yourData": { 
        "sensor_id": 104, 
        "status": "active" 
      }
    }
    ```

---

## 📊 Performance & Benchmarks

The server was load-tested using [`hey`](https://github.com/rakyll/hey) with **100 concurrent clients** sending **10,000 total requests** on a standard local environment.

**Command:** `hey -n 10000 -c 100 http://localhost:8080/api/status`

### Benchmark Results

| Metric | Result | Notes |
| :--- | :--- | :--- |
| **Total Requests** | `10,000` | 100% Success Rate |
| **Concurrency Level** | `100` | Handled safely via channel queues |
| **Throughput (RPS)** | `~5923 req/sec` | Sustained throughput with zero GC pressure |
| **Error Rate** | `0.00%` | Zero dropped TCP connections |

### Latency Distribution
Thanks to the zero-allocation routing optimization, the server demonstrates highly consistent, flat tail-latencies under heavy load.

*   **Average:** `16.8 ms`
*   **p50 (Median):** `15.1 ms`
*   **p90:** `21.3 ms`
*   **p99:** `49.3 ms`

---

<p align="center"><i>Built with Go, for developers curious about low‑level backend systems ✨</i></p>