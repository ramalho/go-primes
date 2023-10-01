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

var semiprimesNear = []struct {
	target    uint64
	semiprime uint64
	err       error
}{
	{1, 4, nil},
	{2, 4, nil},
	{3, 4, nil},
	{5, 4, nil},
	{6, 4, nil},
	{7, 9, nil},
	{8, 9, nil},
	{10, 9, nil},
	{12, 9, nil},
	{13, 15, nil},
	{14, 15, nil},
	{15, 15, nil},
	{16, 15, nil},
	{17, 15, nil},
	{18, 21, nil},
	{20, 21, nil},
	{21, 25, nil},
	{30, 25, nil},
	{31, 35, nil},
	{100, 119, nil},
	{121, 121, nil},
	{130, 121, nil},
	{900, 899, nil},
	{1000, 1147, nil},
	{10000000, 9997619, nil},
	{100000000, 100099001, nil},
	{10000000000000000000, 10000000011584185403, nil},
	{18000000000000000000, 18000000066870082769, nil},
	{18446744000000000000, 18446744065119616769, nil},
	{18446744070000000000, 18446744030759878681, nil},
	{MaxUint64Prime, 18446744030759878681, nil},
	{MaxUint64, 18446744030759878681, nil},
}

func TestSemiprimeNear(t *testing.T) {
	t.Parallel()
	for _, testCase := range semiprimesNear {
		tc := testCase
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			t.Parallel()
			got, err := SemiprimeNear(tc.target)
			if got != tc.semiprime {
				t.Errorf("near %v expected: %v, got: %v",
					tc.target, tc.semiprime, got)
			}
			if got == 0 && tc.err.Error() != err.Error() {
				t.Errorf("%v expected error %q, got: %q",
					tc.target, tc.err, err)
			}
		})
	}
}
