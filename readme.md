# Swiss Knife Library
Swiss knife is just a random bunch of things i noticed i kept on re writting in diffrent projects.

## Installation
To use SwissKnife in your project, install it using go get:
```bash
go get github.com/joey1123455/swissknife
```

## Usage 
### Encrypting Data
The `EncryptAESGCM` function encrypts plaintext using AES-GCM. It requires a key and returns the ciphertext, nonce, and an error (if any).

```go
package main

import (
    "fmt"
    swissknife "github.com/joey1123455/swiss-knife/lib/encryptions"
)

func main() {
    key := []byte("examplekey123456") // 16, 24, or 32 bytes for AES
    plaintext := []byte("Hello, World!")

    ciphertext, nonce, err := swissknife.EncryptAESGCM(plaintext, key)
    if err != nil {
        fmt.Println("Error encrypting:", err)
        return
    }

    fmt.Printf("Ciphertext: %x\n", ciphertext)
    fmt.Printf("Nonce: %x\n", nonce)
}
```

### Decrypting Data
The `DecryptAESGCM` function decrypts ciphertext using AES-GCM. It requires the ciphertext, nonce, and key, and returns the plaintext or an error.

```go
package main

import (
    "fmt"
    swissknife "github.com/joey1123455/swiss-knife/lib/encryptions"
)

func main() {
    key := []byte("examplekey123456") // Same key used for encryption
    ciphertext := []byte{...}         // Ciphertext from encryption
    nonce := []byte{...}              // Nonce from encryption

    plaintext, err := swissknife.DecryptAESGCM(ciphertext, nonce, key)
    if err != nil {
        fmt.Println("Error decrypting:", err)
        return
    }

    fmt.Printf("Plaintext: %s\n", plaintext)
}
```

## API Refrence
`EncryptAESGCM`
```go
func EncryptAESGCM(plaintext []byte, key []byte) ([]byte, []byte, error)
```

- Parameters:
    - `plaintext`: The data to encrypt.
    - `key`: The encryption key (16, 24, or 32 bytes for AES).
- Returns:
    - `ciphertext`: The encrypted data.
    - `nonce`: The randomly generated nonce used for encryption.
    - `error`: An error if encryption fails.

`DecryptAESGCM`
```go
func DecryptAESGCM(ciphertext []byte, nonce []byte, key []byte) ([]byte, error)
```
- Parameters:
    - `ciphertext`: The encrypted data.
    - `nonce`: The nonce used during encryption.
    - `key`: The decryption key (same as the encryption key).
- Returns:
    - `plaintext`: The decrypted data.
    - `error`: An error if decryption fails.


## License
This library is licensed under the MIT License. See the LICENSE file for details.

## Author
Developed by zah_gopher.