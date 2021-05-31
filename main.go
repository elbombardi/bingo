package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/elbombardi/siego/util"
)

var (
	dataPath   string
	outputPath string
)

func init() {
	flag.StringVar(&dataPath, "data", ".", "Path of the data to be indexed")
	flag.StringVar(&outputPath, "output", "output", "Path of the output directory where the indexes will be stored (if exists will be replaced)")
	flag.Parse()

	data, err := os.Open(dataPath)
	if err != nil {
		fmt.Println("Error: Cannot open data files.", err)
		os.Exit(1)
	}
	info, err := data.Stat()
	if err != nil {
		fmt.Println("Error: Cannot open data files.", err)
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Println("Error : Data path should point to a directory.", dataPath)
	}
	defer data.Close()
	err = os.RemoveAll(outputPath)
	if err != nil && !os.IsNotExist(err) {
		fmt.Println("Error : Cannot delete the output directory.", outputPath)
		os.Exit(1)
	}
	err = os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error : Cannot create the ouput directory.", outputPath)
		os.Exit(1)
	}
}

func createFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error : Cannot create file ", filePath, err)
		os.Exit(1)
	}
	return file
}

func browseDir(dirPath string) []string {
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
			result = append(result, browseDir(entryPath)...)
		} else {
			result = append(result, entryPath)
		}
	}
	return result
}

func stepGenerateMapFile() string {
	const MAP_FILE_NAME = "index.map"
	mapFileFullName := path.Join(outputPath, MAP_FILE_NAME)
	mapFile := createFile(mapFileFullName)
	defer mapFile.Close()
	for _, filePath := range browseDir(dataPath) {
		if util.IsTextFile(filePath) {
			fmt.Fprintln(mapFile, filePath)
		} else {
			fmt.Println(filePath, "is not a text file")
		}
	}
	return mapFileFullName
}

func stringToPath(str string) string {
	result := ""
	for _, c := range str {
		result += "/" + string(c)
	}
	return result
}

// func preprocessField(field string) string {
// 	strings.ReplaceAll()
// }

func removeNonWordCharacters(str string) string {
	re, _ := regexp.Compile(`[^A-Za-z]`)
	return string(re.ReplaceAll([]byte(str), []byte(" ")))
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func generateWordList(source string, fileIndex int, targetDir string) {
	file, err := os.Open(source)
	if err != nil {
		fmt.Println("Error while reading file.", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		line = strings.ToUpper(removeNonWordCharacters(line))
		for _, field := range strings.Fields(line) {
			field = field[0:min(4, len(field))]
			fieldDir := path.Join(targetDir, stringToPath(field))
			err := os.MkdirAll(fieldDir, os.ModePerm)
			if err != nil && !os.IsExist(err) {
				fmt.Println("Cannot created directory : ", err)
				os.Exit(1)
			}
			ioutil.WriteFile(path.Join(fieldDir, strconv.Itoa(fileIndex)), []byte(""), os.ModePerm)
		}
	}
}

func stepGenerateWordLists(mapFileName string) {
	const WORD_LIST_DIR = "words"
	wordListFullPath := path.Join(outputPath, WORD_LIST_DIR)
	mapFile, err := os.Open(mapFileName)
	if err != nil {
		fmt.Println("Error while opening file :", err)
		os.Exit(1)
	}
	defer mapFile.Close()

	err = os.MkdirAll(wordListFullPath, os.ModePerm)
	if err != nil {
		fmt.Println("Error : Cannot create the output directory.", outputPath)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(mapFile)
	scanner.Split(bufio.ScanLines)
	index := 1
	for scanner.Scan() {
		fileName := scanner.Text()
		generateWordList(fileName, index, wordListFullPath)
		index++
	}
}

func main() {
	//Map the files, store the map into index.map. Replace full path with a page indentifier : UPI
	mapFile := stepGenerateMapFile()

	//For each file, generate a list of words store in <UPI>.lst
	stepGenerateWordLists(mapFile)

	//Regroup all the lists in one single list (index.all)
	// stepGenerateInvertedIndex()

	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for inputScanner.Scan() {
		fmt.Println(inputScanner.Text())
		fmt.Print("> ")
	}

}
