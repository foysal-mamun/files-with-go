package files

import (
	"archive/zip"
	"compress/gzip"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// check if error, then stop
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

// Just create an empty file
func CreateEmptly(fileName string) {

	newFile, err := os.Create(fileName)
	checkError(err)
	defer newFile.Close()

	log.Println(newFile)
}

// Truncate a file by given size
func Truncate(fileName string, size int64) {

	err := os.Truncate(fileName, size)
	checkError(err)
}

// Display file information
func GetInfo(fileName string) {

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

// Move file to new location
func Move(oldLoc string, newLoc string) {

	err := os.Rename(oldLoc, newLoc)
	checkError(err)
}

// Delete a file
func Delete(fileName string) {

	err := os.Remove(fileName)
	checkError(err)

}

// Example of how to open a file.
func Open(fileName string) {

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(os.IsNotExist(err))
	}
	file.Close()

	file, err = os.OpenFile(fileName, os.O_APPEND, 0666)
	checkError(err)
	file.Close()
}

// Check write permission
func CheckPermission(fileName string) {

	file, err := os.OpenFile(fileName, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(os.IsPermission(err))
	}
	file.Close()
}

// Change file permission
func ChangePermission(fileName string, mode int) {

	err := os.Chmod(fileName, os.FileMode(mode))
	checkError(err)
}

// Change a file ownership
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

// Change file timestamps
func ChangeTime(fileName string, atime, mtime time.Time) {
	if err := os.Chtimes(fileName, atime, mtime); err != nil {
		log.Fatal(err)
	}
}

// Create a hard link
func HardLink(oldname, newname string) {
	if oldname == "" || newname == "" {
		return
	}

	if err := os.Link(oldname, newname); err != nil {
		log.Fatal(err)
	}
}

// Create symbolic link
func SymLink(oldname, newname string) {
	if err := os.Symlink(oldname, newname); err != nil {
		log.Fatal(err)
	}
}

// Copy file (hard link)
func Copy(oldname, newname string) (err error) {

	oldfile, err := os.Stat(oldname)
	checkError(err)
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

// Copy content to new file
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

// seek file
func Seek(fileName string, offset int64, whence int) int64 {

	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()

	newPos, err := file.Seek(offset, whence)
	checkError(err)

	return newPos
}

// Write content to a file
func Write(fileName string, content []byte) {

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	checkError(err)
	defer file.Close()

	_, err = file.Write(content)
	checkError(err)
}

// Read file by given length
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

// Archive given files
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

// Extract a zip archive file
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

// Compress "fileName" to zip "gzFileName"
func CompressFile(gzFileName string, fileName string) {

	outfile, err := os.Create(gzFileName)
	checkError(err)

	gzipWriter := gzip.NewWriter(outfile)
	defer gzipWriter.Close()

	inFile, err := os.Open(fileName)
	checkError(err)

	data, err := ioutil.ReadAll(inFile)
	checkError(err)

	_, err = gzipWriter.Write(data)
	checkError(err)
}

// Uncompress file "gzFileName" to "newFileName"
func UncompressFile(gzFileName string, newFileName string) {

	gzipFile, err := os.Open(gzFileName)
	checkError(err)
	defer gzipFile.Close()

	gzipReader, err := gzip.NewReader(gzipFile)
	checkError(err)
	defer gzipReader.Close()

	outFileWriter, err := os.Create(newFileName)
	checkError(err)
	defer outFileWriter.Close()

	_, err = io.Copy(outFileWriter, gzipReader)
	checkError(err)

}

// Create a temporary file which will remove at file end
func CerateTempFile(fileName string) {

	tempFile, err := ioutil.TempFile("", fileName)
	checkError(err)
	defer func() {
		tempFile.Close()
		removeTempFile(tempFile.Name())
	}()

}

// remove a afile
func removeTempFile(fileName string) {
	err := os.Remove(fileName)
	checkError(err)
}

// Read file from HTTP and write content to a file.
func FileFromHTTP(fileName, url string) {

	newFile, err := os.Create(fileName)
	checkError(err)
	defer newFile.Close()

	res, err := http.Get(url)
	checkError(err)
	defer res.Body.Close()

	_, err = io.Copy(newFile, res.Body)
	checkError(err)
}

// Create checksum from a file content
func ChecksumFileContent(fileName string) {

	data, err := ioutil.ReadFile(fileName)
	checkError(err)

	fmt.Printf("Md5: %x\n\n", md5.Sum(data))
	fmt.Printf("Sha1: %x\n\n", sha1.Sum(data))
	fmt.Printf("Sha256: %x\n\n", sha256.Sum256(data))
	fmt.Printf("Sha512: %x\n\n", sha512.Sum512(data))
}

// Create checksum by file handler
func ChecksumFile(fileName string) {
	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()

	hasher := md5.New()
	_, err = io.Copy(hasher, file)
	checkError(err)

	sum := hasher.Sum(nil)
	fmt.Printf("Md5 checksum: %x\n", sum)
}
