# Swiss Knife Library Documentation

This utility library provides things i keep on rewriting.

## API Implementation

### 🔐 TLS RPC

#### Client

**Function Signatures:**

```go
func NewITlsRpcClient(certPath, address, name string) (ITlsRpcClient, error)

type ITlsRpcClient interface {
    ConnectToRpcServerTls(serviceMethod string, args []byte) ([]byte, error)
    CloseClient()
    SetLogger(logger Logger)
}
```

**Descriptions:**

* **`NewITlsRpcClient(certPath, address, name)`**
  Creates and connects a new TLS RPC client.

* **`ConnectToRpcServerTls(serviceMethod, args)`**
  Calls an RPC method on the connected server.

* **`CloseClient()`**
  Closes the client connection.

* **`SetLogger(logger)`**
  Sets a custom logger for the client.

---

#### Server

**Function Signatures:**

```go
func NewITlsRpcServer(certPath, keyPath, port string) (ITlsRpcServer, error)

type ITlsRpcServer interface {
    RegisterMethod(serviceName string, service any) error
    Serve()
    CloseServer()
    SetLogger(logger Logger)
}
```

**Descriptions:**

* **`NewITlsRpcServer(certPath, keyPath, port)`**
  Creates a new TLS RPC server.

* **`RegisterMethod(serviceName, service)`**
  Registers a service with the server.

* **`Serve()`**
  Starts the server and listens for connections.

* **`CloseServer()`**
  Gracefully shuts down the server.

* **`SetLogger(logger)`**
  Attaches a custom logger.

---

### 📓 Logger

#### LogLevel

```go
type LogLevel int

const (
    LogLevelDebug LogLevel = iota
    LogLevelInfo
    LogLevelWarn
    LogLevelError
    LogLevelFatal
)
```

#### DefaultLogger

```go
type DefaultLogger struct {
    infoLogger  *log.Logger
    errorLogger *log.Logger
    level       LogLevel
}

func NewDefaultLogger() *DefaultLogger

func (dl *DefaultLogger) SetLevel(level LogLevel)
func (dl *DefaultLogger) GetLevel() LogLevel
func (dl *DefaultLogger) shouldLog(level LogLevel) bool

func (dl *DefaultLogger) Debug(v ...interface{})
func (dl *DefaultLogger) Debugf(format string, v ...interface{})
func (dl *DefaultLogger) Info(v ...interface{})
func (dl *DefaultLogger) Warn(v ...interface{})
func (dl *DefaultLogger) Error(v ...interface{})
func (dl *DefaultLogger) Fatal(v ...interface{})
```

**Example:**

```go
logger := NewDefaultLogger()
logger.SetLevel(LogLevelDebug)
logger.Debug("Debugging...")
```

---

### 🔒 AES-GCM Encryption (via `swissknife`)

```go
func EncryptAESGCM(plaintext, key []byte) (ciphertext, nonce []byte, err error)
func DecryptAESGCM(ciphertext, nonce, key []byte) ([]byte, error)
```

**Example:**

```go
key := make([]byte, 32)
_, _ = io.ReadFull(rand.Reader, key)

plaintext := []byte("Hello, World!")
ciphertext, nonce, _ := swissknife.EncryptAESGCM(plaintext, key)
decrypted, _ := swissknife.DecryptAESGCM(ciphertext, nonce, key)

fmt.Println(string(decrypted)) // Hello, World!
```

**Notes:**

* Save nonce alongside ciphertext
* Use secure RNG for key and nonce
* AES-GCM provides encryption and integrity


## Common Errors

* `failed to load cert pool`
* `failed to connect to server`
* `failed to call RPC method`
* `failed to register RPC service`
* `failed to listen on port`
