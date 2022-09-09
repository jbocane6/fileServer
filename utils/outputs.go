package utils

import "fmt"

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}

func compare(bytes, stat, length int64) {
	if bytes != stat {
		fmt.Printf("Error! Wrote %d bytes but length of name is %d!\n", bytes, length)
	}
}
