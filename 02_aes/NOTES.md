# Simplified AES Algorithm - Theory and Example

This implementation demonstrates a simplified version of the AES (Advanced Encryption Standard) encryption algorithm. AES is a symmetric key cipher, which means the same key is used for both encryption and decryption. It operates on blocks of data and transforms them through a series of mathematical operations to secure the data. In this simplified version, we reduce AES down to smaller blocks (16 bits) and simpler operations to illustrate the basic principles.

## Key Components of AES

### 1. **Substitution Box (S-Box)**
   The S-Box is a fundamental part of AES, used to introduce non-linearity into the encryption process. It replaces each nibble (4-bit unit) of the state with another nibble based on a lookup table. This process, called substitution, provides confusion to the encryption process.

   In our simplified version, we use a 4-bit S-Box with 16 possible values:

   ```go
   var sBox = [16]int{
	   0x9, 0x4, 0xA, 0xB, 0xD, 0x1, 0x8, 0x5, 0x6, 0x2, 0x0, 0x3, 0xC, 0xE, 0xF, 0x7,
   }
   ```

   The inverse S-Box is used for decryption, reversing the substitution:

   ```go
   var sBoxI = [16]int{
	   0xA, 0x5, 0x9, 0xB, 0x1, 0x7, 0x8, 0xF, 0x6, 0x0, 0x2, 0x3, 0xC, 0x4, 0xD, 0xE,
   }
   ```

### 2. **Shift Rows**
   This operation shifts the rows of the state. In the simplified AES, we swap the second nibble and the fourth nibble to shift the "rows." This introduces diffusion in the encryption process.

   ```go
   func shiftRows(state []int) []int {
	   return []int{state[0], state[1], state[3], state[2]}
   }
   ```

### 3. **Mix Columns**
   In this step, the columns of the state are mixed. In the simplified AES, we use a simple multiplication in the Galois field (GF(2^4)) to mix the columns. This further diffuses the data and ensures that each output bit depends on every input bit.

   ```go
   func mixColumns(state []int) []int {
	   return []int{
		   state[0] ^ gfMult(4, state[2]),
		   state[1] ^ gfMult(4, state[3]),
		   state[2] ^ gfMult(4, state[0]),
		   state[3] ^ gfMult(4, state[1]),
	   }
   }
   ```

   For decryption, we use the inverse of the mix columns operation:

   ```go
   func inverseMixColumns(state []int) []int {
	   return []int{
		   gfMult(9, state[0]) ^ gfMult(2, state[2]),
		   gfMult(9, state[1]) ^ gfMult(2, state[3]),
		   gfMult(9, state[2]) ^ gfMult(2, state[0]),
		   gfMult(9, state[3]) ^ gfMult(2, state[1]),
	   }
   }
   ```

### 4. **Add Round Key**
   In the add round key step, the state is XORed with the round key. This step introduces the secret key into the transformation and ensures that encryption is dependent on the key.

   ```go
   func addRoundKey(s1, s2 []int) []int {
	   result := make([]int, len(s1))
	   for i := range s1 {
		   result[i] = s1[i] ^ s2[i]
	   }
	   return result
   }
   ```

### 5. **Key Expansion**
   AES uses a series of round keys derived from the original secret key. The process of generating these round keys is called key expansion. In this simplified version, the key expansion generates three keys: the pre-round key, round 1 key, and round 2 key. These keys are derived by using the `subWord`, `rotWord`, and XOR operations with round constants.

   ```go
   func keyExpansion(key int) ([]int, []int, []int) {
	   Rcon1 := 0x80
	   Rcon2 := 0x30
	   w := make([]int, 6)
	   w[0] = (key & 0xFF00) >> 8
	   w[1] = key & 0x00FF
	   w[2] = w[0] ^ (subWord(rotWord(w[1])) ^ Rcon1)
	   w[3] = w[2] ^ w[1]
	   w[4] = w[2] ^ (subWord(rotWord(w[3])) ^ Rcon2)
	   w[5] = w[4] ^ w[3]
	   return intToState((w[0] << 8) + w[1]), intToState((w[2] << 8) + w[3]), intToState((w[4] << 8) + w[5])
   }
   ```

## Example of Encryption and Decryption

Let's walk through the encryption and decryption process step by step using the following parameters:

- **Key**: `0100101011110101` (16 bits)
- **Plaintext**: `1101011100101000` (16 bits)

### Step 1: Key Expansion

We begin by expanding the key into three round keys: the **pre-round key**, **round 1 key**, and **round 2 key**. This is done through a process of substitution, rotation, and XOR operations.

For this key, the expanded keys are:

- **Pre-round key**: `0x4AF5`
- **Round 1 key**: `0xAE99`
- **Round 2 key**: `0x4CB4`

### Step 2: Initial Add Round Key

The plaintext is XORed with the pre-round key:

```plaintext ⊕ pre-round key = 1101011100101000 ⊕ 0100101011110101 = 1001110111011101```

### Step 3: Substitution (S-Box)

The nibbles of the state are substituted using the S-Box. After substitution, the state becomes:

```1001110111011101 → 1010110111111110```

### Step 4: Shift Rows

The second and fourth nibbles of the state# Simplified AES Algorithm - Theory and Example

This implementation demonstrates a simplified version of the AES (Advanced Encryption Standard) encryption algorithm. AES is a symmetric key cipher, which means the same key is used for both encryption and decryption. It operates on blocks of data and transforms them through a series of mathematical operations to secure the data. In this simplified version, we reduce AES down to smaller blocks (16 bits) and simpler operations to illustrate the basic principles.

## Key Components of AES

### 1. **Substitution Box (S-Box)**
   The S-Box is a fundamental part of AES, used to introduce non-linearity into the encryption process. It replaces each nibble (4-bit unit) of the state with another nibble based on a lookup table. This process, called substitution, provides confusion to the encryption process.

   In our simplified version, we use a 4-bit S-Box with 16 possible values:

   ```go
   var sBox = [16]int{
	   0x9, 0x4, 0xA, 0xB, 0xD, 0x1, 0x8, 0x5, 0x6, 0x2, 0x0, 0x3, 0xC, 0xE, 0xF, 0x7,
   }
   ```

   The inverse S-Box is used for decryption, reversing the substitution:

   ```go
   var sBoxI = [16]int{
	   0xA, 0x5, 0x9, 0xB, 0x1, 0x7, 0x8, 0xF, 0x6, 0x0, 0x2, 0x3, 0xC, 0x4, 0xD, 0xE,
   }
   ```

### 2. **Shift Rows**
   This operation shifts the rows of the state. In the simplified AES, we swap the second nibble and the fourth nibble to shift the "rows." This introduces diffusion in the encryption process.

   ```go
   func shiftRows(state []int) []int {
	   return []int{state[0], state[1], state[3], state[2]}
   }
   ```

### 3. **Mix Columns**
   In this step, the columns of the state are mixed. In the simplified AES, we use a simple multiplication in the Galois field (GF(2^4)) to mix the columns. This further diffuses the data and ensures that each output bit depends on every input bit.

   ```go
   func mixColumns(state []int) []int {
	   return []int{
		   state[0] ^ gfMult(4, state[2]),
		   state[1] ^ gfMult(4, state[3]),
		   state[2] ^ gfMult(4, state[0]),
		   state[3] ^ gfMult(4, state[1]),
	   }
   }
   ```

   For decryption, we use the inverse of the mix columns operation:

   ```go
   func inverseMixColumns(state []int) []int {
	   return []int{
		   gfMult(9, state[0]) ^ gfMult(2, state[2]),
		   gfMult(9, state[1]) ^ gfMult(2, state[3]),
		   gfMult(9, state[2]) ^ gfMult(2, state[0]),
		   gfMult(9, state[3]) ^ gfMult(2, state[1]),
	   }
   }
   ```

### 4. **Add Round Key**
   In the add round key step, the state is XORed with the round key. This step introduces the secret key into the transformation and ensures that encryption is dependent on the key.

   ```go
   func addRoundKey(s1, s2 []int) []int {
	   result := make([]int, len(s1))
	   for i := range s1 {
		   result[i] = s1[i] ^ s2[i]
	   }
	   return result
   }
   ```

### 5. **Key Expansion**
   AES uses a series of round keys derived from the original secret key. The process of generating these round keys is called key expansion. In this simplified version, the key expansion generates three keys: the pre-round key, round 1 key, and round 2 key. These keys are derived by using the `subWord`, `rotWord`, and XOR operations with round constants.

   ```go
   func keyExpansion(key int) ([]int, []int, []int) {
	   Rcon1 := 0x80
	   Rcon2 := 0x30
	   w := make([]int, 6)
	   w[0] = (key & 0xFF00) >> 8
	   w[1] = key & 0x00FF
	   w[2] = w[0] ^ (subWord(rotWord(w[1])) ^ Rcon1)
	   w[3] = w[2] ^ w[1]
	   w[4] = w[2] ^ (subWord(rotWord(w[3])) ^ Rcon2)
	   w[5] = w[4] ^ w[3]
	   return intToState((w[0] << 8) + w[1]), intToState((w[2] << 8) + w[3]), intToState((w[4] << 8) + w[5])
   }
   ```

## Example of Encryption and Decryption

Let's walk through the encryption and decryption process step by step using the following parameters:

- **Key**: `0100101011110101` (16 bits)
- **Plaintext**: `1101011100101000` (16 bits)

### Step 1: Key Expansion

We begin by expanding the key into three round keys: the **pre-round key**, **round 1 key**, and **round 2 key**. This is done through a process of substitution, rotation, and XOR operations.

For this key, the expanded keys are:

- **Pre-round key**: `0x4AF5`
- **Round 1 key**: `0xAE99`
- **Round 2 key**: `0x4CB4`

### Step 2: Initial Add Round Key

The plaintext is XORed with the pre-round key:

```plaintext ⊕ pre-round key = 1101011100101000 ⊕ 0100101011110101 = 1001110111011101```

### Step 3: Substitution (S-Box)

The nibbles of the state are substituted using the S-Box. After substitution, the state becomes:

```1001110111011101 → 1010110111111110```

### Step 4: Shift Rows

The second and fourth nibbles of the state are swapped:

```1010110111111110 → 1010111111111011```

### Step 5: Mix Columns

We apply the mix columns operation, which mixes the columns of the state using Galois field multiplication. The result is:

```1010111111111011 → 1100110101010100```

### Step 6: Add Round Key (Round 1)

We XOR the result with the round 1 key:

```1100110101010100 ⊕ 1010111010011001 = 0110001111001101```

### Step 7: Substitution (S-Box)

We substitute the nibbles again using the S-Box. The result becomes:

```0110001111001101 → 1001011011011111```

### Step 8: Shift Rows

We swap the second and fourth nibbles again:

```1001011011011111 → 1001011111011011```

### Step 9: Add Round Key (Round 2)

Finally, we XOR the result with the round 2 key:

```1001011111011011 ⊕ 0100110010110100 = 1101101101101111```

The ciphertext after encryption is:

```
Ciphertext: 1101101101101111
```

### Decryption Process

Decryption in AES is the reverse of the encryption process. The ciphertext is passed through the inverse transformations, using the same round keys in reverse order. Here's how the decryption proceeds:

1. **Add Round Key (Round 2)**: XOR with the round 2 key.
2. **Inverse Shift Rows**: Swap the nibbles back.
3. **Inverse Substitution (Inverse S-Box)**: Replace the nibbles using the inverse S-Box.
4. **Add Round Key (Round 1)**: XOR with the round 1 key.
5. **Inverse Mix Columns**: Reverse the column mixing using inverse Galois field multiplication.
6. **Inverse Shift Rows**: Swap the nibbles again.
7. **Inverse Substitution (Inverse S-Box)**: Apply the inverse S-Box again.
8. **Add Round Key (Pre-Round Key)**: XOR with the pre-round key to recover the original plaintext.

After following the decryption process, the original plaintext `1101011100101000` is recovered.

## Conclusion

This simplified version of AES illustrates the basic principles of modern symmetric-key block ciphers:

- **Substitution (S-Box)** introduces non-linearity.
- **Shift Rows** and **Mix Columns** introduce diffusion, ensuring that small changes in the plaintext or key result in significant changes to the ciphertext

.
- **Add Round Key** incorporates the secret key into the encryption process.

The simplified AES is a helpful tool for learning how real-world ciphers work, though it lacks the security features of full AES. are swapped:

```1010110111111110 → 1010111111111011```

### Step 5: Mix Columns

We apply the mix columns operation, which mixes the columns of the state using Galois field multiplication. The result is:

```1010111111111011 → 1100110101010100```

### Step 6: Add Round Key (Round 1)

We XOR the result with the round 1 key:

```1100110101010100 ⊕ 1010111010011001 = 0110001111001101```

### Step 7: Substitution (S-Box)

We substitute the nibbles again using the S-Box. The result becomes:

```0110001111001101 → 1001011011011111```

### Step 8: Shift Rows

We swap the second and fourth nibbles again:

```1001011011011111 → 1001011111011011```

### Step 9: Add Round Key (Round 2)

Finally, we XOR the result with the round 2 key:

```1001011111011011 ⊕ 0100110010110100 = 1101101101101111```

The ciphertext after encryption is:

```
Ciphertext: 1101101101101111
```

### Decryption Process

Decryption in AES is the reverse of the encryption process. The ciphertext is passed through the inverse transformations, using the same round keys in reverse order. Here's how the decryption proceeds:

1. **Add Round Key (Round 2)**: XOR with the round 2 key.
2. **Inverse Shift Rows**: Swap the nibbles back.
3. **Inverse Substitution (Inverse S-Box)**: Replace the nibbles using the inverse S-Box.
4. **Add Round Key (Round 1)**: XOR with the round 1 key.
5. **Inverse Mix Columns**: Reverse the column mixing using inverse Galois field multiplication.
6. **Inverse Shift Rows**: Swap the nibbles again.
7. **Inverse Substitution (Inverse S-Box)**: Apply the inverse S-Box again.
8. **Add Round Key (Pre-Round Key)**: XOR with the pre-round key to recover the original plaintext.

After following the decryption process, the original plaintext `1101011100101000` is recovered.

## Conclusion

This simplified version of AES illustrates the basic principles of modern symmetric-key block ciphers:

- **Substitution (S-Box)** introduces non-linearity.
- **Shift Rows** and **Mix Columns** introduce diffusion, ensuring that small changes in the plaintext or key result in significant changes to the ciphertext

.
- **Add Round Key** incorporates the secret key into the encryption process.

The simplified AES is a helpful tool for learning how real-world ciphers work, though it lacks the security features of full AES.