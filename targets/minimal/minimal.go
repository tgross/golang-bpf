package main

import (
	"bufio"
	"fmt"
	"os"
)

// function that has multiple returns and calls other functions, makes
// syscalls, and includes a defer
func read(p string) (string, error) {

	f, err := os.Open(p)
	defer f.Close() // <-- a defer! how's that work?
	if err != nil {
		return "", err
	}

	r := bufio.NewReader(f)
	buf, err := r.Peek(5)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("5 bytes: [% x]\n", string(buf)), nil
}

// easy function
func add(x int, y int) int {
	return x + y
}

// function that returns strings
func swap(x, y string) (string, string) {
	return y, x
}

func main() {
	fmt.Println(add(42, 13))

	a, b := swap("hello", "world")
	fmt.Println(a, b)

	val, err := read("/etc/hosts")
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}
