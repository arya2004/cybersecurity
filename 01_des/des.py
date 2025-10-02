import math


def permutation(bits, positions):
    return [bits[p-1] for p in positions]

def left_shift(bits, n):
    return bits[n:] + bits[:n]

def xor(bits1, bits2):
    return [b1 ^ b2 for b1, b2 in zip(bits1, bits2)]

def bin_to_int(bits):
    val = 0
    for i, bit in enumerate(reversed(bits)):
        val += bit * (2 ** i)
    return val

def sbox(bits, s_matrix):
    row = bin_to_int([bits[0], bits[3]])
    col = bin_to_int([bits[1], bits[2]])
    val = s_matrix[row][col]
    return [int(x) for x in f"{val:02b}"]

def generate_keys(key):
    P10 = [3,5,2,7,4,10,1,9,8,6]
    P8  = [6,3,7,4,8,5,10,9]
    LS1 = 1
    LS2 = 2

    key = permutation(key, P10)
    left, right = key[:5], key[5:]

    left1, right1 = left_shift(left, LS1), left_shift(right, LS1)
    k1 = permutation(left1 + right1, P8)

    left2, right2 = left_shift(left1, LS2), left_shift(right1, LS2)
    k2 = permutation(left2 + right2, P8)

    return k1, k2

def function_f(left, right, key):
    EP = [4,1,2,3,2,3,4,1]
    P4 = [2,4,3,1]
    S0 = [
        [1,0,3,2],
        [3,2,1,0],
        [0,2,1,3],
        [3,1,3,2]
    ]
    S1 = [
        [0,1,2,3],
        [2,0,1,3],
        [3,0,1,0],
        [2,1,0,3]
    ]
    temp = permutation(right, EP)
    temp = xor(temp, key)
    left_s0, right_s1 = temp[:4], temp[4:]
    left_s0 = sbox(left_s0, S0)
    right_s1 = sbox(right_s1, S1)
    f_out = permutation(left_s0 + right_s1, P4)
    return xor(left, f_out) + right

def swap(left, right):
    return right, left

def sdes(input_bits, key_bits, decrypt=False):
    IP = [2,6,3,1,4,8,5,7]
    IP_inv = [4,1,3,5,7,2,8,6]

    k1, k2 = generate_keys(key_bits)
    if decrypt:
        k1, k2 = k2, k1

    permuted = permutation(input_bits, IP)
    left, right = permuted[:4], permuted[4:]
    f_out = function_f(left, right, k1)
    left, right = f_out[:4], f_out[4:]
    left, right = swap(left, right)
    f_out = function_f(left, right, k2)
    output_bits = permutation(f_out, IP_inv)
    return output_bits


def parse_input_array(input_str):
    try:
        return [int(c) for c in input_str if c in '01']
    except ValueError:
        print("Invalid input! Use only 0 and 1.")
        return []

def parse_bool(input_str):
    return input_str.lower() == "true"


def main():
    input_bits = []
    key_bits = []
    decrypt_mode = False

    while True:
        print("\nMenu:")
        print("1. Enter input array (8 bits, e.g., 10001011)")
        print("2. Enter key array (10 bits, e.g., 1001100111)")
        print("3. Set decrypt mode (true/false)")
        print("4. Call S-DES function")
        print("5. Exit")
        choice = input("Enter your choice: ")

        if choice == "1":
            inp = input("Enter input array: ")
            input_bits = parse_input_array(inp)
        elif choice == "2":
            key = input("Enter key array: ")
            key_bits = parse_input_array(key)
        elif choice == "3":
            val = input("Enter decrypt mode (true/false): ")
            decrypt_mode = parse_bool(val)
        elif choice == "4":
            if len(input_bits) != 8 or len(key_bits) != 10:
                print("Input must be 8 bits and key must be 10 bits.")
            else:
                result = sdes(input_bits, key_bits, decrypt_mode)
                print("Result:", result)
        elif choice == "5":
            print("Exiting program.")
            break
        else:
            print("Invalid choice. Try again.")

if __name__ == "__main__":
    main()
