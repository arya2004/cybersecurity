import java.util.Arrays;

public class SimplifiedAES {

    private static final int[] sBox = {
        0x9, 0x4, 0xA, 0xB,
        0xD, 0x1, 0x8, 0x5,
        0x6, 0x2, 0x0, 0x3,
        0xC, 0xE, 0xF, 0x7
    };

    private static final int[] sBoxI = {
        0xA, 0x5, 0x9, 0xB,
        0x1, 0x7, 0x8, 0xF,
        0x6, 0x0, 0x2, 0x3,
        0xC, 0x4, 0xD, 0xE
    };

    private final int[] preRoundKey;
    private final int[] round1Key;
    private final int[] round2Key;

    public SimplifiedAES(int key) {
        int[][] keys = keyExpansion(key & 0xFFFF);
        this.preRoundKey = keys[0];
        this.round1Key = keys[1];
        this.round2Key = keys[2];
    }

    private static int subWord(int word) {
        int hi = (word >> 4) & 0xF;
        int lo = word & 0xF;
        return ((sBox[hi] << 4) | sBox[lo]) & 0xFF;
    }

    private static int rotWord(int word) {
        return (((word & 0x0F) << 4) | ((word & 0xF0) >> 4)) & 0xFF;
    }

    private static int[][] keyExpansion(int key) {
        int Rcon1 = 0x80;
        int Rcon2 = 0x30;

        int[] w = new int[6];
        w[0] = (key >> 8) & 0xFF;
        w[1] = key & 0xFF;

        w[2] = w[0] ^ (subWord(rotWord(w[1])) ^ Rcon1);
        w[3] = w[2] ^ w[1];
        w[4] = w[2] ^ (subWord(rotWord(w[3])) ^ Rcon2);
        w[5] = w[4] ^ w[3];

        int k0 = ((w[0] << 8) | w[1]) & 0xFFFF;
        int k1 = ((w[2] << 8) | w[3]) & 0xFFFF;
        int k2 = ((w[4] << 8) | w[5]) & 0xFFFF;

        return new int[][] { intToState(k0), intToState(k1), intToState(k2) };
    }

    private static int gfMult(int a, int b) {
        int product = 0;
        a &= 0xF;
        b &= 0xF;

        while (b != 0) {
            if ((b & 1) != 0) product ^= a;
            a <<= 1;
            if ((a & (1 << 4)) != 0) a ^= 0b10011;
            b >>= 1;
        }
        return product & 0xF;
    }

    private static int[] intToState(int n) {
        return new int[] {
            (n >> 12) & 0xF,
            (n >>  4) & 0xF,
            (n >>  8) & 0xF,
            (n      ) & 0xF
        };
    }

    private static int stateToInt(int[] s) {
        return ((s[0] << 12) | (s[2] << 8) | (s[1] << 4) | s[3]) & 0xFFFF;
    }

    private static int[] addRoundKey(int[] a, int[] b) {
        int[] r = new int[4];
        for (int i = 0; i < 4; i++) r[i] = (a[i] ^ b[i]) & 0xF;
        return r;
    }

    private static int[] subNibbles(int[] sbox, int[] state) {
        int[] r = new int[4];
        for (int i = 0; i < 4; i++) r[i] = sbox[state[i] & 0xF] & 0xF;
        return r;
    }

    private static int[] shiftRows(int[] state) {
        return new int[] { state[0], state[1], state[3], state[2] };
    }

    private static int[] mixColumns(int[] s) {
        return new int[] {
            (s[0] ^ gfMult(4, s[2])) & 0xF,
            (s[1] ^ gfMult(4, s[3])) & 0xF,
            (s[2] ^ gfMult(4, s[0])) & 0xF,
            (s[3] ^ gfMult(4, s[1])) & 0xF
        };
    }

    private static int[] inverseMixColumns(int[] s) {
        return new int[] {
            (gfMult(9, s[0]) ^ gfMult(2, s[2])) & 0xF,
            (gfMult(9, s[1]) ^ gfMult(2, s[3])) & 0xF,
            (gfMult(9, s[2]) ^ gfMult(2, s[0])) & 0xF,
            (gfMult(9, s[3]) ^ gfMult(2, s[1])) & 0xF
        };
    }

    public int encrypt(int plaintext) {
        int[] state = addRoundKey(preRoundKey, intToState(plaintext & 0xFFFF));
        state = mixColumns(shiftRows(subNibbles(sBox, state)));
        state = addRoundKey(round1Key, state);
        state = shiftRows(subNibbles(sBox, state));
        state = addRoundKey(round2Key, state);
        return stateToInt(state);
    }

    public int decrypt(int ciphertext) {
        int[] state = addRoundKey(round2Key, intToState(ciphertext & 0xFFFF));
        state = subNibbles(sBoxI, shiftRows(state));
        state = inverseMixColumns(addRoundKey(round1Key, state));
        state = subNibbles(sBoxI, shiftRows(state));
        state = addRoundKey(preRoundKey, state);
        return stateToInt(state);
    }

    private static String toBinary16BitString(int x) {
        String s = Integer.toBinaryString(x & 0xFFFF);
        return "0".repeat(16 - s.length()) + s;
    }

    public static void main(String[] args) {
        int key = 0b0100101011110101;
        int plaintext = 0b1101011100101000;

        SimplifiedAES saes = new SimplifiedAES(key);

        int enc = saes.encrypt(plaintext);
        System.out.println("Encrypted: " + toBinary16BitString(enc));

        int dec = saes.decrypt(enc);
        System.out.println("Decrypted: " + toBinary16BitString(dec));
    }
}
