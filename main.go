package main

import (
	"fmt"
	"github.com/foysal-mamun/files-with-go/files"
	//"time"
)

func main() {
	fmt.Println("files")

	//files.CreateEmptly("test.txt")
	//files.Truncate("test.txt", 100)

	//files.GetInfo("test.txt")

	//files.Move("test.txt", "text.txt")
	//files.GetInfo("text.txt")

	//files.Delete("text.txt")

	//files.Open("test.txt1")

	//files.CheckPermission("test.txt")
	//files.ChangePermission("test.txt", 0777)
	//files.ChangeOwnership("test.txt", 0, 0)
	//files.ChangeTime("test.txt", time.Now(), time.Now())

	/*
		err := files.Copy("test.txt", "text.txt")
		if err != nil {
			fmt.Println("Copy File failed %q", err)
		} else {
			fmt.Println("Copy file succeeded")
		}
	*/

	files.Seek("test.txt1", 5, 0)

}
