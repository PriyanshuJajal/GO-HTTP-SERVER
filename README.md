# ⚡ Go HTTP Server

A high-performance, multi-threaded HTTP web server built entirely from scratch in Go. 

This project intentionally bypasses the standard `net/http` library to demonstrate core backend engineering principles. It implements layer-7 HTTP parsing, concurrent request routing, and thread-safe data handling directly on top of raw TCP sockets (`net` package).

---

## 🧠 Architecture Overview

This server demonstrates a deep understanding of system design and concurrency:

*   **Custom M:N Worker Pool:** Instead of creating a heavy OS thread per connection, the server multiplexes connections across a predefined pool of lightweight Goroutines, minimizing memory footprint and context-switching overhead.
*   **Buffered Channel Queues:** Incoming TCP connections are accepted by a main dispatcher and pushed into a thread-safe `chan net.Conn` job queue, providing built-in backpressure to prevent system overload.
*   **Layer-7 Protocol Parsing:** Custom string parsing utilizing `bufio` to extract HTTP methods, URIs, and intelligently read `Content-Length` headers for dynamic `POST` body extraction.
*   **O(1) Request Router:** A custom hash map (`map[string]RouteHandler`) that instantly matches HTTP methods and URIs to specific handler functions.
*   **Zero Dependencies:** Built using only the Go standard library (`net`, `bufio`, `sync`).

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

To validate the efficiency of the custom Goroutine worker pool and non-blocking I/O, the server was load-tested using [`hey`](https://github.com/rakyll/hey). 

The test simulated **100 concurrent clients** sending a total of **10,000 requests** to the `GET /api/status` endpoint on a standard local environment.

**Command:** 
```bash
hey -n 10000 -c 100 http://localhost:8080/api/status
```

### Sample Benchmark Results

| Metric | Result | Notes |
| :--- | :--- | :--- |
| **Total Requests** | `10,000` | 100% Success Rate |
| **Concurrency Level** | `100` | Handled via custom channel queue |
| **Throughput (RPS)** | `~6,049 req/sec` | High sustained throughput |
| **Error Rate** | `0.00%` | Zero dropped TCP connections or HTTP 500s |

### Latency Distribution
The server demonstrates highly consistent, low-latency responses under heavy concurrent load, proving the efficiency of Go's non-blocking network poller and `bufio` implementation.

*   **Average:** `16.4 ms`
*   **p50 (Median):** `14.7 ms`
*   **p90:** `20.6 ms`
*   **p99:** `46.1 ms`

---

<p align="center"><i>Built with Go, for developers curious about low‑level backend systems ✨</i></p>