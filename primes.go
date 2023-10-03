package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

const MaxUint64 = ^uint64(0) // using two's complement
const MaxUint64Prime = 18446744073709551557

// LPF returns the Least Prime Factor of `n“.
// If `n` is prime, `LPF(n)` returns `n`.
func LPF(n uint64) uint64 {
	switch {
	case n == 1:
		return 1
	case n%2 == 0:
		return 2
	case n%3 == 0:
		return 3
	}
	limit := uint64(math.Sqrt(float64(n)))
	for i := uint64(5); i <= limit; i += 6 {
		if n%i == 0 {
			return i
		}
		j := i + 2
		if n%j == 0 {
			return j
		}
	}
	return n
}

// IsPrime returns true if n is prime, false otherwise.
func IsPrime(n uint64) bool {
	if n <= 1 {
		return false
	}
	return LPF(n) == n
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
		a, _ = NextPrime(root)
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

const shortTime = 0.0001

func demo() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage:\t%v n\n", os.Args[0])
		fmt.Println("\tfind primes closest to n")
		os.Exit(1)
	}
	n, err := strconv.ParseUint(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// previous prime
	start := time.Now()
	prev, err := PreviousPrime(n)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	elapsed := time.Since(start).Seconds()
	var msg string
	if elapsed < shortTime {
		msg = fmt.Sprintf("< %vs", shortTime)
	} else {
		msg = fmt.Sprintf("%0.3fs", elapsed)
	}
	if prev == n {
		fmt.Printf("%v is prime (%v)\n", n, msg)
		os.Exit(0)
	}
	fmt.Printf("%20d  # previous prime (%v)\n", prev, msg)
	// next prime
	start = time.Now()
	next, err := NextPrime(n)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	elapsed = time.Since(start).Seconds()
	if elapsed < shortTime {
		msg = fmt.Sprintf("< %vs", shortTime)
	} else {
		msg = fmt.Sprintf("%0.3fs", elapsed)
	}
	fmt.Printf("%20d  # next prime (%v)\n", next, msg)
}

var interestingIntegers = []uint64{
	2,
	3,
	17592186044416,       // 2 ** 44
	281474976710656,      // 2 ** 48
	4503599627370496,     // 2 ** 52
	72057594037927936,    // 2 ** 56
	1152921504606846976,  // 2 ** 60
	18446744073709551615, // 2 ** 64 - 1

}

func interestingReport() {
	lineNum := 1
	for _, n := range interestingIntegers {
		pp, _ := PreviousPrime(n)
		if pp != n {
			reportLine(lineNum, pp)
			lineNum++
		}
		reportLine(lineNum, n)
		lineNum++
		sp, _ := SemiprimeNear(n)
		reportLine(lineNum, sp)
		lineNum++
	}
}

func reportLine(i int, n uint64) {
	start := time.Now()
	lpf := LPF(n)
	var label string
	if n == lpf {
		label = "P"
	}
	elapsed := time.Since(start).Seconds()
	fmt.Printf("%2d  %20d  %1s %20d  // %0.6fs\n", i, n, label, lpf, elapsed)

}

func main() {
	interestingReport()
}
