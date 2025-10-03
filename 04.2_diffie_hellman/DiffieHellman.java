//To Run this code use command-line and cd into this folder and Run: javac DiffieHellman.java && java DiffieHellman

import java.math.BigInteger;

public class DiffieHellman {

    // Power function to return value of a^b mod P
    public static long power(long a, long b, long P) {
        BigInteger A = BigInteger.valueOf(a);
        BigInteger B = BigInteger.valueOf(b);
        BigInteger Pmod = BigInteger.valueOf(P);

        BigInteger result = A.modPow(B, Pmod);
        return result.longValue();
    }

    // Encrypt/Decrypt function using XOR
    public static String encryptDecrypt(String message, long key) {
        StringBuilder result = new StringBuilder();
        for (char c : message.toCharArray()) {
            result.append((char) (c ^ key));
        }
        return result.toString();
    }

    public static void main(String[] args) {
        long P, G, x, a, y, b, ka, kb;

        // Both the persons will agree upon the public keys G and P
        P = 23; // A prime number
        System.out.println("The value of P: " + P);

        G = 9; // A primitive root for P
        System.out.println("The value of G: " + G);

        // Alice will choose the private key a
        a = 4; // Alice's private key
        System.out.println("The private key a for Alice: " + a);

        x = power(G, a, P); // Alice's public key
        System.out.println("The public key x for Alice: " + x);

        // Bob will choose the private key b
        b = 3; // Bob's private key
        System.out.println("The private key b for Bob: " + b);

        y = power(G, b, P); // Bob's public key
        System.out.println("The public key y for Bob: " + y);

        // Generating the secret key after exchanging public keys
        ka = power(y, a, P); // Secret key for Alice
        kb = power(x, b, P); // Secret key for Bob
        System.out.println("Secret key for Alice is: " + ka);
        System.out.println("Secret key for Bob is:   " + kb);

        // Alice encrypts a message
        String message = "Hello Bob!";
        String encryptedMessage = encryptDecrypt(message, ka);
        System.out.println("Encrypted Message: " + encryptedMessage);

        // Bob decrypts the message
        String decryptedMessage = encryptDecrypt(encryptedMessage, kb);
        System.out.println("Decrypted Message: " + decryptedMessage);
    }
}

