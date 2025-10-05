public class rsa {

    // Function to calculate the greatest common divisor (GCD) of two numbers
    public static long gcd(long a, long h) {
        while (h != 0) {
            long temp = h;
            h = a % h;
            a = temp;
        }
        return a;
    }

    public static void main(String[] args) {
        // Variables for the prime numbers p and q
        long p = 3;
        long q = 7;

        // Variables for the public key exponent 'e', a multiplier 'k', and the message 'msg'
        long e = 2;
        long k = 2;
        long msg = 12;

        // Calculate n which is the product of p and q
        long n = p * q;
        System.out.println("Calculated n (p * q) = " + n);

        // Calculate Euler's Totient function (phi) for n
        long phi = (p - 1) * (q - 1);
        System.out.println("Calculated phi ((p - 1) * (q - 1)) = " + phi);

        // Find a value for e such that 1 < e < phi and gcd(e, phi) == 1
        if (gcd(e, phi) != 1) {
            System.out.println("e = " + e + " is not coprime with phi. Incrementing e.");
        }

        for (; e < phi; ) {
            if (gcd(e, phi) == 1) {
                System.out.println("Chosen e = " + e + " as it is coprime with phi");
                break;
            } else {
                e++;
            }
        }

        // Calculate the private key exponent 'd' using the formula: d = (1 + (k * phi)) / e
        long d = (1 + (k * phi)) / e;
        System.out.println("Calculated d (private key) = " + d);

        // Print the original message data
        System.out.println("Message data = " + msg);

        // Encrypt the message using the formula: c = (msg^e) % n
        double pow1 = Math.pow((double) msg, (double) e);
        long c = (long) (pow1 % (double) n);
        System.out.println("Encrypted data = " + c);

        // Decrypt the message using the formula: m = (c^d) % n
        double pow2 = Math.pow((double) c, (double) d);
        long m = (long) (pow2 % (double) n);
        System.out.println("Original Message Sent = " + m);
    }
}
