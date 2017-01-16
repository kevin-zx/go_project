package main

import (
	"fmt"
	"bufio"
	// "io"
	"os"
	"strings"
	
)

func main() {
	fileName:="data/config.txt"
	f, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	arra := [7]string{}
	buf := bufio.NewReader(f)
	// defer buf.Close()
	i := 0
	for {

		line, err := buf.ReadString('\n')
		
		if err != nil {
			break
		}
		line = strings.TrimSpace(line) 
		arra[i] = line
		i=i+1
	}

}
