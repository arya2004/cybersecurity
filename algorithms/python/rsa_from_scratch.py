import random


def is_prime(num):
   
    if num < 2:
        return False
    if num in (2, 3):
        return True
    if num % 2 == 0 or num % 3 == 0:
        return False
    i = 5
    while i * i <= num:
        if num % i == 0 or num % (i + 2) == 0:
            return False
        i += 6
    return True

def gcd(a, b):
  
    while b:
        a, b = b, a % b
    return a

def modular_inverse(e, phi):
   
    r0, r1 = phi, e
    x0, x1 = 0, 1

    while r1 != 0:
        q = r0 // r1
       
        r0, r1 = r1, r0 - q * r1

        x0, x1 = x1, x0 - q * x1

    if r0 != 1:
        raise Exception('Modular inverse does not exist (e and phi are not coprime)')

    if x0 < 0:
        x0 += phi

    return x0

def power(a, b, m):

    res = 1
    a %= m
    while b > 0:
        if b & 1:
            res = (res * a) % m
        
        b >>= 1  
        a = (a * a) % m
    return res


def generate_keypair(p, q):

    if not (is_prime(p) and is_prime(q)):
        raise ValueError("Both numbers must be prime.")
    if p == q:
        raise ValueError("p and q cannot be equal.")

    n = p * q

    phi = (p - 1) * (q - 1)

    e = 65537
    
    d = modular_inverse(e, phi)

    return ((e, n), (d, n))


def encrypt(pk, plaintext):
 
    e, n = pk
    cipher = [power(ord(char), e, n) for char in plaintext]
    return cipher

def decrypt(pk, ciphertext):

    d, n = pk
    plain = [chr(power(char, d, n)) for char in ciphertext]
    return ''.join(plain)

p = 101
q = 103

try:
    public, private = generate_keypair(p, q)
    e, n_pub = public
    d, n_priv = private 

    print("--- RSA Key Generation & Encryption Demo ---")
    print(f"Primes (p, q): ({p}, {q})")
    print(f"Modulus (n): {n_pub}")
    print(f"Phi (Totient): {(p - 1) * (q - 1)}")
    print(f"Public Key (e, n): ({e}, {n_pub})")
    print(f"Private Key (d, n): ({d}, {n_priv})")
    print("-" * 45)

    message = "Hello, RSA!"

    encrypted_msg = encrypt(public, message)
    print(f"Original Message: '{message}'")
    print(f"Encrypted (list of numbers): {encrypted_msg}")

    decrypted_msg = decrypt(private, encrypted_msg)
    print(f"Decrypted Message: '{decrypted_msg}'")
    print("-" * 45)

except Exception as err:
    print(f"An error occurred: {err}")