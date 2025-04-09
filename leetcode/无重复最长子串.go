package main

import (
	"bufio"
	"fmt"
	"os"
)

func lengthOfLongestSubstring(s string) int {
	cnt := [128]int{}
	left := 0
	maxLen := 0
	for right, c := range s {
		cnt[c]++
		for cnt[c] > 1 {
			cnt[s[left]]--
			left++
		}
		maxLen = max(maxLen, right-left+1)
	}
	return maxLen
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	s := scanner.Text()
	result := lengthOfLongestSubstring(s)
	fmt.Println(result)
}
