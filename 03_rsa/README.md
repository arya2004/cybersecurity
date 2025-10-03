# RSA Algorithm Implementations

This folder contains implementations of the RSA encryption algorithm in multiple programming languages.

## Available Implementations

### 1. Go Implementation (`main.go`)
The original Go implementation demonstrating basic RSA encryption and decryption.

**Run with:**
```bash
go run main.go
```

### 2. Python Implementation (`rsa_implementation.py`)
A comprehensive Python implementation with object-oriented design and detailed documentation.

**Run with:**
```bash
python rsa_implementation.py
```

**Features:**
- Clean, modular class-based design
- Detailed documentation and comments
- Step-by-step demonstration output
- Easy-to-understand method separation
- Input validation and error handling

## RSA Algorithm Overview

The RSA algorithm consists of the following key steps:

1. **Choose two prime numbers** `p` and `q`
2. **Calculate modulus** `n = p × q`
3. **Calculate Euler's totient** `φ(n) = (p-1) × (q-1)`
4. **Choose public exponent** `e` such that `gcd(e, φ(n)) = 1`
5. **Calculate private exponent** `d` such that `d × e ≡ 1 (mod φ(n))`
6. **Encrypt:** `ciphertext = message^e mod n`
7. **Decrypt:** `message = ciphertext^d mod n`

## Example Output

Both implementations use the same test values:
- Prime numbers: p=3, q=7
- Message: 12
- Expected encrypted value: 12 (in this small example)

## Security Note

⚠️ **These implementations use small prime numbers for educational purposes only.** 

Real-world RSA implementations should use:
- Prime numbers with at least 1024 bits (preferably 2048 or 4096 bits)
- Proper padding schemes (OAEP)
- Secure random number generation
- Key length of at least 2048 bits for security

## Contributing

Feel free to contribute implementations in other languages such as:
- Java
- JavaScript 
- C++
- Rust
- And more!

Make sure your implementation follows the same algorithmic structure for consistency.