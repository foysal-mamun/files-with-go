package files

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

func CreateEmptly(fileName string) {

	if fileName == "" {
		return
	}

	newFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(newFile)
	newFile.Close()
}

func Truncate(fileName string, size int64) {

	if fileName == "" {
		return
	}

	err := os.Truncate(fileName, size)
	if err != nil {
		log.Fatal(err)
	}
}

func GetInfo(fileName string) {

	if fileName == "" {
		return
	}

	fileInfo, err := os.Stat(fileName)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("File name:", fileInfo.Name())
	fmt.Println("Size in bytes:", fileInfo.Size())
	fmt.Println("Permissions:", fileInfo.Mode())
	fmt.Println("Last modified:", fileInfo.ModTime())
	fmt.Println("Is Directory: ", fileInfo.IsDir())
	fmt.Printf("System interface type: %T\n", fileInfo.Sys())
	fmt.Printf("System info: %+v\n\n", fileInfo.Sys())

}

func Move(oldLoc string, newLoc string) {

	if oldLoc == "" || newLoc == "" {
		return
	}

	err := os.Rename(oldLoc, newLoc)
	if err != nil {
		log.Fatal(err)
	}
}

func Delete(fileName string) {

	if fileName == "" {
		return
	}
	err := os.Remove(fileName)
	if err != nil {
		log.Fatal(err)
	}

}

func Open(fileName string) {

	if fileName == "" {
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(os.IsNotExist(err))
	}
	file.Close()

	file, err = os.OpenFile(fileName, os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

func CheckPermission(fileName string) {

	if fileName == "" {
		return
	}

	file, err := os.OpenFile(fileName, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(os.IsPermission(err))
	}
	file.Close()
}

func ChangePermission(fileName string, mode int) {

	if fileName == "" {
		return
	}

	err := os.Chmod(fileName, os.FileMode(mode))
	if err != nil {
		log.Fatal(err)
	}
}

func ChangeOwnership(fileName string, uid, gid int) {

	if uid == 0 {
		uid = os.Getuid()
	}
	if gid == 0 {
		gid = os.Getgid()
	}

	err := os.Chown(fileName, uid, gid)
	if err != nil {
		log.Fatal(err)
	}
}

func ChangeTime(fileName string, atime, mtime time.Time) {
	if err := os.Chtimes(fileName, atime, mtime); err != nil {
		log.Fatal(err)
	}
}

func HardLink(oldname, newname string) {
	if oldname == "" || newname == "" {
		return
	}

	if err := os.Link(oldname, newname); err != nil {
		log.Fatal(err)
	}
}

func SymLink(oldname, newname string) {
	if err := os.Symlink(oldname, newname); err != nil {
		log.Fatal(err)
	}
}

func Copy(oldname, newname string) (err error) {

	oldfile, err := os.Stat(oldname)
	if err != nil {
		return err
	}
	if !oldfile.Mode().IsRegular() {
		return fmt.Errorf("Copy: non-regular old file %s (%q)", oldfile.Name(), oldfile.Mode().String())
	}

	newfile, err := os.Stat(newname)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(newfile.Mode().IsRegular()) {
			return fmt.Errorf("Copy: non-regular new file %s (%q)", oldfile.Name(), oldfile.Mode().String())
		}
		if os.SameFile(oldfile, newfile) {
			return
		}
	}

	if err = os.Link(oldname, newname); err != nil {
		return
	}

	err = copyFileContents(oldname, newname)

	return
}

func copyFileContents(oldname, newname string) (err error) {

	oldfile, err := os.Open(oldname)
	if err != nil {
		return
	}
	defer oldfile.Close()

	newfile, err := os.Create(newname)
	if err != nil {
		return
	}
	defer func() {
		nerr := newfile.Close()
		if err == nil {
			err = nerr
		}
	}()

	if _, err = io.Copy(newfile, oldfile); err != nil {
		return
	}
	err = newfile.Sync()

	return
}

// offset
// how many bytes to move, can be positive or negative
// whence
// 0 = Beginning of file
// 1 = Current position
// 2 = End of file
func Seek(fileName string, offset int64, whence int) int64 {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	newPos, err := file.Seek(offset, whence)
	if err != nil {
		log.Fatal(err)
	}

	return newPos
}
