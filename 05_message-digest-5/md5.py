import math
import struct

# Shift amounts
S = [
    7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
    5,  9, 14, 20, 5,  9, 14, 20, 5,  9, 14, 20, 5,  9, 14, 20,
    4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
    6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21
]

T = [int(abs(math.sin(i + 1)) * 2**32) & 0xFFFFFFFF for i in range(64)]

# Non-linear functions
def F(x, y, z): return (x & y) | (~x & z)
def G(x, y, z): return (x & z) | (y & ~z)
def H(x, y, z): return x ^ y ^ z
def I(x, y, z): return y ^ (x | ~z)


def rotate_left(x, n):
    x &= 0xFFFFFFFF
    return ((x << n) | (x >> (32 - n))) & 0xFFFFFFFF


def encode(input_bytes):
    length = (len(input_bytes) * 8) & 0xFFFFFFFFFFFFFFFF
    input_bytes += b'\x80'
    while (len(input_bytes) % 64) != 56:
        input_bytes += b'\x00'
    input_bytes += struct.pack('<Q', length)
    return input_bytes

def md5_transform(input_bytes):
    x = encode(input_bytes)
    a, b, c, d = 0x67452301, 0xEFCDAB89, 0x98BADCFE, 0x10325476

    for i in range(0, len(x), 64):
        chunk = x[i:i+64]
        M = list(struct.unpack('<16I', chunk))
        A, B, C, D = a, b, c, d

        for j in range(16):
            k = j
            s = S[j]
            temp = (A + F(B, C, D) + M[k] + T[j]) & 0xFFFFFFFF
            A, D, C, B = D, C, B, (B + rotate_left(temp, s)) & 0xFFFFFFFF

       
        for j in range(16, 32):
            k = (5*j + 1) % 16
            s = S[j]
            temp = (A + G(B, C, D) + M[k] + T[j]) & 0xFFFFFFFF
            A, D, C, B = D, C, B, (B + rotate_left(temp, s)) & 0xFFFFFFFF

        for j in range(32, 48):
            k = (3*j + 5) % 16
            s = S[j]
            temp = (A + H(B, C, D) + M[k] + T[j]) & 0xFFFFFFFF
            A, D, C, B = D, C, B, (B + rotate_left(temp, s)) & 0xFFFFFFFF

        for j in range(48, 64):
            k = (7*j) % 16
            s = S[j]
            temp = (A + I(B, C, D) + M[k] + T[j]) & 0xFFFFFFFF
            A, D, C, B = D, C, B, (B + rotate_left(temp, s)) & 0xFFFFFFFF

        a = (a + A) & 0xFFFFFFFF
        b = (b + B) & 0xFFFFFFFF
        c = (c + C) & 0xFFFFFFFF
        d = (d + D) & 0xFFFFFFFF

    return struct.pack('<4I', a, b, c, d)

# Calculatate MD5 hash 
def calculate_md5(input_str):
    md5_bytes = md5_transform(input_str.encode())
    return ''.join(f'{b:02x}' for b in md5_bytes)


input_text = input("Enter text to hash with MD5: ")  

md5_hash = calculate_md5(input_text)
print(f"MD5 Hash for input: '{input_text}' is: {md5_hash}")
