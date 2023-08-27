package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
)

func GetDirPath() string {

	tempDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(tempDir, "data")
}

func Diff(filea string, fileb string) (string, uint64, uint64) {

	// a.txt와 b.txt 파일을 엽니다.
	a, err := os.Open(filea)
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	b, err := os.Open(fileb)
	if err != nil {
		log.Fatal(err)
	}
	defer b.Close()

	// 두 파일의 내용을 문자열로 변환합니다.
	text1, err := io.ReadAll(a)
	if err != nil {
		log.Fatal(err)
	}
	text2, err := io.ReadAll(b)
	if err != nil {
		log.Fatal(err)
	}

	// // 두 파일의 내용을 문자열로 변환합니다.
	file1 := string(text1)
	file2 := string(text2)

	slice1 := strings.Split(file1, "\n")
	slice2 := strings.Split(file2, "\n")
	diff := cmp.Diff(slice1, slice2)

	// result.txt 파일을 생성하고, 변경사항을 씁니다.
	log.Println("diff : ", diff)

	getSlice := strings.Split(diff, "\n")
	insertSum := 0
	deleteSum := 0
	for i := 0; i < len(getSlice)-1; i++ {
		if getSlice[i][:1] == "+" {
			insertSum++
		}
		if getSlice[i][:1] == "-" {
			deleteSum++
		}
	}
	return diff, uint64(insertSum), uint64(deleteSum)
}

func GetSize(filepath string) {

}
