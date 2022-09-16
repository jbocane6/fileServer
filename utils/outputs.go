package utils

import "fmt"

func check(err error) {
	// Checks if getting any error and shows error content
	if err != nil {
		fmt.Println(err)
		return
	}
}

func compareSizeWitten(fileSize, bytesWritten int64) {
	// Checks if there's a difference between file size and bytes written
	if fileSize != bytesWritten {
		fmt.Println(ErrFileSize)
	}
}

func compare(bytes, stat, length int64) {
	// Validates if bytes is the same as the content size
	if bytes != stat {
		fmt.Printf(ErrFileNameSize, bytes, length)
	}
}
