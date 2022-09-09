package utils

import "fmt"

func check(err error) {
	if err != nil {
		fmt.Println(err)
		return
	}
}
