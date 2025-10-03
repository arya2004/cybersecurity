class SimplifiedAES {
    constructor(key) {
        this.sBox = [
            0x9, 0x4, 0xA, 0xB, 0xD, 0x1, 0x8, 0x5, 
            0x6, 0x2, 0x0, 0x3, 0xC, 0xE, 0xF, 0x7
        ];
        
        this.sBoxI = [
            0xA, 0x5, 0x9, 0xB, 0x1, 0x7, 0x8, 0xF, 
            0x6, 0x0, 0x2, 0x3, 0xC, 0x4, 0xD, 0xE
        ];
        
        const [preRoundKey, round1Key, round2Key] = this.keyExpansion(key);
        this.preRoundKey = preRoundKey;
        this.round1Key = round1Key;
        this.round2Key = round2Key;
    }
    
    subWord(word) {
        return (this.sBox[word >> 4] << 4) + this.sBox[word & 0x0F];
    }
    
    rotWord(word) {
        return ((word & 0x0F) << 4) + ((word & 0xF0) >> 4);
    }
    
    keyExpansion(key) {
        const Rcon1 = 0x80;
        const Rcon2 = 0x30;
        
        const w = new Array(6);
        w[0] = (key & 0xFF00) >> 8;
        w[1] = key & 0x00FF;
        w[2] = w[0] ^ (this.subWord(this.rotWord(w[1])) ^ Rcon1);
        w[3] = w[2] ^ w[1];
        w[4] = w[2] ^ (this.subWord(this.rotWord(w[3])) ^ Rcon2);
        w[5] = w[4] ^ w[3];
        
        return [
            this.intToState((w[0] << 8) + w[1]),
            this.intToState((w[2] << 8) + w[3]),
            this.intToState((w[4] << 8) + w[5])
        ];
    }
    
    gfMult(a, b) {
        let product = 0;
        a = a & 0x0F;
        b = b & 0x0F;
        
        while (a !== 0 && b !== 0) {
            if (b & 1) {
                product = product ^ a;
            }
            a = a << 1;
            if (a & (1 << 4)) {
                a = a ^ 0b10011;
            }
            b = b >> 1;
        }
        
        return product;
    }
    
    intToState(n) {
        return [
            (n >> 12) & 0xF,
            (n >> 4) & 0xF,
            (n >> 8) & 0xF,
            n & 0xF
        ];
    }
    
    stateToInt(m) {
        return (m[0] << 12) + (m[2] << 8) + (m[1] << 4) + m[3];
    }
    
    addRoundKey(s1, s2) {
        return s1.map((val, i) => val ^ s2[i]);
    }
    
    subNibbles(sbox, state) {
        return state.map(val => sbox[val]);
    }
    
    shiftRows(state) {
        return [state[0], state[1], state[3], state[2]];
    }
    
    mixColumns(state) {
        return [
            state[0] ^ this.gfMult(4, state[2]),
            state[1] ^ this.gfMult(4, state[3]),
            state[2] ^ this.gfMult(4, state[0]),
            state[3] ^ this.gfMult(4, state[1])
        ];
    }
    
    inverseMixColumns(state) {
        return [
            this.gfMult(9, state[0]) ^ this.gfMult(2, state[2]),
            this.gfMult(9, state[1]) ^ this.gfMult(2, state[3]),
            this.gfMult(9, state[2]) ^ this.gfMult(2, state[0]),
            this.gfMult(9, state[3]) ^ this.gfMult(2, state[1])
        ];
    }
    
    encrypt(plaintext) {
        let state = this.addRoundKey(this.preRoundKey, this.intToState(plaintext));
        state = this.mixColumns(this.shiftRows(this.subNibbles(this.sBox, state)));
        state = this.addRoundKey(this.round1Key, state);
        state = this.shiftRows(this.subNibbles(this.sBox, state));
        state = this.addRoundKey(this.round2Key, state);
        return this.stateToInt(state);
    }
    
    decrypt(ciphertext) {
        let state = this.addRoundKey(this.round2Key, this.intToState(ciphertext));
        state = this.subNibbles(this.sBoxI, this.shiftRows(state));
        state = this.inverseMixColumns(this.addRoundKey(this.round1Key, state));
        state = this.subNibbles(this.sBoxI, this.shiftRows(state));
        state = this.addRoundKey(this.preRoundKey, state);
        return this.stateToInt(state);
    }
}

function main() {
    const key = 0b0100101011110101;
    const plaintext = 0b1101011100101000;
    
    const saes = new SimplifiedAES(key);
    
    const encrypted = saes.encrypt(plaintext);
    console.log(`Encrypted: ${encrypted.toString(2).padStart(16, '0')}`);
    
    const decrypted = saes.decrypt(encrypted);
    console.log(`Decrypted: ${decrypted.toString(2).padStart(16, '0')}`);
}

if (typeof module !== 'undefined' && module.exports) {
    module.exports = SimplifiedAES;
} else {
    main();
}
