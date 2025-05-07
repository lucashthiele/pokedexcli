package logs

import "fmt"

func Log(values ...any) {
	fmt.Printf("-------- Log ---------\n")
	for _, value := range values {
		fmt.Printf("%q\n", value)
	}
	fmt.Printf("----------------------\n")
}
