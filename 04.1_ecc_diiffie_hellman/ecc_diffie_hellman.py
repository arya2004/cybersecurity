# To run this script, use the command: python3 ecc_diffie_hellman.py

import random

"""
An educational implementation of the Elliptic Curve Diffie-Hellman (ECDH) key exchange protocol in Python.

This script implements the necessary elliptic curve arithmetic from scratch to demonstrate the
underlying principles of ECDH. It uses a small, simple curve for clarity.

WARNING: This implementation is for educational purposes only. It is NOT secure for real-world use
as it lacks many security features (e.g., proper random number generation, protection against
side-channel attacks, use of standardized secure curves).
"""

class EllipticCurve:
    """
    Represents an elliptic curve defined by the equation y^2 = x^3 + ax + b (mod p).
    """
    def __init__(self, a, b, p, base_point):
        self.a = a
        self.b = b
        self.p = p
        self.base_point = base_point  # The generator point G

        # The point at infinity, which serves as the identity element
        self.infinity = None

    def _inverse_mod_p(self, n):
        """
        Calculates the modular multiplicative inverse of n modulo p.
        Uses Python's built-in pow(n, -1, p) which is efficient.
        """
        # Ensure the number is not negative before calculating the inverse.
        n = n % self.p
        if n < 0:
            n += self.p
        return pow(n, -1, self.p)

    def is_on_curve(self, point):
        """
        Checks if a given point (x, y) lies on the curve.
        """
        if point == self.infinity:
            return True
        x, y = point
        # This is the core check: y^2 === x^3 + ax + b (mod p)
        return (y * y - (x * x * x + self.a * x + self.b)) % self.p == 0

    def add_points(self, p1, p2):
        """
        Adds two points on the elliptic curve (p1 + p2).
        """
        # If one point is the point at infinity, the sum is the other point.
        if p1 == self.infinity:
            return p2
        if p2 == self.infinity:
            return p1

        x1, y1 = p1
        x2, y2 = p2

        # If the points are inverses of each other, their sum is the point at infinity.
        if x1 == x2 and y1 != y2:
            return self.infinity

        # The calculation of the slope 'm' depends on whether we are adding two distinct points
        # or doubling a single point.
        if x1 == x2:  # Point doubling
            # Calculate the slope of the tangent line: m = (3*x1^2 + a) * (2*y1)^-1 mod p
            numerator = (3 * x1 * x1 + self.a) % self.p
            denominator = (2 * y1) % self.p
            m = (numerator * self._inverse_mod_p(denominator)) % self.p
        else:  # Point addition
            # Calculate the slope of the line between p1 and p2: m = (y2 - y1) * (x2 - x1)^-1 mod p
            numerator = (y2 - y1) % self.p
            denominator = (x2 - x1) % self.p
            m = (numerator * self._inverse_mod_p(denominator)) % self.p

        # Use the slope 'm' to calculate the coordinates of the new point (x3, y3).
        x3 = (m * m - x1 - x2) % self.p
        y3 = (m * (x1 - x3) - y1) % self.p

        return (x3, y3)

    def scalar_multiply(self, n, point):
        """
        Performs scalar multiplication (n * P) using the double-and-add algorithm.
        This is the core of ECC, providing a one-way function. It's easy to compute
        n * P, but hard to find 'n' given P and (n * P).
        """
        result = self.infinity
        current = point

        # Process the integer n in its binary representation.
        while n > 0:
            # If the current bit of 'n' is 1, we add the current point to the result.
            if n % 2 == 1:
                result = self.add_points(result, current)

            # We "double" the point for the next bit of 'n'.
            current = self.add_points(current, current)
            n //= 2

        return result

def xor_encrypt_decrypt(message, key):
    """Simple XOR encryption/decryption using the shared secret key."""
    # Use the first byte of the key for a simple XOR cipher
    key_byte = key.to_bytes((key.bit_length() + 7) // 8, 'big')[0]
    return ''.join(chr(ord(c) ^ key_byte) for c in message)

def find_all_points(curve):
    """
    Finds and prints all the points that lie on the given elliptic curve.
    This function iterates through every possible (x, y) coordinate pair
    and uses the is_on_curve method to check if it satisfies the curve equation.
    """
    points = []
    # Iterate through all possible x and y values in the finite field Z_p.
    for x in range(curve.p):
        for y in range(curve.p):
            point = (x, y)
            if curve.is_on_curve(point):
                points.append(point)
                print(point)
    return points

def main():
    """
    Main function to demonstrate finding points and the ECDH key exchange.
    """
    # --- Part 1: Find all points on a user-defined curve ---
    print("--- Elliptic Curve Point Finder ---")
    try:
        a = int(input("Enter curve parameter 'a': "))
        b = int(input("Enter curve parameter 'b': "))
        p = int(input("Enter curve prime modulus 'p': "))
    except ValueError:
        print("Invalid input. Please enter integers only.")
        return

    # Create a curve object to find the points. The base point is not needed yet.
    curve_for_points = EllipticCurve(a, b, p, None)
    print(f"\nAll points on the curve y^2 = x^3 + {a}x + {b} (mod {p}) are:")
    find_all_points(curve_for_points)

    # --- Part 2: Demonstrate ECDH with the user's curve ---
    print("\n--- Elliptic Curve Diffie-Hellman Key Exchange ---")
    print("To proceed, please choose a base point (G) from the list above.")
    try:
        gx = int(input("Enter the x-coordinate of your chosen base point G: "))
        gy = int(input("Enter the y-coordinate of your chosen base point G: "))
    except ValueError:
        print("Invalid input. Please enter integers only.")
        return

    base_point_g = (gx, gy)

    # Create the final curve object for the key exchange
    curve = EllipticCurve(a, b, p, base_point_g)

    # Verify that the user-chosen point is actually on the curve before proceeding.
    if not curve.is_on_curve(base_point_g):
        print("\n❌ Error: The chosen point is not on the curve. Exiting.")
        return

    print(f"\nUsing Base Point G: {base_point_g}\n")

    # --- Alice's Side ---
    # 1. Alice chooses a secret private key (a random integer, nA)
    private_key_alice = random.randint(1, p - 1)
    print(f"Alice's private key (nA): {private_key_alice}")

    # 2. Alice computes her public key: PA = nA * G
    public_key_alice = curve.scalar_multiply(private_key_alice, curve.base_point)
    print(f"Alice's public key (PA = nA*G): {public_key_alice}\n")

    # --- Bob's Side ---
    # 1. Bob chooses his secret private key (nB)
    private_key_bob = random.randint(1, p - 1)
    print(f"Bob's private key (nB): {private_key_bob}")

    # 2. Bob computes his public key: PB = nB * G
    public_key_bob = curve.scalar_multiply(private_key_bob, curve.base_point)
    print(f"Bob's public key (PB = nB*G): {public_key_bob}\n")

    # --- Key Exchange ---
    print("--- Shared Secret Calculation ---")
    # Alice and Bob now exchange their public keys.

    # 3. Alice computes the shared secret using her private key and Bob's public key: S = nA * PB
    shared_secret_alice = curve.scalar_multiply(private_key_alice, public_key_bob)
    print(f"Alice calculates shared secret (S = nA*PB): {shared_secret_alice}")

    # 4. Bob computes the shared secret using his private key and Alice's public key: S = nB * PA
    shared_secret_bob = curve.scalar_multiply(private_key_bob, public_key_alice)
    print(f"Bob calculates shared secret (S = nB*PA):   {shared_secret_bob}\n")

    # 5. Verification
    if shared_secret_alice == shared_secret_bob:
        print("✅ Success! Both parties have calculated the same shared secret.")
        # The x-coordinate of the shared point is typically used as the symmetric key.
        shared_key = shared_secret_alice[0]
        print(f"The shared symmetric key is the x-coordinate: {shared_key}\n")

        # --- Example of using the shared key ---
        message = "Hello, Bob!"
        print(f"Original message from Alice: '{message}'")
        encrypted = xor_encrypt_decrypt(message, shared_key)
        print(f"Encrypted message sent to Bob: '{encrypted}'")
        decrypted = xor_encrypt_decrypt(encrypted, shared_key)
        print(f"Bob decrypts the message: '{decrypted}'")
    else:
        print("❌ Failure! The shared secrets do not match.")

if __name__ == "__main__":
    main()