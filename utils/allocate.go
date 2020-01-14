package main

import (
	"fmt"
	"github.com/c2h5oh/datasize"
	"math/rand"
	"os"
	"strconv"
	"time"

)

func check(err error){
	if err != nil{
		panic(err)
	}
}

func main() {

	mb, err := strconv.ParseInt(os.Args[1], 10, 32)
	check(err)

	chunk := 20
	var arrs [][]byte

	inc := datasize.ByteSize(mb) * datasize.MB / datasize.ByteSize(chunk)

	for i := 0; i < chunk; i++{

		arrs = append(arrs, make([]byte,  inc.Bytes() ))

		for j := range arrs[i]{
			arrs[i][j] = byte(rand.Int())
		}
		fmt.Println("Allocated", (datasize.ByteSize(i+1)*inc).String())

	}

	fmt.Println("Done, cleaning up in 30 sec")
	time.Sleep(30*time.Second)
}
