# go-primes
Naïve functions to check primality and find prime numbers.
Useful to generate data to study concurrency with for
CPU-intensive computations.


## Idea

The time to determine primality depends on the size of
the smallest prime factor. For composite numbers,
the smallest prime factor is the value that makes `IsPrime`
return.

We should generate a list of pairs `(number, smallestPrimeFactor)`.
If `number` is prime, `smallestPrimeFactor == number`.
