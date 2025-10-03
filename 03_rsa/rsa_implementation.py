#!/usr/bin/env python3
"""
RSA Algorithm Implementation in Python

This module implements a basic version of the RSA encryption and decryption algorithm.
RSA is a widely used public-key cryptosystem for secure data transmission.

Author: GitHub Copilot
Date: October 2025
"""

import math


class RSA:
    """RSA encryption/decryption implementation"""
    
    def __init__(self, p=3, q=7, e=2, k=2, msg=12):
        """
        Initialize RSA with given parameters
        
        Args:
            p (int): First prime number
            q (int): Second prime number  
            e (int): Initial public key exponent
            k (int): Multiplier for private key calculation
            msg (int): Message to encrypt/decrypt
        """
        self.p = p
        self.q = q
        self.e = e
        self.k = k
        self.msg = msg
        self.n = 0
        self.phi = 0
        self.d = 0
    
    def gcd(self, a, h):
        """
        Calculate the greatest common divisor (GCD) of two numbers using Euclidean algorithm
        
        Args:
            a (int): First number
            h (int): Second number
            
        Returns:
            int: GCD of a and h
        """
        while h != 0:
            temp = a % h
            a = h
            h = temp
        return a
    
    def calculate_n(self):
        """Calculate n which is the product of p and q"""
        self.n = self.p * self.q
        print(f"Calculated n (p * q) = {self.n}")
        return self.n
    
    def calculate_phi(self):
        """Calculate Euler's Totient function (phi) for n"""
        self.phi = (self.p - 1) * (self.q - 1)
        print(f"Calculated phi ((p - 1) * (q - 1)) = {self.phi}")
        return self.phi
    
    def find_public_exponent(self):
        """
        Find a value for e such that 1 < e < phi and gcd(e, phi) == 1
        """
        if self.gcd(self.e, self.phi) != 1:
            print(f"e = {self.e} is not coprime with phi. Incrementing e.")
        
        while self.e < self.phi:
            if self.gcd(self.e, self.phi) == 1:
                print(f"Chosen e = {self.e} as it is coprime with phi")
                break
            else:
                self.e += 1
        
        return self.e
    
    def calculate_private_key(self):
        """
        Calculate the private key exponent 'd' using the formula: d = (1 + (k * phi)) / e
        """
        self.d = (1 + (self.k * self.phi)) // self.e
        print(f"Calculated d (private key) = {self.d}")
        return self.d
    
    def encrypt(self, message=None):
        """
        Encrypt the message using the formula: c = (msg^e) % n
        
        Args:
            message (int, optional): Message to encrypt. Uses self.msg if None.
            
        Returns:
            int: Encrypted message
        """
        if message is None:
            message = self.msg
            
        print(f"Message data = {message}")
        
        # Encrypt using modular exponentiation: c = (msg^e) % n
        encrypted = pow(message, self.e, self.n)
        print(f"Encrypted data = {encrypted}")
        
        return encrypted
    
    def decrypt(self, ciphertext):
        """
        Decrypt the message using the formula: m = (c^d) % n
        
        Args:
            ciphertext (int): Encrypted message to decrypt
            
        Returns:
            int: Decrypted message
        """
        # Decrypt using modular exponentiation: m = (c^d) % n
        decrypted = pow(ciphertext, self.d, self.n)
        print(f"Original Message Sent = {decrypted}")
        
        return decrypted
    
    def run_rsa_demo(self):
        """Run complete RSA encryption/decryption demonstration"""
        print("=" * 50)
        print("RSA Algorithm Demonstration")
        print("=" * 50)
        
        # Step 1: Calculate n
        self.calculate_n()
        
        # Step 2: Calculate phi
        self.calculate_phi()
        
        # Step 3: Find suitable public exponent e
        self.find_public_exponent()
        
        # Step 4: Calculate private key d
        self.calculate_private_key()
        
        # Step 5: Encrypt the message
        encrypted_msg = self.encrypt()
        
        # Step 6: Decrypt the message
        decrypted_msg = self.decrypt(encrypted_msg)
        
        print("=" * 50)
        print("RSA Demo Complete!")
        print(f"Public Key: (n={self.n}, e={self.e})")
        print(f"Private Key: (n={self.n}, d={self.d})")
        print("=" * 50)
        
        return {
            'public_key': (self.n, self.e),
            'private_key': (self.n, self.d),
            'original_message': self.msg,
            'encrypted_message': encrypted_msg,
            'decrypted_message': decrypted_msg
        }


def main():
    """Main function to run RSA demonstration"""
    # Create RSA instance with the same parameters as Go implementation
    rsa = RSA(p=3, q=7, e=2, k=2, msg=12)
    
    # Run the complete demonstration
    results = rsa.run_rsa_demo()
    
    # Verify that decryption worked correctly
    if results['original_message'] == results['decrypted_message']:
        print("✅ RSA encryption/decryption successful!")
    else:
        print("❌ RSA encryption/decryption failed!")


if __name__ == "__main__":
    main()