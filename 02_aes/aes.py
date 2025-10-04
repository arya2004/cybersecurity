class SimplifiedAES:
    sBox = [
        0x9, 0x4, 0xA, 0xB,
        0xD, 0x1, 0x8, 0x5,
        0x6, 0x2, 0x0, 0x3,
        0xC, 0xE, 0xF, 0x7
    ]

    sBoxI = [
        0xA, 0x5, 0x9, 0xB,
        0x1, 0x7, 0x8, 0xF,
        0x6, 0x0, 0x2, 0x3,
        0xC, 0x4, 0xD, 0xE
    ]

    def __init__(self, key: int):
        self.preRoundKey, self.round1Key, self.round2Key = self.keyExpansion(key)


    def subWord(word: int) -> int:
        return (SimplifiedAES.sBox[(word >> 4)] << 4) + SimplifiedAES.sBox[word & 0x0F]

    def rotWord(word: int) -> int:
        return ((word & 0x0F) << 4) + ((word & 0xF0) >> 4)

    def keyExpansion(key: int):
        Rcon1 = 0x80
        Rcon2 = 0x30

        w = [0] * 6
        w[0] = (key & 0xFF00) >> 8
        w[1] = key & 0x00FF
        w[2] = w[0] ^ (SimplifiedAES.subWord(SimplifiedAES.rotWord(w[1])) ^ Rcon1)
        w[3] = w[2] ^ w[1]
        w[4] = w[2] ^ (SimplifiedAES.subWord(SimplifiedAES.rotWord(w[3])) ^ Rcon2)
        w[5] = w[4] ^ w[3]

        return (SimplifiedAES.intToState((w[0] << 8) + w[1]),
                SimplifiedAES.intToState((w[2] << 8) + w[3]),
                SimplifiedAES.intToState((w[4] << 8) + w[5]))


    def gfMult(a: int, b: int) -> int:
        product = 0
        a &= 0x0F
        b &= 0x0F
        while a != 0 and b != 0:
            if b & 1:
                product ^= a
            a <<= 1
            if a & (1 << 4):
                a ^= 0b10011
            b >>= 1
        return product


    def intToState(n: int):
        return [(n >> 12) & 0xF, (n >> 4) & 0xF, (n >> 8) & 0xF, n & 0xF]

    def stateToInt(m):
        return (m[0] << 12) + (m[2] << 8) + (m[1] << 4) + m[3]

    def addRoundKey(s1, s2):
        return [a ^ b for a, b in zip(s1, s2)]

    def subNibbles(sbox, state):
        return [sbox[x] for x in state]

    def shiftRows(state):
        return [state[0], state[1], state[3], state[2]]

    def mixColumns(state):
        return [
            state[0] ^ SimplifiedAES.gfMult(4, state[2]),
            state[1] ^ SimplifiedAES.gfMult(4, state[3]),
            state[2] ^ SimplifiedAES.gfMult(4, state[0]),
            state[3] ^ SimplifiedAES.gfMult(4, state[1]),
        ]
    
    def inverseMixColumns(state):
        return [
            SimplifiedAES.gfMult(9, state[0]) ^ SimplifiedAES.gfMult(2, state[2]),
            SimplifiedAES.gfMult(9, state[1]) ^ SimplifiedAES.gfMult(2, state[3]),
            SimplifiedAES.gfMult(9, state[2]) ^ SimplifiedAES.gfMult(2, state[0]),
            SimplifiedAES.gfMult(9, state[3]) ^ SimplifiedAES.gfMult(2, state[1]),
        ]

    def Encrypt(self, plaintext: int) -> int:
        state = self.addRoundKey(self.preRoundKey, self.intToState(plaintext))
        state = self.mixColumns(self.shiftRows(self.subNibbles(self.sBox, state)))
        state = self.addRoundKey(self.round1Key, state)
        state = self.shiftRows(self.subNibbles(self.sBox, state))
        state = self.addRoundKey(self.round2Key, state)
        return self.stateToInt(state)

    def Decrypt(self, ciphertext: int) -> int:
        state = self.addRoundKey(self.round2Key, self.intToState(ciphertext))
        state = self.subNibbles(self.sBoxI, self.shiftRows(state))
        state = self.inverseMixColumns(self.addRoundKey(self.round1Key, state))
        state = self.subNibbles(self.sBoxI, self.shiftRows(state))
        state = self.addRoundKey(self.preRoundKey, state)
        return self.stateToInt(state)


if __name__ == "__main__":
    key = 0b0100101011110101
    plaintext = 0b1101011100101000

    saes = SimplifiedAES(key)

    enc = saes.Encrypt(plaintext)
    print(f"Encrypted: {enc:016b}")

    dec = saes.Decrypt(enc)
    print(f"Decrypted: {dec:016b}")
