## Environment Variables
- `HTTP_PORT`: Port to respond to HTTP queries. Defaults to 8080.
- `UDP_PORT`: Port for UDP logs to receive. Defaults to 10000.
- `UDP_BUFFER`: Size of UDP buffer for incoming packets. Minimum of 1024 (1K), defaults to 65000 (almost 65K).
- `KEEP_LOGS`: Number of log lines to keep. Minimum of 1, defaults to 5,000 logs.

## HTTP API
Any path is supported - current recommendation is just `/` (root). The `q` query parameter allows you to filter your query to specific comma-delimited strings that at least one of must be present.
