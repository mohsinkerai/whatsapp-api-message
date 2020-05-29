
package main

import (
	//"bufio"
	"encoding/csv"
	"fmt"
	// "io"
	"log"
	"os"
)

func main() {
	// Open the file
	res := parseCsv("input.csv")

	for _, r := range res {
		fmt.Printf("\nPhone number is %s", r[0])
	}
}

func parseCsv(fileName string) [][]string {
	csvfile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln("Couldn't open the csv file", err)
	}

	// Parse the file
	r := csv.NewReader(csvfile)
	//r := csv.NewReader(bufio.NewReader(csvfile))

	rows, err := r.ReadAll()

	return rows
}
