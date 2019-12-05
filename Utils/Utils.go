package Utils

import (
	"os"
	"time"
)

//
func PathIsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func CurrentTime() string {
	t := time.Now()
	str := t.Format("2006-01-02 15:04:05")

	return str
}

func DeleteFile(filePath string) bool {
	errRemove := os.Remove(filePath)

	if errRemove != nil {
		return false
	}
	return true
}
