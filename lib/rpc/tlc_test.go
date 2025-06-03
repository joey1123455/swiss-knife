package swissknife

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeebo/assert"
)

func newTestLogger(bufOut, bufErr *bytes.Buffer, level LogLevel) *DefaultLogger {
	return &DefaultLogger{
		infoLogger:  log.New(bufOut, "[TLS-RPC] ", log.LstdFlags),
		errorLogger: log.New(bufErr, "[TLS-RPC] ", log.LstdFlags),
		level:       level,
	}
}

func TestLogLevelFiltering(t *testing.T) {
	var out, err bytes.Buffer
	logger := newTestLogger(&out, &err, LogLevelInfo)

	logger.Debug("this is debug")
	if strings.Contains(out.String(), "DEBUG") {
		t.Error("expected no DEBUG logs at INFO level")
	}

	logger.Info("this is info")
	if !strings.Contains(out.String(), "INFO") {
		t.Error("expected INFO log to be printed")
	}
}

func TestWarnAndErrorLogging(t *testing.T) {
	var out, err bytes.Buffer
	logger := newTestLogger(&out, &err, LogLevelDebug)

	logger.Warn("warn message")
	if !strings.Contains(err.String(), "warn message") {
		t.Error("WARN message not found in stderr")
	}

	if !strings.Contains(err.String(), "WARN:") {
		t.Error("WARN message not found in stderr")
	}

	logger.Error("error message")
	if !strings.Contains(err.String(), "error message") {
		t.Error("ERROR message not found in stderr")
	}

	if !strings.Contains(err.String(), "ERROR:") {
		t.Error("WARN message not found in stderr")
	}
}

func TestFormattedLogging(t *testing.T) {
	var out, err bytes.Buffer
	logger := newTestLogger(&out, &err, LogLevelDebug)

	logger.Infof("info %s", "formatted")
	if !strings.Contains(out.String(), "info formatted") {
		t.Error("Formatted INFO log not printed correctly")
	}
	if !strings.Contains(out.String(), "INFO:") {
		t.Error("Formatted INFO log not printed correctly")
	}

	logger.Warnf("warn %s", "formatted")
	if !strings.Contains(err.String(), "warn formatted") {
		t.Error("Formatted WARN log not printed correctly")
	}

	if !strings.Contains(err.String(), "WARN:") {
		t.Error("Formatted INFO log not printed correctly")
	}
}

func TestSetGetLevel(t *testing.T) {
	logger := NewDefaultLogger()
	logger.SetLevel(LogLevelDebug)

	if logger.GetLevel() != LogLevelDebug {
		t.Errorf("expected level to be DEBUG, got %v", logger.GetLevel())
	}
}

func TestFatalLogging(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		var out, err bytes.Buffer
		logger := newTestLogger(&out, &err, LogLevelDebug)
		logger.Fatalf("fatal error %s", "test")
		t.Fatal("os.Exit(1) should have exited the process")
		return
	}

	cmd := os.Args[0]
	args := []string{"-test.run=TestFatalLogging"}
	env := append(os.Environ(), "TEST_FATAL=1")

	proc := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Env:   env,
	}

	p, err := os.StartProcess(cmd, args, proc)
	if err != nil {
		t.Fatal(err)
	}
	ps, err := p.Wait()
	if err != nil {
		t.Fatal(err)
	}

	if ps.ExitCode() != 1 {
		t.Errorf("expected exit code 1 from Fatalf, got %d", ps.ExitCode())
	}
}

type TestService struct{}

type Args struct {
	A, B int
}

type Reply struct {
	Sum int
}

func (s *TestService) Add(args *[]byte, reply *[]byte) error {
	var argsStruct Args
	json.Unmarshal(*args, &argsStruct)
	var replyStruct Reply
	replyStruct.Sum = argsStruct.A + argsStruct.B
	replyData, err := json.Marshal(replyStruct)
	if err != nil {
		return err
	}

	*reply = replyData
	return nil
}

func getCertPaths(t *testing.T) (certPath, keyPath string) {
	base := filepath.Join("testdata", "tls")
	cert := filepath.Join(base, "server.crt")
	key := filepath.Join(base, "server.key")
	if _, err := os.Stat(cert); os.IsNotExist(err) {
		t.Skip("TLS cert not found")
	}
	if _, err := os.Stat(key); os.IsNotExist(err) {
		t.Skip("TLS key not found")
	}
	return cert, key
}

func TestTlsRpcClientServer(t *testing.T) {
	certPath, keyPath := getCertPaths(t)

	port := "7000"

	server, err := NewITlsRpcServer(certPath, keyPath, port)
	require.NoError(t, err)
	defer server.CloseServer()

	err = server.RegisterMethod("TestService", new(TestService))
	require.NoError(t, err)

	go server.Serve()
	time.Sleep(500 * time.Millisecond) // Give time to start

	client, err := NewITlsRpcClient(certPath, "localhost:"+port, "test-client")
	require.NoError(t, err)
	defer client.CloseClient()

	args := Args{A: 5, B: 3}
	argsData, err := json.Marshal(args)
	require.NoError(t, err)

	replyAny, err := client.ConnectToRpcServerTls("TestService.Add", argsData)
	require.NoError(t, err)

	var reply Reply
	err = json.Unmarshal(replyAny, &reply)
	require.NoError(t, err)
	assert.Equal(t, 8, reply.Sum)
}

func TestServerRejectsInvalidCert(t *testing.T) {
	certPath, keyPath := getCertPaths(t)
	port := "7001"

	server, err := NewITlsRpcServer(certPath, keyPath, port)
	require.NoError(t, err)
	defer server.CloseServer()

	go server.Serve()
	time.Sleep(500 * time.Millisecond)

	// Use an empty cert pool (invalid)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // Force connection even with invalid cert
	}

	_, err = tls.Dial("tcp", "localhost:"+port, tlsConfig)
	require.Error(t, err, "Expected connection to fail due to Incorrect certificate")
}
