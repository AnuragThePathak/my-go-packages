# my-go-packages

A collection of reusable Go utilities.

**Requires Go 1.25+**

```bash
go get github.com/AnuragThePathak/my-go-packages@latest
```

---

## Packages

### `srv` — HTTP Server

Production-ready HTTP server with TLS support, graceful shutdown, and pluggable cleanup hooks.

#### Production usage

```go
import "github.com/AnuragThePathak/my-go-packages/srv"

type myDB struct{ ... }
func (db *myDB) Shutdown(ctx context.Context) error { return db.Close() }

s := srv.NewServer(router, srv.ServerConfig{
    Port:        8080,
    TLSEnabled:  true,
    TLSCertPath: "/etc/ssl/cert.pem",
    TLSKeyPath:  "/etc/ssl/key.pem",
})

s.StartWithGracefulShutdown(ctx, 10*time.Second, &myDB{})
// Blocks until SIGINT/SIGTERM, runs cleanup handlers, then shuts down.
```

#### Test usage

```go
srv, err := s.Start()
if err != nil {
    t.Fatal(err)
}
defer srv.Shutdown(context.Background())

// Server is guaranteed to be listening at this point.
resp, _ := http.Get("http://localhost:8080/health")
```

#### `ServerConfig`

| Field         | Type     | Description                        |
|---------------|----------|------------------------------------|
| `Port`        | `int`    | Port to bind                       |
| `TLSEnabled`  | `bool`   | Enable TLS                         |
| `TLSCertPath` | `string` | Path to TLS certificate file       |
| `TLSKeyPath`  | `string` | Path to TLS private key file       |

#### `CleanupHandler`

Implement this interface to hook into graceful shutdown:

```go
type CleanupHandler interface {
    Shutdown(ctx context.Context) error
}
```

All handlers run in parallel during shutdown.

---

### `env` — Environment Variables

Type-safe helpers to read environment variables with optional defaults.

```go
import "github.com/AnuragThePathak/my-go-packages/env"

// Required — returns error if not set
port, err := env.GetEnvAsInt("PORT")

// With default
port, err := env.GetEnvAsInt("PORT", 8080)

// String
dsn, err := env.GetEnv("DATABASE_URL")

// Bool
debug, err := env.GetEnvAsBool("DEBUG", false)
```

---

## License

MIT
