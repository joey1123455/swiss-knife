package swissknife

import (
	"crypto/tls"
	"net"
)

type ITlsRpcServer interface {
	CloseServer()
	RegisterMethod(serviceName string, service any) error
	Serve()
	SetLogger(logger Logger)
}

type ITlsRpcClient interface {
	CloseClient()
	SetLogger(logger Logger)
	ConnectToRpcServerTls(serviceMethod string, args []byte) ([]byte, error)
}

type tlsRpcClient struct {
	conn   *tls.Conn
	logger Logger
}

type tlsRpcServer struct {
	listener net.Listener
	logger   Logger
}

type EncryptedRPCStream struct {
	EncryptedStream []byte
	Nonce           []byte
}

type Logger interface {
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Warn(v ...interface{})
	Warnf(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	SetLevel(level LogLevel)
	GetLevel() LogLevel
}
