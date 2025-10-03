class Point {
  constructor(x, y) {
    this.x = x;
    this.y = y;
  }
}

// Check if point lies on the curve
function isPoint(a1, a, b, p) {
  let y = (a1.y * a1.y) % p;
  let x = ((a1.x * a1.x * a1.x) + (a * a1.x) + b) % p;
  return x === y;
}

// Print all points on the curve
function allPoints(p, a, b) {
  for (let i = 0; i < p; i++) {
    for (let j = 0; j < p; j++) {
      let a2 = new Point(i, j);
      if (isPoint(a2, a, b, p)) {
        console.log(i, j);
      }
    }
  }
}

const readline = require("readline").createInterface({
  input: process.stdin,
  output: process.stdout
});

readline.question("Enter a: ", (a) => {
  readline.question("Enter b: ", (b) => {
    readline.question("Enter p: ", (p) => {
      a = parseInt(a);
      b = parseInt(b);
      p = parseInt(p);

      console.log("All Points on the curve are:");
      allPoints(p, a, b);

      readline.close();
    });
  });
});
