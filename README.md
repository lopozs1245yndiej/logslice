# logslice

Stream and filter structured JSON logs from multiple sources with a unified query syntax.

---

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

---

## Usage

Pipe logs directly or point logslice at a file, socket, or remote source and filter using the query syntax:

```bash
# Filter logs by level and service
logslice --source ./app.log 'level == "error" && service == "api"'

# Stream from multiple sources simultaneously
logslice --source ./app.log --source tcp://localhost:5170 'status_code >= 500'

# Pretty-print matched log entries
logslice --pretty --source ./app.log 'user_id == "abc123"'
```

### Query Syntax

| Expression | Description |
|---|---|
| `field == "value"` | Exact string match |
| `field >= 500` | Numeric comparison |
| `&&` / `\|\|` | Logical AND / OR |
| `field ~= "regex"` | Regular expression match |

### Supported Sources

- Local files
- stdin (pipe-friendly)
- TCP / UDP sockets
- HTTP log endpoints

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss significant changes.

---

## License

MIT © 2024 yourusername