package utils

import (
	"fmt"
	"os"
	"path"
)

func CreateFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error : Cannot create file ", filePath, err)
		os.Exit(1)
	}
	return file
}

func BrowseDir(dirPath string) []string {
	dir, err := os.Open(dirPath)
	if err != nil {
		fmt.Println("Error while opening directory :", err)
		os.Exit(1)
	}
	defer dir.Close()
	entries, err := dir.ReadDir(-1)
	if err != nil {
		fmt.Println("Error while opening directory :", err)
		os.Exit(1)
	}
	result := []string{}
	for _, entry := range entries {
		entryPath := path.Join(dirPath, entry.Name())
		if entry.IsDir() {
			result = append(result, BrowseDir(entryPath)...)
		} else {
			result = append(result, entryPath)
		}
	}
	return result
}

// textExt[x] is true if the extension x indicates a text file, and false otherwise.
var textExt = map[string]bool{
	".md":  true, // must be served raw
	".txt": true, // must be served raw
	".csv": true, // must be served raw

}

// IsTextFile reports whether the file has a known extension indicating
// a text file, or if a significant chunk of the specified file looks like
// correct UTF-8; that is, if it is likely that the file contains human-
// readable text.
func IsTextFile(filename string) bool {
	// if the extension is known, use it for decision making
	if isText, found := textExt[path.Ext(filename)]; found {
		return isText
	}
	return false
	// the extension is not known; read an initial chunk
	// of the file and check if it looks like text
	/*f, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer f.Close()

	var buf [1024]byte
	n, err := f.Read(buf[0:])
	if err != nil {
		return false
	}

	return IsText(buf[0:n])*/
}
