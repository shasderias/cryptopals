package main

import "fmt"

type challengeFunc func() error

func main() {
	challenges := []challengeFunc{
		runC17,
	}

	for n, challengeFunc := range challenges {
		fmt.Println("Challenge", n+17)
		if err := challengeFunc(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("")
	}
}
