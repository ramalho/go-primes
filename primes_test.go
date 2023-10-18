package main

import (
	"errors"
	"flag"
	"fmt"
	"testing"
)

var sample = []struct {
	n       uint64
	isPrime bool
}{
	{2, true},
	{142702110479723, true},
	{299593572317531, true},
	{3333333333333301, true},
	{3333333333333333, false},
	{3333335652092209, false},
	{4444444444444423, true},
	{4444444444444444, false},
	{4444444488888889, false},
	{5555553133149889, false},
	{5555555555555503, true},
	{5555555555555555, false},
	{6666666666666666, false},
	{6666666666666719, true},
	{6666667141414921, false},
	{7777777536340681, false},
	{7777777777777753, true},
	{7777777777777777, false},
	{9999999999999917, true},
	{9999999999999999, false},
	{MaxUint64, false},
}

func TestIsPrime(t *testing.T) {
	t.Parallel()
	for _, testCase := range sample {
		tc := testCase
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			t.Parallel()
			got := IsPrime(tc.n)
			if got != tc.isPrime {
				t.Errorf("%v got: %v",
					tc, got)
			}
		})
	}
}

// First primes: 2, 3, 5, 7, 11, 13, 17, 19, 23
// Source: https://oeis.org/A000040

var nextPrimes = []struct {
	n    uint64
	next uint64
	err  error
}{
	{0, 2, nil},
	{1, 2, nil},
	{2, 2, nil},
	{3, 3, nil},
	{4, 5, nil},
	{20, 23, nil},
	{7777, 7789, nil},
	{6666666666666666, 6666666666666719, nil},
	{MaxUint64Prime + uint64(1), 0,
		errors.New(
			fmt.Sprintf("no primes >= %v in uint64 range", MaxUint64Prime+uint64(1)))},
}

func TestNextPrime(t *testing.T) {
	t.Parallel()
	for _, testCase := range nextPrimes {
		tc := testCase
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			t.Parallel()
			got, err := NextPrime(tc.n)
			if got != tc.next {
				t.Errorf("%v expected: %v, got: %v",
					tc.n, tc.next, got)
			}
			if got == 0 && tc.err.Error() != err.Error() {
				t.Errorf("%v expected error %q, got: %q",
					tc.n, tc.err, err)
			}
		})
	}
}

// First primes: 2, 3, 5, 7, 11, 13, 17, 19, 23

type prevPrime struct {
	n    uint64
	prev uint64
	err  error
}

var prevPrimes = []prevPrime{
	{2, 2, nil},
	{3, 3, nil},
	{4, 3, nil},
	{5, 5, nil},
	{6, 5, nil},
	{22, 19, nil},
	{23, 23, nil},
	{7777, 7759, nil},
	{6666666666666720, 6666666666666719, nil},
	{1, 0, errors.New("no primes < 2")},
}

// run `go test -slow` to set this flag to true
var slow = flag.Bool("slow", false, "perform slow tests")

func TestPreviousPrime(t *testing.T) {
	t.Parallel()
	if *slow {
		var slowCase = prevPrime{MaxUint64, MaxUint64Prime, nil}
		prevPrimes = append(prevPrimes, slowCase)
	}
	for _, testCase := range prevPrimes {
		tc := testCase
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			t.Parallel()
			got, err := PreviousPrime(tc.n)
			if got != tc.prev {
				t.Errorf("before %v expected: %v, got: %v",
					tc.n, tc.prev, got)
			}
			if got == 0 && tc.err.Error() != err.Error() {
				t.Errorf("%v expected error %q, got: %q",
					tc.n, tc.err, err)
			}
		})
	}
}

// Semiprimes:
// 4, 6, 9, 10, 14, 15, 21, 22, 25, 26, 33...

var nextSemiprimes = []struct {
	arg uint64
	pr  PrimeResult
}{
	{4, PrimeResult{4, 2}},
	{5, PrimeResult{6, 2}},
}

func TestNextSemiprime(t *testing.T) {
	t.Parallel()
	for _, testCase := range nextSemiprimes {
		tc := testCase
		t.Run(fmt.Sprintf("%v", testCase), func(t *testing.T) {
			t.Parallel()
			got := NextSemiprime(tc.arg)
			if got != tc.pr {
				t.Errorf("expected: %v, got: %v",
					tc, got)
			}
		})
	}
}

var semiprimesNear = []struct {
	arg uint64
	pr  PrimeResult
}{
	{1, PrimeResult{4, 2}},
	{2, PrimeResult{4, 2}},
	{3, PrimeResult{4, 2}},
	{4, PrimeResult{4, 2}},
	{5, PrimeResult{4, 2}},
	{6, PrimeResult{6, 2}},
	{7, PrimeResult{6, 2}},
	{8, PrimeResult{9, 3}},
	{9, PrimeResult{9, 3}},
	{10, PrimeResult{10, 2}},
	{11, PrimeResult{10, 2}},
	{12, PrimeResult{10, 2}},
	{13, PrimeResult{14, 2}},
	{14, PrimeResult{14, 3}},
	{15, PrimeResult{15, 3}},
	// {16, PrimeResult{15, 3}},
	// {17, PrimeResult{15, 3}},
	// {18, PrimeResult{21, 3}},
	// {20, PrimeResult{21, 3}},
	// {21, PrimeResult{25, 5}},
	// {30, PrimeResult{25, 5}},
	// {31, PrimeResult{35, 5}},
	// {39, PrimeResult{35, 5}},
	// {40, PrimeResult{55, 5}},
	// {40, PrimeResult{55, 5}},
	// {10000000, 9997619},
	// {100000000, 100099001},
	// {10000000000000000000, 10000000011584185403},
	// {18000000000000000000, 18000000066870082769},
	// {18446744000000000000, 18446744065119616769},
	// {18446744070000000000, 18446744030759878681},
	// {MaxUint64Prime, 18446744030759878681},
	//{MaxUint64, 18446744030759878681},
}

// func TestSemiprimeNearest(t *testing.T) {
// 	t.Parallel()
// 	for _, testCase := range semiprimesNear {
// 		tc := testCase
// 		t.Run(fmt.Sprintf("%v", testCase), func(t *testing.T) {
// 			t.Parallel()
// 			got := SemiprimeNearest(tc.arg)
// 			if got != tc.pr {
// 				t.Errorf("expected: %v, got: %v",
// 					tc, got)
// 			}
// 		})
// 	}
// }
