package swissknife

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/rpc"
)

// NewITlsRpcClient creates a new TLS RPC client and connects to the specified
// RPC server, using the specified certificate path to load the TLS
// configuration. The client is named for logging purposes.
//
// The returned error is non-nil if the client fails to connect to the server.
func NewITlsRpcClient(certPath, address, name string) (ITlsRpcClient, error) {
	logger := NewDefaultLogger()

	certPool, err := loadCertPool(certPath)
	if err != nil {
		logger.Errorf("Failed to load cert pool for client %s: %v", name, err)
		return nil, fmt.Errorf("failed to load cert pool: %w", err)
	}

	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}

	logger.Infof("Connecting to TLS RPC server at %s for client %s", address, name)
	logger.Debugf("TLS configuration loaded for client %s", name)

	conn, err := tls.Dial("tcp", address, tlsConfig)
	if err != nil {
		logger.Errorf("Connection failed for client %s to %s: %v", name, address, err)
		return nil, fmt.Errorf("failed to connect to server %s: %w", name, err)
	}

	logger.Infof("Successfully connected to TLS RPC server at %s for client %s", address, name)
	return &tlsRpcClient{
		conn:   conn,
		logger: logger,
	}, nil
}

// NewITlsRpcServer creates a new TLS RPC server that listens on the specified
// port and uses the specified certificate and key files to configure the TLS
// connection. The returned error is non-nil if the server fails to listen on
// the specified port.
//
// The returned ITlsRpcServer object is ready to use for RPC registrations and
// serving.
func NewITlsRpcServer(certPath, keyPath, port string) (ITlsRpcServer, error) {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS certificate and key: %w", err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	listener, err := tls.Listen("tcp", fmt.Sprintf(":%s", port), config)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	return &tlsRpcServer{
		listener: listener,
		logger:   NewDefaultLogger(), // Set default logger
	}, nil
}

// CloseClient closes the TLS RPC client connection. If the connection is
// already closed, this function has no effect.
func (c *tlsRpcClient) CloseClient() {
	if c.conn != nil {
		c.logger.Info("Closing TLS RPC client connection")
		c.conn.Close()
	}
}

// SetLogger sets the logger used by the TLS RPC client. If not set, the
// client will use the default logger.
func (c *tlsRpcClient) SetLogger(logger Logger) {
	c.logger = logger
}

// ConnectToRpcServerTls calls the specified RPC method on the connected TLS
// RPC server using the provided arguments. The returned error is non-nil if
// the RPC call fails.
func (c *tlsRpcClient) ConnectToRpcServerTls(serviceMethod string, args []byte) ([]byte, error) {
	c.logger.Infof("Calling RPC method: %s", serviceMethod)
	c.logger.Debugf("RPC method %s called with args: %+v", serviceMethod, args)

	client := rpc.NewClient(c.conn)
	defer client.Close()

	var reply []byte
	err := client.Call(serviceMethod, &args, &reply)
	if err != nil {
		c.logger.Errorf("RPC call failed for method %s: %v", serviceMethod, err)
		return nil, fmt.Errorf("failed to call RPC method %s: %w", serviceMethod, err)
	}

	c.logger.Infof("RPC call successful for method: %s", serviceMethod)
	c.logger.Debugf("RPC method %s returned: %+v", serviceMethod, reply)
	return reply, nil
}

// CloseServer closes the TLS RPC server listener. If the server is already
// closed, this method has no effect.
func (s *tlsRpcServer) CloseServer() {
	if s.listener != nil {
		s.logger.Info("Closing TLS RPC server")
		s.listener.Close()
	}
}

// RegisterMethod registers a new RPC service with the server. The
// serviceName parameter specifies the service name that clients will use
// to access the service. The service parameter is a pointer to the actual service
// implementation.
//
// The returned error is non-nil if the registration fails.
func (s *tlsRpcServer) RegisterMethod(serviceName string, service any) error {
	err := rpc.RegisterName(serviceName, service)
	if err != nil {
		s.logger.Errorf("Failed to register RPC service %s: %v", serviceName, err)
		return fmt.Errorf("failed to register RPC service %s: %w", serviceName, err)
	}
	s.logger.Infof("Successfully registered RPC service: %s", serviceName)
	return nil
}

// SetLogger assigns a custom logger to the TLS RPC server. If not set, the server
// will use the default logger. This allows for customizable logging behavior.

func (s *tlsRpcServer) SetLogger(logger Logger) {
	s.logger = logger
}

// Serve starts the TLS RPC server and listens for incoming connections. The
// server will log a message when it starts, and will log errors if any occur
// while accepting connections. The server will also log a message when a new
// client connects or disconnects. The server will automatically spawn a new
// goroutine to handle each incoming connection.
func (s *tlsRpcServer) Serve() {
	s.logger.Infof("TLS RPC server started, listening on %s", s.listener.Addr().String())
	s.logger.Debugf("Server configuration: TLS enabled, RPC protocol")

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			// Log the error and continue, unless it's due to listener being closed
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "accept" {
				s.logger.Info("Server stopped accepting connections")
				break
			}
			s.logger.Errorf("Failed to accept connection: %v", err)
			continue
		}

		s.logger.Infof("New client connected from %s", conn.RemoteAddr().String())
		s.logger.Debugf("Client connection details: %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
		go s.handleConnection(conn)
	}
}

// handleConnection manages a single RPC connection for the TLS RPC server.
// It logs the connection details, serves the RPC requests, and ensures
// the connection is closed after use. The function also captures and logs
// any panics that may occur during RPC handling to prevent the server from
// crashing.
func (s *tlsRpcServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.logger.Infof("Connection closed for client %s", conn.RemoteAddr().String())
	}()

	s.logger.Infof("Serving RPC connection for client %s", conn.RemoteAddr().String())
	s.logger.Debugf("Starting RPC handler for client %s", conn.RemoteAddr().String())

	// Handle potential panics in RPC serving
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("RPC handler panic for client %s: %v", conn.RemoteAddr().String(), r)
		}
	}()

	rpc.ServeConn(conn)
	s.logger.Debugf("RPC handler finished for client %s", conn.RemoteAddr().String())
}
