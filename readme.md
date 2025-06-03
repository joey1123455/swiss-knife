# Swiss Knife
Swiss Knife is a collection of utility functions and libraries for Go programming language. It provides a set of reusable and modular components that can be used to simplify and accelerate development of various applications.

## Table of Contents
* [Introduction](#introduction)
* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
* [API Reference](#api-reference)
* [License](#license)
* [Author](#author)
* [Contributing](#contributing)

## Introduction
Swiss Knife is designed to be a versatile and flexible library that can be used in a wide range of applications, from command-line tools to web applications. It provides a set of utility functions and libraries that can be used to perform various tasks, such as encryption, decryption, logging, and more.

## Features
* **Encryption and Decryption**: Swiss Knife provides a set of encryption and decryption functions using AES-GCM algorithm.
* **Logging**: Swiss Knife provides a logging library that can be used to log messages at different levels (debug, info, warn, error, fatal).
* **TLS Support**: Swiss Knife provides support for TLS connections, including certificate generation and verification.
* **RPC Support**: Swiss Knife provides support for RPC connections, including client and server implementations.

## Installation
To install Swiss Knife, you can use the following command:
```bash
go get github.com/joey1123455/swiss-knife
```
## Usage
Swiss Knife is designed to be easy to use. Here are some examples of how to use the library:

### Encryption and Decryption
```go
package main

import (
	"fmt"
	"github.com/joey1123455/swiss-knife/lib/encryptions"
)

func main() {
	key := []byte("examplekey123456")
	plaintext := []byte("Hello, World!")

	ciphertext, nonce, err := swissknife.EncryptAESGCM(plaintext, key)
	if err != nil {
		fmt.Println("Error encrypting:", err)
		return
	}

	fmt.Printf("Ciphertext: %x\n", ciphertext)
	fmt.Printf("Nonce: %x\n", nonce)

	decrypted, err := swissknife.DecryptAESGCM(ciphertext, nonce, key)
	if err != nil {
		fmt.Println("Error decrypting:", err)
		return
	}

	fmt.Printf("Decrypted: %s\n", decrypted)
}
```

### Logging
```go
package main

import (
	"fmt"
	"github.com/joey1123455/swiss-knife/lib/rpc/logger"
)

func main() {
	log := logger.NewDefaultLogger()
	log.Info("This is an info message")
	log.Warn("This is a warning message")
	log.Error("This is an error message")
}
```

### RPC Support
```go
package main

import (
	"fmt"
	"github.com/joey1123455/swiss-knife/lib/rpc/types"
)

func main() {
	rpcServer := types.NewITlsRpcServer("localhost:8080")
	rpcServer.RegisterMethod("TestService", new(TestService))

	go rpcServer.Serve()

	rpcClient := types.NewITlsRpcClient("localhost:8080")
	rpcClient.ConnectToRpcServerTls("TestService.Add", []byte("Hello, World!"))
}
```

## API Reference
Swiss Knife provides a comprehensive API reference that can be found in the [API Reference](https://github.com/joey1123455/swiss-knife/blob/main/API.md) document.

## License
Swiss Knife is licensed under the MIT License. See the [LICENSE](https://github.com/joey1123455/swiss-knife/blob/main/LICENSE) file for details.

## Author
Swiss Knife was developed by [zah_gopher](https://github.com/joey1123455).

## Contributing
Contributions to Swiss Knife are welcome. If you would like to contribute, please fork the repository and submit a pull request.

please run test using the make command
```bash
make tests
```
