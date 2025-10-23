#include <bits/stdc++.h>
using namespace std;

class SimplifiedAES {
public:
    static constexpr uint8_t sBox[16] = {
        0x9, 0x4, 0xA, 0xB,
        0xD, 0x1, 0x8, 0x5,
        0x6, 0x2, 0x0, 0x3,
        0xC, 0xE, 0xF, 0x7
    };

    static constexpr uint8_t sBoxI[16] = {
        0xA, 0x5, 0x9, 0xB,
        0x1, 0x7, 0x8, 0xF,
        0x6, 0x0, 0x2, 0x3,
        0xC, 0x4, 0xD, 0xE
    };

    vector<uint8_t> preRoundKey, round1Key, round2Key;

    SimplifiedAES(uint16_t key) {
        tie(preRoundKey, round1Key, round2Key) = keyExpansion(key);
    }

    static uint8_t subWord(uint8_t word) {
        return (sBox[(word >> 4)] << 4) | sBox[word & 0x0F];
    }

    static uint8_t rotWord(uint8_t word) {
        return ((word & 0x0F) << 4) | ((word & 0xF0) >> 4);
    }

    static tuple<vector<uint8_t>, vector<uint8_t>, vector<uint8_t>> keyExpansion(uint16_t key) {
        uint8_t Rcon1 = 0x80;
        uint8_t Rcon2 = 0x30;

        uint8_t w[6];
        w[0] = (key & 0xFF00) >> 8;
        w[1] = key & 0x00FF;

        w[2] = w[0] ^ (subWord(rotWord(w[1])) ^ Rcon1);
        w[3] = w[2] ^ w[1];
        w[4] = w[2] ^ (subWord(rotWord(w[3])) ^ Rcon2);
        w[5] = w[4] ^ w[3];

        auto k0 = intToState((w[0] << 8) + w[1]);
        auto k1 = intToState((w[2] << 8) + w[3]);
        auto k2 = intToState((w[4] << 8) + w[5]);

        return {k0, k1, k2};
    }

    static uint8_t gfMult(uint8_t a, uint8_t b) {
        uint8_t product = 0;
        a &= 0x0F;
        b &= 0x0F;

        while (b) {
            if (b & 1)
                product ^= a;
            a <<= 1;
            if (a & (1 << 4))
                a ^= 0b10011;
            b >>= 1;
        }
        return product;
    }

    static vector<uint8_t> intToState(uint16_t n) {
        return {
            uint8_t((n >> 12) & 0xF),
            uint8_t((n >> 4) & 0xF),
            uint8_t((n >> 8) & 0xF),
            uint8_t(n & 0xF)
        };
    }

    static uint16_t stateToInt(const vector<uint8_t> &m) {
        return (m[0] << 12) + (m[2] << 8) + (m[1] << 4) + m[3];
    }

    static vector<uint8_t> addRoundKey(const vector<uint8_t> &s1, const vector<uint8_t> &s2) {
        vector<uint8_t> res(4);
        for (int i = 0; i < 4; i++)
            res[i] = s1[i] ^ s2[i];
        return res;
    }

    static vector<uint8_t> subNibbles(const uint8_t sbox[16], const vector<uint8_t> &state) {
        vector<uint8_t> res(4);
        for (int i = 0; i < 4; i++)
            res[i] = sbox[state[i]];
        return res;
    }

    static vector<uint8_t> shiftRows(const vector<uint8_t> &state) {
        return {state[0], state[1], state[3], state[2]};
    }

    static vector<uint8_t> mixColumns(const vector<uint8_t> &state) {
        return {
            uint8_t(state[0] ^ gfMult(4, state[2])),
            uint8_t(state[1] ^ gfMult(4, state[3])),
            uint8_t(state[2] ^ gfMult(4, state[0])),
            uint8_t(state[3] ^ gfMult(4, state[1]))
        };
    }

    static vector<uint8_t> inverseMixColumns(const vector<uint8_t> &state) {
        return {
            uint8_t(gfMult(9, state[0]) ^ gfMult(2, state[2])),
            uint8_t(gfMult(9, state[1]) ^ gfMult(2, state[3])),
            uint8_t(gfMult(9, state[2]) ^ gfMult(2, state[0])),
            uint8_t(gfMult(9, state[3]) ^ gfMult(2, state[1]))
        };
    }

    uint16_t Encrypt(uint16_t plaintext) {
        auto state = addRoundKey(preRoundKey, intToState(plaintext));
        state = mixColumns(shiftRows(subNibbles(sBox, state)));
        state = addRoundKey(round1Key, state);
        state = shiftRows(subNibbles(sBox, state));
        state = addRoundKey(round2Key, state);
        return stateToInt(state);
    }

    uint16_t Decrypt(uint16_t ciphertext) {
        auto state = addRoundKey(round2Key, intToState(ciphertext));
        state = subNibbles(sBoxI, shiftRows(state));
        state = inverseMixColumns(addRoundKey(round1Key, state));
        state = subNibbles(sBoxI, shiftRows(state));
        state = addRoundKey(preRoundKey, state);
        return stateToInt(state);
    }
};


int main() {
    uint16_t key = 0b0100101011110101;
    uint16_t plaintext = 0b1101011100101000;

    SimplifiedAES saes(key);

    uint16_t enc = saes.Encrypt(plaintext);
    cout << "Encrypted: " << bitset<16>(enc) << endl;

    uint16_t dec = saes.Decrypt(enc);
    cout << "Decrypted: " << bitset<16>(dec) << endl;

    return 0;
}
