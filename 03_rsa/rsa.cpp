#include <iostream>
#include <cmath>
using namespace std;

// Function to calculate GCD (Greatest Common Divisor)
long long gcd(long long a, long long h) {
    while (h != 0) {
        long long temp = h;
        h = a % h;
        a = temp;
    }
    return a;
}

int main() {
    // Prime numbers
    long long p = 3;
    long long q = 7;

    // Public key exponent and multiplier
    long long e = 2;
    long long k = 2;

    // Message to encrypt
    long long msg = 12;

    // Calculate n = p * q
    long long n = p * q;
    cout << "Calculated n (p * q) = " << n << endl;

    // Calculate phi = (p - 1) * (q - 1)
    long long phi = (p - 1) * (q - 1);
    cout << "Calculated phi ((p - 1) * (q - 1)) = " << phi << endl;

    // Find e such that gcd(e, phi) == 1
    if (gcd(e, phi) != 1) {
        cout << "e = " << e << " is not coprime with phi. Incrementing e..." << endl;
    }

    for (; e < phi; e++) {
        if (gcd(e, phi) == 1) {
            cout << "Chosen e = " << e << " as it is coprime with phi" << endl;
            break;
        }
    }

    // Calculate d = (1 + (k * phi)) / e
    long long d = (1 + (k * phi)) / e;
    cout << "Calculated d (private key) = " << d << endl;

    // Original message
    cout << "Message data = " << msg << endl;

    // Encrypt message: c = (msg ^ e) % n
    double pow1 = pow((double)msg, (double)e);
    long long c = (long long)fmod(pow1, (double)n);
    cout << "Encrypted data = " << c << endl;

    // Decrypt message: m = (c ^ d) % n
    double pow2 = pow((double)c, (double)d);
    long long m = (long long)fmod(pow2, (double)n);
    cout << "Original Message Sent = " << m << endl;

    return 0;
}
