const p = 3;
const q = 7;

let e = 2;
const k = 2;
const msg = 12;

function gcd(a, h) {
  let temp;
  while (true) {
    temp = a % h;
    if (temp === 0) {
      return h;
    }
    a = h;
    h = temp;
  }
}

function main() {
  const n = p * q;
  console.log("Calculated n (p * q) = ", n);

  const phi = (p - 1) * (q - 1);
  console.log("Calculated phi ((p - 1) * (q - 1)) = ", phi);

  if (gcd(e, phi) !== 1) {
    console.log("e =", e, "is not coprime with phi. Incrementing e.");
  }

  while (e < phi) {
    if (gcd(e, phi) === 1) {
      console.log("Chosen e = ", e, " as it is coprime with phi");
      break;
    } else {
      e++;
    }
  }

  const d = Math.floor((1 + k * phi) / e);
  console.log("Calculated d (private key) = ", d);

  console.log("Message data = ", msg);

  const c = Math.pow(msg, e) % n;
  console.log("Encrypted data = ", c);

  const m = Math.pow(c, d) % n;
  console.log("Original Message Sent = ", m);
}

main();
