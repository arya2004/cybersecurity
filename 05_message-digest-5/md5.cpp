#include <bits/stdc++.h>
using namespace std;

static const uint32_t S[64] = {
    7,12,17,22, 7,12,17,22, 7,12,17,22, 7,12,17,22,
    5,9,14,20, 5,9,14,20, 5,9,14,20, 5,9,14,20,
    4,11,16,23, 4,11,16,23, 4,11,16,23, 4,11,16,23,
    6,10,15,21, 6,10,15,21, 6,10,15,21, 6,10,15,21
};

static uint32_t K[64];

static uint32_t leftRotate(uint32_t x, uint32_t c) {
    return (x << c) | (x >> (32 - c));
}

static vector<uint8_t> digest(const vector<uint8_t> &message) {
    uint32_t a0 = 0x67452301;
    uint32_t b0 = 0xefcdab89;
    uint32_t c0 = 0x98badcfe;
    uint32_t d0 = 0x10325476;

    uint64_t originalLenBits = (uint64_t)message.size() * 8ULL;
    size_t paddingLen = ((56 - (message.size() + 1) % 64) + 64) % 64;
    size_t totalLen = message.size() + 1 + paddingLen + 8;

    vector<uint8_t> padded(totalLen, 0);
    memcpy(padded.data(), message.data(), message.size());
    padded[message.size()] = 0x80;

    for (int i = 0; i < 8; i++) {
        padded[totalLen - 8 + i] = (uint8_t)((originalLenBits >> (8 * i)) & 0xFF);
    }

    size_t numChunks = padded.size() / 64;
    for (size_t i = 0; i < numChunks; i++) {
        uint32_t M[16];
        size_t offset = i * 64;
        for (int j = 0; j < 16; j++) {
            size_t idx = offset + j * 4;
            M[j] = (uint32_t)padded[idx]
                 | ((uint32_t)padded[idx + 1] << 8)
                 | ((uint32_t)padded[idx + 2] << 16)
                 | ((uint32_t)padded[idx + 3] << 24);
        }

        uint32_t A = a0, B = b0, C = c0, D = d0;

        for (int j = 0; j < 64; j++) {
            uint32_t F;
            uint32_t g;
            if (j <= 15) {
                F = (B & C) | (~B & D);
                g = (uint32_t)j;
            } else if (j <= 31) {
                F = (D & B) | (~D & C);
                g = (uint32_t)((5 * j + 1) % 16);
            } else if (j <= 47) {
                F = B ^ C ^ D;
                g = (uint32_t)((3 * j + 5) % 16);
            } else {
                F = C ^ (B | ~D);
                g = (uint32_t)((7 * j) % 16);
            }
            uint32_t temp = D;
            D = C;
            C = B;
            uint32_t sum = A + F + K[j] + M[g];
            B = B + leftRotate(sum, S[j]);
            A = temp;
        }

        a0 += A;
        b0 += B;
        c0 += C;
        d0 += D;
    }

    vector<uint8_t> out(16);
    auto putLE32 = [&](int pos, uint32_t v) {
        out[pos + 0] = (uint8_t)(v & 0xFF);
        out[pos + 1] = (uint8_t)((v >> 8) & 0xFF);
        out[pos + 2] = (uint8_t)((v >> 16) & 0xFF);
        out[pos + 3] = (uint8_t)((v >> 24) & 0xFF);
    };

    putLE32(0, a0);
    putLE32(4, b0);
    putLE32(8, c0);
    putLE32(12, d0);
    return out;
}

static string toHex(const vector<uint8_t> &bytes) {
    static const char *hex = "0123456789abcdef";
    string s;
    s.reserve(bytes.size() * 2);
    for (uint8_t b : bytes) {
        s.push_back(hex[(b >> 4) & 0xF]);
        s.push_back(hex[b & 0xF]);
    }
    return s;
}

int main() {
    for (int i = 0; i < 64; i++) {
        K[i] = (uint32_t)(uint64_t)( (1ULL << 32) * fabs(sin((double)(i + 1))) );
    }

    cout << "Enter text to hash using MD5: ";
    string input;
    getline(cin, input);

    vector<uint8_t> msg(input.begin(), input.end());
    vector<uint8_t> d = digest(msg);
    cout << "MD5 Hash: " << toHex(d) << "\n";
    return 0;
}
