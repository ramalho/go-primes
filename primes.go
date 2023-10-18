package main

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"slices"
	"strconv"
	"time"
)

const MaxUint64 = ^uint64(0) // using two's complement
const MaxUint64Prime = 18446744073709551557

// PrimeResult represents a number and its least prime factor
type PrimeResult struct {
	N   uint64
	LPF uint64
}

func (p PrimeResult) IsPrime() bool {
	return p.N >= 1 && p.N == p.LPF
}

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

	// n.ProbablyPrime(0) uses the Baillie-PSW primality test
	// https://en.wikipedia.org/wiki/Baillie%E2%80%93PSW_primality_test
	if n < math.MaxInt64 && big.NewInt(int64(n)).ProbablyPrime(0) {
		return n
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

// IsPrime returns true if `n“ is prime, false otherwise.
func IsPrime(n uint64) bool {
	if n <= 1 {
		return false
	}
	return LPF(n) == n
}

// NextPrime finds the next prime number starting at `n`.
// If `n` is prime, it returns `n“.
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
	return 0, fmt.Errorf("no primes >= %v in uint64 range", n)
}

// PreviousPrime finds the previous prime number starting at `n“.
// If `n“ is prime, it returns `n“.
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

// NextSemiprime finds the next semi prime after `n`.
// If `n` is semi prime, it returns `n`.
func NextSemiprime(n uint64) PrimeResult {
	a := LPF(n)
	b := n / a
	if b < 2 {
		b = 2
	}
	prod := a * b
	if prod == n && IsPrime(a) {
		return PrimeResult{prod, a}
	}
	for prod < n {
		b, err := NextPrime(b + 1)
		if err != nil {
			panic(err)
		}
		prod = a * b
	}
	return PrimeResult{prod, a}

}

const shortTime = 0.0001

func FindPrimes() {
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

func UintPow(n, exponent uint64) uint64 {
	if exponent == 0 {
		return 1
	}
	result := n
	for i := uint64(2); i <= exponent; i++ {
		result *= n
	}
	return result
}

func genTargetsLinear(queue chan<- uint64) {
	const step = MaxUint64 / 16
	var previous uint64 = 12
	queue <- previous
	for n := step; n <= MaxUint64; n += step {
		if n < previous {
			break // handle uint64 overflow
		}
		queue <- n
		previous = n
	}
	close(queue)
}

func genTargetsExp2(queue chan<- uint64) {
	for e := uint64(30); e < 64; e++ {
		queue <- UintPow(2, e)
	}
	queue <- 18446744073709551615 // 2 ** 64 - 1
	close(queue)
}

func isPowerOf2(n uint64) (bool, uint64) {
	for i := uint64(2); i < 64; i++ {
		if n == UintPow(2, i) {
			return true, i
		}
	}
	return false, 0
}

type reportLine = struct {
	n       uint64
	comment string
}

func report() {
	go isPowerOf2(3)
	queue := make(chan uint64)
	go genTargetsExp2(queue)
	var lines []reportLine
	for n := range queue {
		pp, _ := PreviousPrime(n)
		if pp != n {
			lines = append(lines, reportLine{pp, "prime"})
		}
		var comment string

		if isP2, p2 := isPowerOf2(n); isP2 {
			comment = fmt.Sprintf("2 ** %v", p2)
		}

		lines = append(lines, reportLine{n, comment})
		sp := NextSemiprime(n)
		if sp.N != n {
			lines = append(lines, reportLine{sp.N, "semiprime"})
		}
		slices.SortFunc(lines, func(a, b reportLine) int {
			return cmp.Compare(a.n, b.n)
		})
		for _, line := range lines {
			displayLine(line)
		}
		lines = nil
	}
	displayLine(reportLine{MaxUint64, "2 ** 64 - 1"})
}

func displayLine(l reportLine) {
	// format line to use as fixture for primes.py in python-eng repo
	var lpf uint64
	if l.comment == "prime" {
		lpf = l.n
	} else {
		lpf = LPF(l.n)
	}

	// Experiment(17592186044416, 2),  # 2 ** 44
	if len(l.comment) > 0 {
		l.comment = "  # " + l.comment
	}
	fmt.Printf("Experiment(%20d, %20d),%v\n", l.n, lpf, l.comment)

}

func main() {
	// report()
	FindPrimes()
}
