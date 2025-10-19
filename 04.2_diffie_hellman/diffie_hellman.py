# To run this code, open a terminal in this folder and run: python3 diffie_hellman.py

def power(a, b, P):
    """
    Calculates (a^b) % P using modular exponentiation.
    Python's built-in pow(a, b, P) is efficient for this.
    """
    return pow(a, b, P)

def encrypt_decrypt(message, key):
    """
    Encrypts or decrypts a message using a simple XOR cipher.
    Applying the same function twice with the same key returns the original message.
    """
    # ''.join() is an efficient way to build the resulting string
    return ''.join(chr(ord(char) ^ key) for char in message)

def main():
    """
    Main function to demonstrate the Diffie-Hellman key exchange.
    """
    # Both parties agree upon the public keys G and P.

    # A prime number P is taken.
    P = 23
    print(f"The value of P: {P}")

    # A primitive root for P, G is taken.
    G = 9
    print(f"The value of G: {G}")
    print("-------------------------")

    # Alice chooses her private key 'a'.
    a = 4
    print(f"The private key 'a' for Alice: {a}")

    # Alice calculates her public key 'x'.
    x = power(G, a, P)
    print(f"The public key 'x' for Alice: {x}")
    print("-------------------------")

    # Bob chooses his private key 'b'.
    b = 3
    print(f"The private key 'b' for Bob: {b}")

    # Bob calculates his public key 'y'.
    y = power(G, b, P)
    print(f"The public key 'y' for Bob: {y}")
    print("-------------------------")

    # After exchanging public keys, they generate the shared secret key.

    # Alice calculates the secret key.
    ka = power(y, a, P)

    # Bob calculates the secret key.
    kb = power(x, b, P)

    print(f"Secret key for Alice is: {ka}")
    print(f"Secret key for Bob is: {kb}")
    print("-------------------------")

    # Now they can use the shared secret key for encryption.
    message = "Hello Bob!"
    print(f"Original Message: {message}")

    # Alice encrypts the message with her key.
    encrypted_message = encrypt_decrypt(message, ka)
    print(f"Encrypted Message: {encrypted_message}")

    # Bob decrypts the message with his key.
    decrypted_message = encrypt_decrypt(encrypted_message, kb)
    print(f"Decrypted Message: {decrypted_message}")

# This ensures the main function is called only when the script is executed directly
if __name__ == '__main__':
    main()