class SimplifiedDES {
    constructor() {
        this.P10 = [3, 5, 2, 7, 4, 10, 1, 9, 8, 6];
        this.P8 = [6, 3, 7, 4, 8, 5, 10, 9];
        this.P4 = [2, 4, 3, 1];
        this.IP = [2, 6, 3, 1, 4, 8, 5, 7];
        this.IP_INV = [4, 1, 3, 5, 7, 2, 8, 6];
        this.E_P = [4, 1, 2, 3, 2, 3, 4, 1];

        this.S0 = [
            [1, 0, 3, 2],
            [3, 2, 1, 0],
            [0, 2, 1, 3],
            [3, 1, 3, 2]
        ];

        this.S1 = [
            [0, 1, 2, 3],
            [2, 0, 1, 3],
            [3, 0, 1, 0],
            [2, 1, 0, 3]
        ];
    }

    permute(arr, perm, size) {
        const result = new Array(size);
        for (let i = 0; i < size; i++) {
            result[i] = arr[perm[i] - 1];
        }
        return result;
    }

    leftShift(bits, shifts) {
        const shifted = new Array(bits.length);
        for (let i = 0; i < bits.length; i++) {
            shifted[i] = bits[(i + shifts) % bits.length];
        }
        return shifted;
    }

    XOR(a, b) {
        const res = new Array(a.length);
        for (let i = 0; i < a.length; i++) res[i] = a[i] ^ b[i];
        return res;
    }

    binToInt(bits) {
        let val = 0;
        for (let i = 0; i < bits.length; i++) {
            val = (val << 1) | bits[i];
        }
        return val;
    }

    intToBin(val, size) {
        const bits = new Array(size);
        for (let i = size - 1; i >= 0; i--) {
            bits[i] = val & 1;
            val >>= 1;
        }
        return bits;
    }

    sBox(input, sMatrix) {
        const row = (input[0] << 1) | input[3];
        const col = (input[1] << 1) | input[2];
        const val = sMatrix[row][col];
        return this.intToBin(val, 2);
    }

    generateKeys(key) {
        const p10 = this.permute(key, this.P10, 10);
        let left = p10.slice(0, 5);
        let right = p10.slice(5);

        left = this.leftShift(left, 1);
        right = this.leftShift(right, 1);
        let combined = left.concat(right);
        const k1 = this.permute(combined, this.P8, 8);

        left = this.leftShift(left, 2);
        right = this.leftShift(right, 2);
        combined = left.concat(right);
        const k2 = this.permute(combined, this.P8, 8);

        return [k1, k2];
    }

    functionF(left, right, key) {
        let temp = this.permute(right, this.E_P, 8);
        temp = this.XOR(temp, key);

        const left4 = temp.slice(0, 4);
        const right4 = temp.slice(4);
        const s0_out = this.sBox(left4, this.S0);
        const s1_out = this.sBox(right4, this.S1);

        let combined = s0_out.concat(s1_out);
        combined = this.permute(combined, this.P4, 4);

        const result = this.XOR(left, combined);
        const f_output = result.concat(right);
        return f_output;
    }

    sdes(input, key, decrypt) {
        let [k1, k2] = this.generateKeys(key);
        if (decrypt) [k1, k2] = [k2, k1];

        const ip = this.permute(input, this.IP, 8);
        const left = ip.slice(0, 4);
        const right = ip.slice(4);

        const f1 = this.functionF(left, right, k1);
        const swapped_left = f1.slice(4);
        const swapped_right = f1.slice(0, 4);

        const f2 = this.functionF(swapped_left, swapped_right, k2);
        return this.permute(f2, this.IP_INV, 8);
    }

    parseBits(s, expected) {
        if (s.length !== expected || /[^01]/.test(s)) {
            console.log(`Invalid input. Must be ${expected} bits (0 or 1 only).`);
            return [];
        }
        const bits = new Array(expected);
        for (let i = 0; i < expected; i++) bits[i] = s.charCodeAt(i) - 48;
        return bits;
    }

    printBits(bits) {
        console.log(bits.join(""));
    }
}

function main() {
    const sdes = new SimplifiedDES();

    const input = sdes.parseBits("10101010", 8);
    const key = sdes.parseBits("1010000010", 10);

    const encrypted = sdes.sdes(input, key, false);
    console.log("Encrypted: " + encrypted.join(""));

    const decrypted = sdes.sdes(encrypted, key, true);
    console.log("Decrypted: " + decrypted.join(""));
}

if (typeof module !== 'undefined' && module.exports) {
    module.exports = SimplifiedDES;
} else {
    main();
}
