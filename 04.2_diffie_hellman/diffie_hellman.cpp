// To compile and run this code, open a terminal in this folder and run:
// g++ diffie_hellman.cpp -o diffie_hellman && ./diffie_hellman

#include <iostream>
#include <string>

// Use long long for wider integer range, suitable for this example's numbers.
using ll = long long;

/**
 * @brief Power function to return value of (base^exp) mod modulus.
 * * @param base The base of the exponentiation.
 * @param exp The exponent.
 * @param modulus The modulus.
 * @return The result of (base^exp) mod modulus.
 */
ll power(ll base, ll exp, ll modulus) {
    ll result = 1;
    base %= modulus;
    while (exp > 0) {
        if (exp % 2 == 1) {
            result = (result * base) % modulus;
        }
        exp >>= 1; // equivalent to exp = exp / 2
        base = (base * base) % modulus;
    }
    return result;
}

/**
 * @brief A simple XOR-based encryption/decryption function.
 * * Applying the same function twice with the same key returns the original message.
 * @param message The string message to be processed.
 * @param key The secret key to use for the XOR operation.
 * @return The processed (encrypted or decrypted) string.
 */
std::string encryptDecrypt(const std::string& message, ll key) {
    std::string result = "";
    for (char c : message) {
        result += c ^ key;
    }
    return result;
}

int main() {
    ll P, G, x, a, y, b, ka, kb;

    // Both parties agree upon the public keys G and P.

    // A prime number P is taken.
    P = 23;
    std::cout << "The value of P: " << P << std::endl;

    // A primitive root for P, G is taken.
    G = 9;
    std::cout << "The value of G: " << G << std::endl;
    std::cout << "-------------------------" << std::endl;

    // Alice chooses her private key 'a'.
    a = 4;
    std::cout << "The private key 'a' for Alice: " << a << std::endl;

    // Alice calculates her public key 'x'.
    x = power(G, a, P);
    std::cout << "The public key 'x' for Alice: " << x << std::endl;
    std::cout << "-------------------------" << std::endl;

    // Bob chooses his private key 'b'.
    b = 3;
    std::cout << "The private key 'b' for Bob: " << b << std::endl;

    // Bob calculates his public key 'y'.
    y = power(G, b, P);
    std::cout << "The public key 'y' for Bob: " << y << std::endl;
    std::cout << "-------------------------" << std::endl;

    // After exchanging public keys, they generate the shared secret key.

    // Alice calculates the secret key.
    ka = power(y, a, P);

    // Bob calculates the secret key.
    kb = power(x, b, P);

    std::cout << "Secret key for Alice is: " << ka << std::endl;
    std::cout << "Secret key for Bob is: " << kb << std::endl;
    std::cout << "-------------------------" << std::endl;

    // Now they can use the shared secret key for encryption.
    std::string message = "Hello Bob!";
    std::cout << "Original Message: " << message << std::endl;

    // Alice encrypts the message with her key.
    std::string encryptedMessage = encryptDecrypt(message, ka);
    std::cout << "Encrypted Message: " << encryptedMessage << std::endl;

    // Bob decrypts the message with his key.
    std::string decryptedMessage = encryptDecrypt(encryptedMessage, kb);
    std::cout << "Decrypted Message: " << decryptedMessage << std::endl;

    return 0;
}