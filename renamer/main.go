package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

func main() {
	// 	fileName := "birthday_001.txt"
	// 	// => Birthday - 1 of 4.txt
	// 	newName, err := match(fileName, 4)
	// 	if err != nil {
	// 		fmt.Println("no match")
	// 		os.Exit(1)
	// 	}
	// 	fmt.Println(newName)
	files, err := ioutil.ReadDir("./sample")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			fmt.Println("Dir: ", file.Name())
		} else {
			tmp, err := match(file.Name(), 4)
			fmt.Println("match: ", tmp, err)
		}
	}
}

func match(fileName string, total int) (string, error) {
	pieces := strings.Split(fileName, ".")
	ext := pieces[len(pieces)-1]
	tmp := strings.Join(pieces[0:len(pieces)-1], ".")
	pieces = strings.Split(tmp, "_")
	name := strings.Join(pieces[0:len(pieces)-1], "_")
	number, err := strconv.Atoi(pieces[len(pieces)-1])
	if err != nil {
		return "", fmt.Errorf("%s didn't match our patter", tmp)
	}
	// Birthday - 1 .txt
	return fmt.Sprintf("%s - %d of %d.%s", strings.Title(name), number, total, ext), nil

	return "", nil
}
