package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/elbombardi/siego/index"
)

var (
	dataPath   string
	indexPath  string
	outputPath string
)

func init() {
	flag.StringVar(&dataPath, "data", "",
		"Path of the data to be indexed")
	flag.StringVar(&indexPath, "index", "",
		"Path of a pre-saved index")
	flag.StringVar(&outputPath, "output", "",
		"Path of the output directory where the indexes will be stored (if exists will be replaced)")
	flag.Parse()

	if dataPath == "" && indexPath == "" {
		fmt.Println("Please specify a data path (-data) or a path to a presaved index file (-index)")
		os.Exit(1)
	}

	if dataPath != "" {
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
	}

	if indexPath != "" {
		index, err := os.Open(indexPath)
		if err != nil {
			fmt.Println("Error : Cannot open index file.", err)
			os.Exit(1)
		}
		defer index.Close()
	}
}

func main() {
	var siegoIndex *index.SiegoIndex
	var err error
	if dataPath != "" {
		siegoIndex = &index.SiegoIndex{
			Target: dataPath,
		}
		siegoIndex.Index()
		if outputPath != "" {
			err := siegoIndex.Save(outputPath)
			if err != nil {
				fmt.Println("Error : Cannot save index ", err)
				os.Exit(1)
			}
		}
	} else if indexPath != "" {
		siegoIndex, err = index.Load(indexPath)
		if err != nil {
			fmt.Println("Error : Cannot load index ", err)
			os.Exit(1)
		}
	}
	// siegoIndex.PrintEntries()

	inputScanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for inputScanner.Scan() {
		query := inputScanner.Text()
		locations, found := siegoIndex.Lookup(query)
		if !found {
			fmt.Print("<No result!>\n> ")
			continue
		}
		fmt.Printf("Found in %v out of %v document(s) : \n",
			len(locations), len(siegoIndex.DocumentsMap))
		if len(locations) > 5 {
			fmt.Print("Do you want to see all? (y/N) > ")
			inputScanner.Scan()
			answer := strings.ToLower(strings.TrimSpace(inputScanner.Text()))
			if answer == "" || answer == "n" {
				fmt.Print("> ")
				continue
			}
		}
		for i := 0; i < len(siegoIndex.DocumentsMap); i++ {
			loc, ok := locations[i]
			if !ok {
				continue
			}
			fmt.Printf("%v occurence(s) found in \"%s\".\n",
				len(loc.Positions), siegoIndex.DocumentsMap[i].Name)
		}
		fmt.Print("> ")
	}
}
