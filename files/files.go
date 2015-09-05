package files

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func CreateEmptly(fileName string) {

	if fileName == "" {
		return
	}

	newFile, err := os.Create(fileName)
	checkError(err)

	log.Println(newFile)
	newFile.Close()
}

func Truncate(fileName string, size int64) {

	if fileName == "" {
		return
	}

	err := os.Truncate(fileName, size)
	checkError(err)
}

func GetInfo(fileName string) {

	if fileName == "" {
		return
	}

	fileInfo, err := os.Stat(fileName)
	checkError(err)

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
	checkError(err)
}

func Delete(fileName string) {

	if fileName == "" {
		return
	}
	err := os.Remove(fileName)
	checkError(err)

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
	checkError(err)
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
	checkError(err)
}

func ChangeOwnership(fileName string, uid, gid int) {

	if uid == 0 {
		uid = os.Getuid()
	}
	if gid == 0 {
		gid = os.Getgid()
	}

	err := os.Chown(fileName, uid, gid)
	checkError(err)
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
	checkError(err)
	defer file.Close()

	newPos, err := file.Seek(offset, whence)
	checkError(err)

	return newPos
}

func Write(fileName string, content []byte) {

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	checkError(err)
	defer file.Close()

	_, err = file.Write(content)
	checkError(err)
}

func Read(fileName string, len int) {

	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()

	byteSlice := make([]byte, len)
	bytesRead, err := file.Read(byteSlice)
	checkError(err)

	log.Println(byteSlice[:len])
	log.Println(bytesRead)
}

/**
 * Archive given files
 * @param {string} zipName   string
 * @param {array of string} fileNames []string
 */
func CreateArchiveFile(zipName string, fileNames []string) {

	file, err := os.OpenFile(zipName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	checkError(err)
	defer file.Close()

	zipWriter := zip.NewWriter(file)

	for _, fileName := range fileNames {

		infile, err := os.Open(fileName)

		data, err := ioutil.ReadAll(infile)
		checkError(err)

		fileWrite, err := zipWriter.Create(fileName)
		checkError(err)

		_, err = fileWrite.Write(data)
		checkError(err)
	}

	err = zipWriter.Close()
	checkError(err)

}

func ExtractArchiveFile(zipName string, targetDirectory string) {

	zipReader, err := zip.OpenReader(zipName)
	checkError(err)
	defer zipReader.Close()

	for _, file := range zipReader.Reader.File {

		zippedfile, err := file.Open()
		checkError(err)
		defer zippedfile.Close()

		if targetDirectory == "" {
			targetDirectory = "./"
		}
		extractedFilePath := filepath.Join(targetDirectory, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extractedFilePath, file.Mode())
		} else {
			outputFile, err := os.OpenFile(extractedFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
			checkError(err)
			defer outputFile.Close()

			_, err = io.Copy(outputFile, zippedfile)
			checkError(err)
		}

	}
}
