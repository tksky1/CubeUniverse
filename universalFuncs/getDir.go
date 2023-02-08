package universalFuncs

import (
	"os"
	"path/filepath"
	"strings"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

// GetParentDir 获取运行目录的父目录
func GetParentDir() string {
	directory, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return substr(directory, 0, strings.LastIndex(directory, "/"))
}

// GetCurrentDir 获取当前运行目录
func GetCurrentDir() string {
	directory, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return directory
}
