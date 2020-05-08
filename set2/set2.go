package main

import "fmt"

type challengeFunc func() error

func main() {
	challenges := []challengeFunc{
		runC9,
		runC10,
		runC11,
		runC12,
		runC13,
		runC14,
		runC16,
	}

	for n, challengeFunc := range challenges {
		fmt.Println("Challenge", n+9)
		if err := challengeFunc(); err != nil {
			fmt.Println(err)
		}
		fmt.Println("")
	}
}
