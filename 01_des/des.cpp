#include <bits/stdc++.h>
using namespace std;

/*
 * Simplified DES (S-DES) implementation in C++
 * ---------------------------------------------
 * This is a pedagogical version demonstrating the principles of block ciphers.
 * DO NOT use this for any real-world encryption.
 *
 * Author: <your-name>
 * Based on Java version by <original author if known>
 */

// --- Permutation and S-Box Constants ---
int P10[] = {3, 5, 2, 7, 4, 10, 1, 9, 8, 6};
int P8[]  = {6, 3, 7, 4, 8, 5, 10, 9};
int P4[]  = {2, 4, 3, 1};
int IP[]  = {2, 6, 3, 1, 4, 8, 5, 7};
int IP_INV[] = {4, 1, 3, 5, 7, 2, 8, 6};
int E_P[] = {4, 1, 2, 3, 2, 3, 4, 1};

int S0[4][4] = {
    {1, 0, 3, 2},
    {3, 2, 1, 0},
    {0, 2, 1, 3},
    {3, 1, 3, 2}
};

int S1[4][4] = {
    {0, 1, 2, 3},
    {2, 0, 1, 3},
    {3, 0, 1, 0},
    {2, 1, 0, 3}
};

// Utility functions
vector<int> permute(const vector<int> &arr, int perm[], int size) {
    vector<int> result(size);
    for (int i = 0; i < size; i++) {
        result[i] = arr[perm[i] - 1];
    }
    return result;
}

vector<int> leftShift(const vector<int> &bits, int shifts) {
    vector<int> shifted(bits.size());
    for (int i = 0; i < bits.size(); i++) {
        shifted[i] = bits[(i + shifts) % bits.size()];
    }
    return shifted;
}

vector<int> XOR(const vector<int> &a, const vector<int> &b) {
    vector<int> res(a.size());
    for (int i = 0; i < a.size(); i++) res[i] = a[i] ^ b[i];
    return res;
}

int binToInt(const vector<int> &bits) {
    int val = 0;
    for (int i = 0; i < bits.size(); i++) {
        val = (val << 1) | bits[i];
    }
    return val;
}

vector<int> intToBin(int val, int size) {
    vector<int> bits(size);
    for (int i = size - 1; i >= 0; i--) {
        bits[i] = val & 1;
        val >>= 1;
    }
    return bits;
}

vector<int> sBox(const vector<int> &input, int sMatrix[4][4]) {
    int row = (input[0] << 1) | input[3];
    int col = (input[1] << 1) | input[2];
    int val = sMatrix[row][col];
    return intToBin(val, 2);
}

// --- Key generation (produces K1, K2) ---
pair<vector<int>, vector<int>> generateKeys(vector<int> key) {
    vector<int> p10 = permute(key, P10, 10);
    vector<int> left(p10.begin(), p10.begin() + 5);
    vector<int> right(p10.begin() + 5, p10.end());

    // LS-1
    left = leftShift(left, 1);
    right = leftShift(right, 1);
    vector<int> combined = left;
    combined.insert(combined.end(), right.begin(), right.end());
    vector<int> k1 = permute(combined, P8, 8);

    // LS-2
    left = leftShift(left, 2);
    right = leftShift(right, 2);
    combined = left;
    combined.insert(combined.end(), right.begin(), right.end());
    vector<int> k2 = permute(combined, P8, 8);

    return {k1, k2};
}

// F-function
vector<int> functionF(const vector<int> &left, const vector<int> &right, const vector<int> &key) {
    vector<int> temp = permute(right, E_P, 8);
    temp = XOR(temp, key);

    vector<int> left4(temp.begin(), temp.begin() + 4);
    vector<int> right4(temp.begin() + 4, temp.end());
    vector<int> s0_out = sBox(left4, S0);
    vector<int> s1_out = sBox(right4, S1);

    vector<int> combined = s0_out;
    combined.insert(combined.end(), s1_out.begin(), s1_out.end());
    combined = permute(combined, P4, 4);

    vector<int> result = XOR(left, combined);
    vector<int> f_output = result;
    f_output.insert(f_output.end(), right.begin(), right.end());
    return f_output;
}

// DES main
vector<int> sdes(vector<int> input, vector<int> key, bool decrypt) {
    auto [k1, k2] = generateKeys(key);
    if (decrypt) swap(k1, k2);

    vector<int> ip = permute(input, IP, 8);
    vector<int> left(ip.begin(), ip.begin() + 4);
    vector<int> right(ip.begin() + 4, ip.end());

    vector<int> f1 = functionF(left, right, k1);
    vector<int> swapped_left(f1.begin() + 4, f1.end());
    vector<int> swapped_right(f1.begin(), f1.begin() + 4);

    vector<int> f2 = functionF(swapped_left, swapped_right, k2);
    return permute(f2, IP_INV, 8);
}

// Input parsing helper
vector<int> parseBits(const string &s, int expected) {
    if (s.size() != expected || s.find_first_not_of("01") != string::npos) {
        cout << "Invalid input. Must be " << expected << " bits (0 or 1 only).\n";
        return {};
    }
    vector<int> bits(expected);
    for (int i = 0; i < expected; i++) bits[i] = s[i] - '0';
    return bits;
}

void printBits(const vector<int> &bits) {
    for (int b : bits) cout << b;
    cout << "\n";
}

// --- Interactive Menu ---
int main() {
    vector<int> input, key;
    bool decrypt = false;
    while (true) {
        cout << "\nMenu:\n"
             << "1. Enter input array (8 bits)\n"
             << "2. Enter key array (10 bits)\n"
             << "3. Set mode (encrypt/decrypt)\n"
             << "4. Run S-DES\n"
             << "5. Exit\n"
             << "Enter choice: ";

        int choice;
        if (!(cin >> choice)) {
            cin.clear(); cin.ignore(1000, '\n');
            cout << "Invalid choice.\n";
            continue;
        }

        if (choice == 1) {
            string s; cout << "Enter 8-bit input: "; cin >> s;
            input = parseBits(s, 8);
        } else if (choice == 2) {
            string s; cout << "Enter 10-bit key: "; cin >> s;
            key = parseBits(s, 10);
        } else if (choice == 3) {
            string mode; cout << "Enter mode (encrypt/decrypt): "; cin >> mode;
            decrypt = (mode == "decrypt");
            cout << "Mode set to " << (decrypt ? "Decrypt" : "Encrypt") << ".\n";
        } else if (choice == 4) {
            if (input.empty() || key.empty()) {
                cout << "Please enter input and key first.\n";
            } else {
                vector<int> result = sdes(input, key, decrypt);
                cout << "Result: ";
                printBits(result);
            }
        } else if (choice == 5) {
            cout << "Exiting.\n";
            break;
        } else {
            cout << "Invalid option.\n";
        }
    }
    return 0;
}
