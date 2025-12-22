class MD5Hash {
    constructor() {
        this.S = [
            7,12,17,22, 7,12,17,22, 7,12,17,22, 7,12,17,22,
            5,9,14,20, 5,9,14,20, 5,9,14,20, 5,9,14,20,
            4,11,16,23, 4,11,16,23, 4,11,16,23, 4,11,16,23,
            6,10,15,21, 6,10,15,21, 6,10,15,21, 6,10,15,21
        ];

        this.K = new Array(64);
        for (let i = 0; i < 64; i++) {
            this.K[i] = Math.floor((2 ** 32) * Math.abs(Math.sin(i + 1))) >>> 0;
        }
    }

    leftRotate(x, c) {
        return ((x << c) | (x >>> (32 - c))) >>> 0;
    }

    digest(message) {
        let a0 = 0x67452301 >>> 0;
        let b0 = 0xefcdab89 >>> 0;
        let c0 = 0x98badcfe >>> 0;
        let d0 = 0x10325476 >>> 0;

        const originalLenBits = (message.length * 8) >>> 0;
        const paddingLen = ((56 - ((message.length + 1) % 64) + 64) % 64) >>> 0;
        const totalLen = message.length + 1 + paddingLen + 8;

        const padded = new Uint8Array(totalLen);
        padded.set(message, 0);
        padded[message.length] = 0x80;

        let len = BigInt(message.length) * 8n;
        for (let i = 0; i < 8; i++) {
            padded[totalLen - 8 + i] = Number((len >> BigInt(8 * i)) & 0xFFn);
        }

        const numChunks = padded.length / 64;
        for (let i = 0; i < numChunks; i++) {
            const M = new Uint32Array(16);
            const offset = i * 64;
            for (let j = 0; j < 16; j++) {
                const idx = offset + j * 4;
                M[j] =
                    (padded[idx] & 0xff) |
                    ((padded[idx + 1] & 0xff) << 8) |
                    ((padded[idx + 2] & 0xff) << 16) |
                    ((padded[idx + 3] & 0xff) << 24);
            }

            let A = a0, B = b0, C = c0, D = d0;

            for (let j = 0; j < 64; j++) {
                let F, g;
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
                F >>>= 0;

                const temp = D;
                D = C;
                C = B;

                const sum = (A + F + this.K[j] + M[g]) >>> 0;
                B = (B + this.leftRotate(sum, this.S[j])) >>> 0;
                A = temp;
            }

            a0 = (a0 + A) >>> 0;
            b0 = (b0 + B) >>> 0;
            c0 = (c0 + C) >>> 0;
            d0 = (d0 + D) >>> 0;
        }

        const out = new Uint8Array(16);
        const dv = new DataView(out.buffer);
        dv.setUint32(0, a0, true);
        dv.setUint32(4, b0, true);
        dv.setUint32(8, c0, true);
        dv.setUint32(12, d0, true);
        return out;
    }

    toHex(bytes) {
        let s = "";
        for (let i = 0; i < bytes.length; i++) {
            s += bytes[i].toString(16).padStart(2, "0");
        }
        return s;
    }
}

async function main() {
    const readline = await import("node:readline/promises");
    const { stdin: inputStream, stdout: outputStream } = await import("node:process");
    const rl = readline.createInterface({ input: inputStream, output: outputStream });

    outputStream.write("Enter text to hash using MD5: ");
    const input = await rl.question("");
    rl.close();

    const md5 = new MD5Hash();
    const msg = new TextEncoder().encode(input);
    const digest = md5.digest(msg);
    outputStream.write("MD5 Hash: " + md5.toHex(digest) + "\n");
}

if (typeof module !== "undefined" && module.exports) {
    module.exports = MD5Hash;
} else {
    main();
}
