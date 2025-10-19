// To run this code, you need Node.js installed.
// Open a terminal in this folder and run: node diffie_hellman.js

/**
 * Power function to return value of (a^b) mod P.
 * This uses BigInt for handling potentially large numbers in cryptography.
 *
 * @param {bigint} a The base.
 * @param {bigint} b The exponent.
 * @param {bigint} P The modulus.
 * @returns {bigint} The result of (a^b) mod P.
 */
function power(a, b, P) {
    // Modular exponentiation
    let result = 1n;
    a = a % P;
    while (b > 0n) {
        if (b % 2n === 1n) {
            result = (result * a) % P;
        }
        b = b / 2n;
        a = (a * a) % P;
    }
    return result;
}

/**
 * A simple XOR-based encryption/decryption function.
 * Applying the same function twice with the same key returns the original message.
 *
 * @param {string} message The string message to be processed.
 * @param {bigint} key The secret key to use for the XOR operation.
 * @returns {string} The processed (encrypted or decrypted) string.
 */
function encryptDecrypt(message, key) {
    let result = '';
    // The key might be a BigInt, so we convert it to a Number for XOR.
    // This is safe for typical key sizes in this example.
    const keyNum = Number(key);
    for (let i = 0; i < message.length; i++) {
        result += String.fromCharCode(message.charCodeAt(i) ^ keyNum);
    }
    return result;
}

function main() {
    // We use 'n' at the end of numbers to denote them as BigInt literals.
    let P, G, x, a, y, b, ka, kb;

    // Both parties agree upon the public keys G and P.

    // A prime number P is taken.
    P = 23n;
    console.log(`The value of P: ${P}`);

    // A primitive root for P, G is taken.
    G = 9n;
    console.log(`The value of G: ${G}`);
    console.log("-------------------------");

    // Alice chooses her private key 'a'.
    a = 4n;
    console.log(`The private key 'a' for Alice: ${a}`);

    // Alice calculates her public key 'x'.
    x = power(G, a, P);
    console.log(`The public key 'x' for Alice: ${x}`);
    console.log("-------------------------");

    // Bob chooses his private key 'b'.
    b = 3n;
    console.log(`The private key 'b' for Bob: ${b}`);

    // Bob calculates his public key 'y'.
    y = power(G, b, P);
    console.log(`The public key 'y' for Bob: ${y}`);
    console.log("-------------------------");

    // After exchanging public keys, they generate the shared secret key.

    // Alice calculates the secret key.
    ka = power(y, a, P);

    // Bob calculates the secret key.
    kb = power(x, b, P);

    console.log(`Secret key for Alice is: ${ka}`);
    console.log(`Secret key for Bob is: ${kb}`);
    console.log("-------------------------");

    // Now they can use the shared secret key for encryption.
    const message = "Hello Bob!";
    console.log(`Original Message: ${message}`);

    // Alice encrypts the message with her key.
    const encryptedMessage = encryptDecrypt(message, ka);
    console.log(`Encrypted Message: ${encryptedMessage}`);

    // Bob decrypts the message with his key.
    const decryptedMessage = encryptDecrypt(encryptedMessage, kb);
    console.log(`Decrypted Message: ${decryptedMessage}`);
}

// Run the main function
main();