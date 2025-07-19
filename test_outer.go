package main

import "fmt"

func main() {
outerLoop:
	for i := 0; i < 5; i++ {
		fmt.Println("Outer loop i:", i)

		for j := 0; j < 5; j++ {
			fmt.Println("  Inner loop j:", j)

			if i == 2 && j == 3 {
				fmt.Println("  Breaking out of both loops!")
				break outerLoop // 🔥 This breaks out of BOTH loops
			}
		}
	}

	fmt.Println("Finished.")
}
