#include <bits/stdc++.h>
using namespace std;

struct Point {
    int x, y;
};

// Check if point lies on the curve
bool isPoint(Point a1, int a, int b, int p) {
    int y = (1LL * a1.y * a1.y) % p;
    int x = ((1LL * a1.x * a1.x % p * a1.x % p) + (1LL * a * a1.x % p) + b) % p;
    return x == y;
}

// Print all points on the curve
void allPoints(int p, int a, int b) {
    for (int i = 0; i < p; i++) {
        for (int j = 0; j < p; j++) {
            Point a2 = {i, j};
            if (isPoint(a2, a, b, p)) {
                cout << i << " " << j << "\n";
            }
        }
    }
}

int main() {
    int a, b, p;
    cout << "Enter a: ";
    cin >> a;
    cout << "Enter b: ";
    cin >> b;
    cout << "Enter p: ";
    cin >> p;

    cout << "All Points on the curve are:\n";
    allPoints(p, a, b);

    return 0;
}
