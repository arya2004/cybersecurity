import java.util.Scanner;

class Point {
    int x, y;

    Point(int x, int y) {
        this.x = x;
        this.y = y;
    }
}

public class Ecc_1 {

    // Check if a point lies on the elliptic curve y² = x³ + ax + b (mod p)
    static boolean isPoint(Point a1, int a, int b, int p) {
        long ySq = (long) a1.y * a1.y % p;
        long xCubed = (long) a1.x * a1.x * a1.x % p;
        long ax = (long) a * a1.x % p;
        long rhs = (xCubed + ax + b) % p;

        return (int) rhs == (int) ySq;
    }

    // Print all points on the curve
    static void allPoints(int p, int a, int b) {
        for (int i = 0; i < p; i++) {
            for (int j = 0; j < p; j++) {
                Point a2 = new Point(i, j);
                if (isPoint(a2, a, b, p)) {
                    System.out.println(i + " " + j);
                }
            }
        }
    }

    public static void main(String[] args) {
        Scanner sc = new Scanner(System.in);
        System.out.print("Enter a: ");
        int a = sc.nextInt();
        System.out.print("Enter b: ");
        int b = sc.nextInt();
        System.out.print("Enter p: ");
        int p = sc.nextInt();

        System.out.println("All Points on the curve are:");
        allPoints(p, a, b);
        sc.close();
    }
}
