package main

import (
	"fmt"
	"github.com/paterson/binder/utils/clientproxy"
)

func main() {
	printSteps()
	clientproxy := clientproxy.New()
	clientproxy.Signup("niall", "password")
	clientproxy.LockFile("/test.png")
	clientproxy.WriteFile("./test.png", "/test.png")
	clientproxy.UnlockFile("/test.png")
	clientproxy.ReadFile("/test.png", "./test1.png")
}

func printSteps() {
	fmt.Println("")
	fmt.Println("*******************************************************************")
	fmt.Println("This default client implementation does a number of things:")
	fmt.Println("* Signups a user")
	fmt.Println("* Locks the file ./test.png")
	fmt.Println("* Writes the file ./test.png. This will create write it to both running file servers (./.files and ./fileserver/.files)")
	fmt.Println("* Unlocks the file ./test.png")
	fmt.Println("* Reads the file back down to test1.png and stores it in cache.")
	fmt.Println("*******************************************************************")
	fmt.Println("")
}
