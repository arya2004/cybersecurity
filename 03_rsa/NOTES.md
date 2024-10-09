# RSA Algorithm - Theory and Example in Depth

This Go program implements a basic version of the RSA encryption and decryption algorithm, which is a widely used public-key cryptosystem for secure data transmission. RSA is based on the mathematical properties of prime numbers and modular arithmetic, providing a way for two parties to securely communicate without having to share a secret key.

In this markdown, we will go through the core components of the RSA algorithm, explain the theoretical underpinnings, and provide an in-depth example using the provided code.

## 1. **Key Concepts of RSA**

### 1.1 Prime Numbers and Modulus \(n\)
RSA relies on selecting two large prime numbers \(p\) and \(q\). These prime numbers are used to calculate a modulus \(n\), which is a critical part of both the public and private keys.

\[
n = p \times q
\]
In the example:
- \(p = 3\)
- \(q = 7\)
- \(n = 3 \times 7 = 21\)

### 1.2 Euler's Totient Function \(\phi(n)\)
Euler's Totient function \(\phi(n)\) is used to calculate how many numbers less than \(n\) are coprime with \(n\) (i.e., numbers that share no common factors with \(n\) other than 1). For RSA, \(\phi(n)\) is calculated as:

\[
\phi(n) = (p - 1) \times (q - 1)
\]
In the example:
- \(p - 1 = 3 - 1 = 2\)
- \(q - 1 = 7 - 1 = 6\)
- \(\phi(n) = 2 \times 6 = 12\)

### 1.3 Public and Private Keys
- **Public Key**: The public key consists of the modulus \(n\) and the exponent \(e\), where \(e\) is chosen such that it is coprime with \(\phi(n)\) (i.e., gcd(\(e\), \(\phi(n)\)) = 1).
- **Private Key**: The private key exponent \(d\) is calculated using the following formula:

\[
d = \frac{1 + (k \times \phi(n))}{e}
\]
Here, \(k\) is an integer multiplier to ensure that \(d\) is an integer.

In the example:
- We start with \(e = 2\).
- The gcd of \(e = 2\) and \(\phi(n) = 12\) is not 1, so \(e\) is incremented until we find a value that is coprime with \(\phi(n)\).
- The chosen \(e = 5\) because \(gcd(5, 12) = 1\).
- The private key exponent \(d\) is calculated as:

\[
d = \frac{1 + (2 \times 12)}{5} = \frac{25}{5} = 5
\]

Thus, the public key is \((n = 21, e = 5)\) and the private key is \((n = 21, d = 5)\).

### 1.4 Encryption and Decryption
RSA encryption and decryption are based on modular exponentiation.

- **Encryption**: The message \(msg\) is encrypted using the public key \((n, e)\) with the formula:

\[
c = (msg^e) \mod n
\]
- **Decryption**: The ciphertext \(c\) is decrypted using the private key \((n, d)\) with the formula:

\[
m = (c^d) \mod n
\]

## 2. **RSA Example with the Code**

### 2.1 Given Variables
In this example, the following variables are used:
- \(p = 3\), \(q = 7\) (prime numbers)
- \(n = p \times q = 3 \times 7 = 21\)
- \(\phi(n) = (p - 1) \times (q - 1) = (3 - 1) \times (7 - 1) = 12\)
- Initial public exponent \(e = 2\)
- Multiplier \(k = 2\)
- Message \(msg = 12\)

### 2.2 Step-by-Step Explanation

#### Step 1: Calculate \(n\)
The modulus \(n\) is calculated as:

\[
n = p \times q = 3 \times 7 = 21
\]

#### Step 2: Calculate Euler's Totient Function \(\phi(n)\)
The totient function \(\phi(n)\) is calculated as:

\[
\phi(n) = (p - 1) \times (q - 1) = (3 - 1) \times (7 - 1) = 12
\]

#### Step 3: Choose the Public Exponent \(e\)
We need to choose \(e\) such that \(1 < e < \phi(n)\) and \(gcd(e, \phi(n)) = 1\). Initially, \(e = 2\), but \(gcd(2, 12) = 2\), so we increment \(e\) until we find a coprime value.

Finally, \(e = 5\) is chosen since \(gcd(5, 12) = 1\).

#### Step 4: Calculate the Private Key Exponent \(d\)
The private key exponent \(d\) is calculated using the formula:

\[
d = \frac{1 + (k \times \phi(n))}{e} = \frac{1 + (2 \times 12)}{5} = \frac{25}{5} = 5
\]

Thus, the private key exponent \(d = 5\).

#### Step 5: Encrypt the Message
The encryption process uses the formula:

\[
c = (msg^e) \mod n = (12^5) \mod 21
\]

To calculate this, we first compute \(12^5 = 248832\), and then:

\[
c = 248832 \mod 21 = 12
\]

The encrypted message is \(c = 12\).

#### Step 6: Decrypt the Ciphertext
The decryption process uses the formula:

\[
m = (c^d) \mod n = (12^5) \mod 21
\]

Similar to the encryption, we compute \(12^5 = 248832\) and:

\[
m = 248832 \mod 21 = 12
\]

The decrypted message is \(m = 12\), which matches the original message.


## 4. **Conclusion**
The RSA algorithm provides secure encryption through public and private key pairs. The public key is used to encrypt data, while the private key is used to decrypt it. This example demonstrated how RSA works, from selecting prime numbers to encrypting and decrypting a message using modular arithmetic. While this implementation uses small values, real-world RSA relies on much larger primes to ensure security.