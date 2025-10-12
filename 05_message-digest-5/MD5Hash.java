import java.nio.ByteBuffer;
import java.nio.ByteOrder;
import java.util.Formatter;
import java.util.Scanner;

public final class MD5Hash {
    private static final int[] S = {
        7,12,17,22, 7,12,17,22, 7,12,17,22, 7,12,17,22,
        5,9,14,20, 5,9,14,20, 5,9,14,20, 5,9,14,20,
        4,11,16,23, 4,11,16,23, 4,11,16,23, 4,11,16,23,
        6,10,15,21, 6,10,15,21, 6,10,15,21, 6,10,15,21
    };

    private static final int[] K = new int[64];
    static {
        for (int i = 0; i < 64; i++) {
            K[i] = (int) (long) ((1L << 32) * Math.abs(Math.sin(i + 1)));
        }
    }

    private static int leftRotate(int x, int c) {
        return (x << c) | (x >>> (32 - c));
    }

    public static byte[] digest(byte[] message) {
        int a0 = 0x67452301;
        int b0 = 0xefcdab89;
        int c0 = 0x98badcfe;
        int d0 = 0x10325476;

        int originalLenBits = message.length * 8;
        int paddingLen = ((56 - (message.length + 1) % 64) + 64) % 64;
        int totalLen = message.length + 1 + paddingLen + 8;
        byte[] padded = new byte[totalLen];

        System.arraycopy(message, 0, padded, 0, message.length);
        padded[message.length] = (byte) 0x80;

        ByteBuffer buf = ByteBuffer.allocate(8).order(ByteOrder.LITTLE_ENDIAN);
        buf.putLong(Integer.toUnsignedLong(originalLenBits));
        byte[] lenBytes = buf.array();
        System.arraycopy(lenBytes, 0, padded, totalLen - 8, 8);

        int numChunks = padded.length / 64;
        for (int i = 0; i < numChunks; i++) {
            int[] M = new int[16];
            int offset = i * 64;
            for (int j = 0; j < 16; j++) {
                int idx = offset + j * 4;
                M[j] = ((padded[idx] & 0xff)) |
                       ((padded[idx + 1] & 0xff) << 8) |
                       ((padded[idx + 2] & 0xff) << 16) |
                       ((padded[idx + 3] & 0xff) << 24);
            }

            int A = a0, B = b0, C = c0, D = d0;

            for (int j = 0; j < 64; j++) {
                int F, g;
                if (j <= 15) {
                    F = (B & C) | (~B & D);
                    g = j;
                } else if (j <= 31) {
                    F = (D & B) | (~D & C);
                    g = (5 * j + 1) % 16;
                } else if (j <= 47) {
                    F = B ^ C ^ D;
                    g = (3 * j + 5) % 16;
                } else {
                    F = C ^ (B | ~D);
                    g = (7 * j) % 16;
                }
                int temp = D;
                D = C;
                C = B;
                int sum = A + F + K[j] + M[g];
                B = B + leftRotate(sum, S[j]);
                A = temp;
            }

            a0 += A;
            b0 += B;
            c0 += C;
            d0 += D;
        }

        ByteBuffer out = ByteBuffer.allocate(16).order(ByteOrder.LITTLE_ENDIAN);
        out.putInt(a0);
        out.putInt(b0);
        out.putInt(c0);
        out.putInt(d0);
        return out.array();
    }

    public static String toHex(byte[] bytes) {
        try (Formatter fmt = new Formatter()) {
            for (byte b : bytes) fmt.format("%02x", b & 0xff);
            return fmt.toString();
        }
    }

    public static void main(String[] args) throws Exception {
        Scanner sc = new Scanner(System.in);
        System.out.print("Enter text to hash using MD5: ");
        String input = sc.nextLine();
        sc.close();

        byte[] digest = digest(input.getBytes("UTF-8"));
        System.out.println("MD5 Hash: " + toHex(digest));
    }
}
