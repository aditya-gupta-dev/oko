package main

import (
	"fmt"
	"os"
	"strings"
)

// func main() {
// 	path := "C:/Users/hyper/progs/ytt/output"

// 	songs, err := song.ListSongFiles(path)

// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}

// 	for _, song := range songs {
// 		fmt.Println(song.Name, song.Duration)
// 	}
// }

func main() {
	file, err := os.Open("C:/Users/hyper/code/golang/th/main.go")

	if err != nil {
		fmt.Println("Error", err.Error())
	}

	iter := strings.Split(file.Name(), "/")
	fmt.Println(iter[len(iter)-1])
}
