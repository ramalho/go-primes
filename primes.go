package main

import (
	"errors"
	"fmt"
	"math"
	"time"
)

const MaxUint64 = ^uint64(0) // using two's complement
const MaxUint64Prime = 18446744073709551557

func IsPrime(n uint64) bool {
	if n == 2 || n == 3 {
		return true
	}
	if n <= 1 || n%2 == 0 || n%3 == 0 {
		return false
	}

	limit := uint64(math.Sqrt(float64(n)))
	for i := uint64(5); i <= limit; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

func NextPrime(n uint64) (uint64, error) {
	if n < 2 {
		return 2, nil
	}
	if IsPrime(n) {
		return n, nil
	}
	start := n
	if start%2 == 0 {
		start++ // skip even
	}
	for i := start; i < MaxUint64; i += 2 {
		if IsPrime(i) {
			return i, nil
		}
	}
	return 0, errors.New(
		fmt.Sprintf("no primes >= %v in uint64 range", n))
}

func PreviousPrime(n uint64) (uint64, error) {
	if IsPrime(n) {
		return n, nil
	}
	if n%2 == 0 {
		n-- // skip even
	}
	for i := n; i > 2; i -= 2 {
		if IsPrime(i) {
			return i, nil
		}
	}
	return 0, errors.New("no primes < 2")
}

func TwoFactorComposite(target uint64) (uint64, error) {
	root := uint64(math.Round(math.Sqrt(float64(target))))
	root = max(root, 2) // 2 is the smallest prime
	if IsPrime(root) {
		return root * root, nil
	}
	a, err := PreviousPrime(root)
	if err != nil {
		a, err = NextPrime(root)
	}
	b, err := NextPrime(target / a)
	if err != nil {
		return 0, err
	}
	if a*b < a*a { // handle uint64 overflow
		return a * a, nil
	}
	return a * b, nil
}

func main() {

	fmt.Println("The largest prime in the uint64 range is...")
	start := time.Now()
	bigPrime, _ := PreviousPrime(MaxUint64)
	elapsed := time.Since(start)
	fmt.Println(bigPrime)
	fmt.Printf("Which is %v less than the largest uint64: %v\n",
		MaxUint64-bigPrime, MaxUint64)
	fmt.Printf("(computed in %.2fs)\n", elapsed.Seconds())

}
