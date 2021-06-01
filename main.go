package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/elbombardi/siego/index"
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

func main() {
	siegoIndex := index.SiegoIndex{
		Target: dataPath,
	}
	siegoIndex.Index()
	// siegoIndex.PrintEntries()
	siegoIndex.Save(path.Join(outputPath, "index.sie"))

	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for inputScanner.Scan() {
		query := inputScanner.Text()
		locations := siegoIndex.Lookup(query)
		if len(locations) == 0 {
			fmt.Println("<No result!>")
		}
		for locInd, loc := range locations {
			fmt.Printf("%v occurence(s) found in '%s'\n", loc.Counter, siegoIndex.LocationsMap[locInd])
		}
		fmt.Print("> ")
	}

}
