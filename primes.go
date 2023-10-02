package main

import (
	"errors"
	"fmt"
	"os"
	"math"
	"strconv"
	"time"
)

const MaxUint64 = ^uint64(0) // using two's complement
const MaxUint64Prime = 18446744073709551557

// IsPrime returns true if n is prime, false otherwise.
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

// NextPrime finds the next prime number starting at n.
// If n is prime, it returns n.
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


// PreviousPrime finds the previous prime number starting at n.
// If n is prime, it returns n.
func PreviousPrime(n uint64) (uint64, error) {
	if IsPrime(n) {
		return n, nil
	}
	start := n
	if start%2 == 0 {
		start-- // skip even
	}
	for i := start; i >= 2; i -= 2 {
		if IsPrime(i) {
			return i, nil
		}
	}
	return 0, errors.New("no primes < 2")
}


// SemiprimeNear finds a semiprime close to the target.
func SemiprimeNear(target uint64) (uint64, error) {
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

// TODO: this does not work as intended.
// It should find a semiprime that takes at least
// dt seconds to check with IsPrime.
// Instead, it counts the time for the entire
// search, not just the primality check.
func findPrime(dt float64) (uint64, float64) {
	var prime uint64
	var err error
	var elapsed float64
	for bits := 32; bits < 64; bits++ {
		start := time.Now()
		prime, err = NextPrime(2<<bits - 1)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		elapsed = time.Since(start).Seconds()
		if elapsed >= dt {
			break
		}
	}
	return prime, elapsed
}

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:\t%v <time-in-seconds>\n", os.Args[0])
		fmt.Println("\tfind a prime that takes at least <time-in-seconds> to find")
		os.Exit(1)
	}
	dt, err := strconv.ParseFloat(os.Args[1], 32)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	prime, elapsed := findPrime(dt)
	fmt.Printf("Found prime %v in %v seconds\n", prime, elapsed)
}	
