package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	s := "2" + 1
	// fmt.Println(s)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		fmt.Println(r.Intn(10))
	}
}
