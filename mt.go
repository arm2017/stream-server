package main

import (
	"fmt"
	moveregister "stream-cam-api/move-register"
)

func main() {
	moveregister.RunTest()
	fmt.Print("OK.")
}
