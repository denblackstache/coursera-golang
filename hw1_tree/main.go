package main

import (
	"fmt"
	"io"
	"os"
)

const entryPrefix = "├───"
const lastEntryPrefix = "└───"
const levelPrefix = "│\t"
const lastLevelPrefix = "\t"

func dirTree(out io.Writer, path string, printFiles bool) error {
	return recursiveDirTree(out, path, printFiles, "")
}

func recursiveDirTree(out io.Writer, path string, printFiles bool, baseLinePrefix string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if !printFiles {
		files = takeDirsOnly(files)
	}

	isLastEntry := false
	for key, file := range files {
		isLastEntry = key == len(files)-1
		line, newPrefix := formatLine(file, baseLinePrefix, isLastEntry)
		_, _ = fmt.Fprint(out, line)

		if file.IsDir() {
			dirPath := path + string(os.PathSeparator) + file.Name()
			err = recursiveDirTree(out, dirPath, printFiles, newPrefix)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func takeDirsOnly(files []os.DirEntry) []os.DirEntry {
	var toPrint []os.DirEntry
	for _, file := range files {
		if file.IsDir() {
			toPrint = append(toPrint, file)
		}
	}

	return toPrint
}

func formatLine(file os.DirEntry, basePrefix string, isLastEntry bool) (string, string) {
	var linePrefix string
	var newPrefix string
	if isLastEntry {
		linePrefix = basePrefix + lastEntryPrefix
		newPrefix = basePrefix + lastLevelPrefix
	} else {
		linePrefix = basePrefix + entryPrefix
		newPrefix = basePrefix + levelPrefix
	}

	line := linePrefix + file.Name()

	if !file.IsDir() {
		fileInfo, err := file.Info()
		fileSize := fileInfo.Size()
		if err != nil {
			line += fmt.Sprintf(" (unknown)")
		}

		if fileSize == 0 {
			line += fmt.Sprintf(" (empty)")
		} else {
			line += fmt.Sprintf(" (%db)", fileSize)
		}
	}
	line += "\n"

	return line, newPrefix
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
