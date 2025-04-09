package main

import (
	"fmt"
)

// 判断一个数是否为质数
func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	if n == 2 || n == 3 {
		return true
	}
	if n%2 == 0 || n%3 == 0 {
		return false
	}
	for i := 5; i*i <= n; i += 6 {
		if n%i == 0 || n%(i+2) == 0 {
			return false
		}
	}
	return true
}

// 判断一个数是否为漂亮数
func isBeautiful(x int) bool {
	// 特殊情况：1不是漂亮数
	if x == 1 {
		return false
	}

	// 如果x本身是质数，那么它是漂亮数（因为它自己就是满足条件的质因子）
	if isPrime(x) {
		return true
	}

	// 检查所有可能的质因子
	for p := 2; p*p <= x; p++ {
		if isPrime(p) && x%p == 0 && p*p >= x {
			return true
		}
	}

	// 检查大于sqrt(x)的质因子
	for p := x / 2; p >= 2; p-- {
		if x%p == 0 && isPrime(p) && p*p >= x {
			return true
		}
	}

	return false
}

func main() {
	var n int
	fmt.Scan(&n)

	count := 0
	for i := 1; i <= n; i++ {
		if isBeautiful(i) {
			count++
		}
	}

	fmt.Println(count)
}
