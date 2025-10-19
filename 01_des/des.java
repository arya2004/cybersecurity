// To compile and run this code from the command line, navigate to the 01_des folder and use:
// javac des.java && java des

import java.util.Arrays;
import java.util.Scanner;

/**
 * An educational implementation of the Simplified DES (S-DES) algorithm in Java.
 * This version is a direct translation of the provided Go implementation's structure and logic.
 *
 * S-DES operates on 8-bit blocks of data and uses a 10-bit key. It is intended
 * for academic purposes to demonstrate the principles of modern block ciphers.
 *
 * WARNING: This algorithm is NOT secure and should not be used for real-world applications.
 */
public class des {

    // --- S-DES Permutation and S-Box Constants ---
    private static final int[] P10 = {3, 5, 2, 7, 4, 10, 1, 9, 8, 6};
    private static final int[] P8 = {6, 3, 7, 4, 8, 5, 10, 9};
    private static final int[] P4 = {2, 4, 3, 1};
    private static final int[] IP = {2, 6, 3, 1, 4, 8, 5, 7};
    private static final int[] IP_INV = {4, 1, 3, 5, 7, 2, 8, 6};
    private static final int[] E_P = {4, 1, 2, 3, 2, 3, 4, 1};

    private static final int[][] S0 = {
            {1, 0, 3, 2},
            {3, 2, 1, 0},
            {0, 2, 1, 3},
            {3, 1, 3, 2}
    };

    private static final int[][] S1 = {
            {0, 1, 2, 3},
            {2, 0, 1, 3},
            {3, 0, 1, 0},
            {2, 1, 0, 3}
    };

    /**
     * The main S-DES encryption/decryption function.
     * @param inputText 8-bit input data.
     * @param key 10-bit key.
     * @param decrypt A boolean flag to switch between encryption and decryption.
     * @return 8-bit output data.
     */
    public static int[] des(int[] inputText, int[] key, boolean decrypt) {
        int[][] keys = generateKeys(key);
        int[] k1 = decrypt ? keys[1] : keys[0];
        int[] k2 = decrypt ? keys[0] : keys[1];

        // 1. Initial Permutation (IP)
        int[] permutedText = permutation(inputText, IP);
        int[] left = Arrays.copyOfRange(permutedText, 0, 4);
        int[] right = Arrays.copyOfRange(permutedText, 4, 8);

        // 2. First Round with F_k1
        int[] f1_output = functionF(left, right, k1);

        // 3. Swap (SW)
        int[] swapped_left = Arrays.copyOfRange(f1_output, 4, 8);
        int[] swapped_right = Arrays.copyOfRange(f1_output, 0, 4);

        // 4. Second Round with F_k2
        int[] f2_output = functionF(swapped_left, swapped_right, k2);

        // 5. Final Permutation (IP^-1)
        return permutation(f2_output, IP_INV);
    }

    /**
     * Performs the complex F function on the left and right halves.
     * @param left 4-bit left half.
     * @param right 4-bit right half.
     * @param key 8-bit subkey.
     * @return A combined 8-bit array after the operations.
     */
    private static int[] functionF(int[] left, int[] right, int[] key) {
        // Expansion/Permutation on the right half
        int[] tempValue = permutation(right, E_P);

        // XOR with the subkey
        tempValue = xor(tempValue, key);

        // S-Box substitutions
        int[] sBox0_input = Arrays.copyOfRange(tempValue, 0, 4);
        int[] sBox1_input = Arrays.copyOfRange(tempValue, 4, 8);
        int[] sBox0_output = sBox(sBox0_input, S0);
        int[] sBox1_output = sBox(sBox1_input, S1);

        // Combine S-Box outputs
        int[] partialOutput = new int[4];
        System.arraycopy(sBox0_output, 0, partialOutput, 0, 2);
        System.arraycopy(sBox1_output, 0, partialOutput, 2, 2);

        // P4 Permutation
        partialOutput = permutation(partialOutput, P4);

        // XOR with the left half
        partialOutput = xor(left, partialOutput);

        // Combine the result with the original right half
        int[] fOutput = new int[8];
        System.arraycopy(partialOutput, 0, fOutput, 0, 4);
        System.arraycopy(right, 0, fOutput, 4, 4);

        return fOutput;
    }

    /**
     * S-Box lookup function.
     * @param bitList 4-bit input.
     * @param sMatrix The S-Box (S0 or S1) to use.
     * @return 2-bit output.
     */
    private static int[] sBox(int[] bitList, int[][] sMatrix) {
        int row = binToInt(new int[]{bitList[0], bitList[3]});
        int column = binToInt(new int[]{bitList[1], bitList[2]});
        int value = sMatrix[row][column];

        // Convert the integer value (0-3) to a 2-bit array
        int[] output = new int[2];
        output[0] = (value & 2) >> 1; // First bit
        output[1] = value & 1;        // Second bit
        return output;
    }

    /**
     * Generates two 8-bit subkeys (k1, k2) from a 10-bit key.
     * @param key The 10-bit initial key.
     * @return A 2D array containing the two subkeys.
     */
    private static int[][] generateKeys(int[] key) {
        // P10
        int[] p10_key = permutation(key, P10);
        int[] leftHalf = Arrays.copyOfRange(p10_key, 0, 5);
        int[] rightHalf = Arrays.copyOfRange(p10_key, 5, 10);

        // LS-1
        int[] ls1_left = leftShift(leftHalf, 1);
        int[] ls1_right = leftShift(rightHalf, 1);
        int[] combined_ls1 = new int[10];
        System.arraycopy(ls1_left, 0, combined_ls1, 0, 5);
        System.arraycopy(ls1_right, 0, combined_ls1, 5, 5);

        // k1 from P8
        int[] k1 = permutation(combined_ls1, P8);

        // LS-2
        int[] ls2_left = leftShift(ls1_left, 2);
        int[] ls2_right = leftShift(ls1_right, 2);
        int[] combined_ls2 = new int[10];
        System.arraycopy(ls2_left, 0, combined_ls2, 0, 5);
        System.arraycopy(ls2_right, 0, combined_ls2, 5, 5);

        // k2 from P8
        int[] k2 = permutation(combined_ls2, P8);

        return new int[][]{k1, k2};
    }

    // --- UTILITY METHODS ---

    private static int[] permutation(int[] list, int[] positions) {
        int[] permutedList = new int[positions.length];
        for (int i = 0; i < positions.length; i++) {
            permutedList[i] = list[positions[i] - 1]; // Adjust for 0-based indexing
        }
        return permutedList;
    }

    private static int[] leftShift(int[] bits, int shifts) {
        int[] shifted = new int[bits.length];
        for (int i = 0; i < bits.length; i++) {
            shifted[i] = bits[(i + shifts) % bits.length];
        }
        return shifted;
    }

    private static int[] xor(int[] a, int[] b) {
        int[] result = new int[a.length];
        for (int i = 0; i < a.length; i++) {
            result[i] = a[i] ^ b[i];
        }
        return result;
    }

    private static int binToInt(int[] binValue) {
        int intValue = 0;
        for (int i = 0; i < binValue.length; i++) {
            intValue += binValue[binValue.length - 1 - i] * Math.pow(2, i);
        }
        return intValue;
    }

    // --- INTERACTIVE MENU ---

    public static void main(String[] args) {
        int[] input = null;
        int[] key = null;
        boolean decrypt = false;
        Scanner scanner = new Scanner(System.in);

        while (true) {
            System.out.println("\nMenu:");
            System.out.println("1. Enter input array (8 bits)");
            System.out.println("2. Enter key array (10 bits)");
            System.out.println("3. Set mode (encrypt/decrypt)");
            System.out.println("4. Call S-DES function");
            System.out.println("5. Exit");
            System.out.print("Enter your choice: ");

            int choice = -1;
            try {
                choice = Integer.parseInt(scanner.nextLine());
            } catch (NumberFormatException e) {
                System.out.println("Invalid choice. Please enter a number.");
                continue;
            }

            switch (choice) {
                case 1:
                    System.out.print("Enter input array (e.g., 10001011): ");
                    input = parseArray(scanner.nextLine(), 8);
                    break;
                case 2:
                    System.out.print("Enter key array (e.g., 1001100111): ");
                    key = parseArray(scanner.nextLine(), 10);
                    break;
                case 3:
                    System.out.print("Enter mode ('encrypt' or 'decrypt'): ");
                    String mode = scanner.nextLine().toLowerCase();
                    if (mode.equals("decrypt")) {
                        decrypt = true;
                        System.out.println("Mode set to Decrypt.");
                    } else {
                        decrypt = false;
                        System.out.println("Mode set to Encrypt.");
                    }
                    break;
                case 4:
                    if (input == null || key == null) {
                        System.out.println("Input or key array is empty. Please enter both before calling S-DES.");
                    } else {
                        System.out.println("Current mode: " + (decrypt ? "Decrypt" : "Encrypt"));
                        int[] result = des(input, key, decrypt);
                        System.out.print("Result from S-DES function: ");
                        for (int bit : result) {
                            System.out.print(bit);
                        }
                        System.out.println();
                    }
                    break;
                case 5:
                    System.out.println("Exiting program.");
                    scanner.close();
                    return;
                default:
                    System.out.println("Invalid choice. Please try again.");
            }
        }
    }

    private static int[] parseArray(String inputString, int expectedLength) {
        if (!inputString.matches("[01]+") || inputString.length() != expectedLength) {
            System.out.println("Invalid input. Please enter a binary string of length " + expectedLength + ".");
            return null;
        }
        int[] result = new int[expectedLength];
        for (int i = 0; i < expectedLength; i++) {
            result[i] = Character.getNumericValue(inputString.charAt(i));
        }
        System.out.println("Successfully parsed: " + Arrays.toString(result));
        return result;
    }
}