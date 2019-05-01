package fileutils

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func MakeInputFile(templateFileName string,workdir string,  values []float64, placeholders []string) (string, error) {
	filename :=  filepath.Base(templateFileName)

	inputFileName := workdir +"/"+ strings.TrimSuffix(filename, ".txt") + "_iter.txt"

	inputFileName = workdir +"/"+ filename


	//err := os.MkdirAll(id,0755)
	//if err != nil{
//		log.Panic(err)
//	}

	stringArgs := floatsToStrings(values)
	substituteInFile(templateFileName, inputFileName, placeholders, stringArgs)
	return inputFileName, nil

}

func substituteInFile(filename string, saveAsFilename string, targets []string, substitutes []string) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}
	output := replaceStrings(input, targets, substitutes)

	err = ioutil.WriteFile(saveAsFilename, output, 0644)
	if err != nil {
		log.Fatalln(err)

	}
}
func replaceStrings(input []byte, targets []string, substitutes []string) []byte {
	//TODO: Strange composition of input types?
	output := input

	for i, target := range targets {
		if found := bytes.Contains(input, []byte(target)); found != true {
			log.Fatal("Did not find key in file")
		}
		output = bytes.Replace(output, []byte(target), []byte(substitutes[i]), -1)
	}
	return output

}

func floatsToStrings(floatList []float64) []string {
	stringList := []string{}
	for _, elm := range floatList {
		floatAsString := strconv.FormatFloat(elm, 'e', -1, 32)
		stringList = append(stringList, floatAsString)
	}
	return stringList
}

func GetKeyFromCSVFile(fileName string, key string) []float64 {

	csvFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	line, err := reader.Read()
	if err != nil {
		log.Fatalf("Reading from file did not work: %s", err)
	}
	ind, err := findWordInd(line, key)
	// Check if key was found
	if err != nil {
		log.Fatalf("Did not find key. %s", err)
	}
	row, err := floatsColumn(reader, ind)
	if err != nil {
		log.Fatalf("Did not find key. %s", err)
	}

	return row
}
// floatsColumn get all floats contained in a csv column, skips non-float values
func floatsColumn(reader *csv.Reader, columnInd int) ([]float64, error) {
	results := []float64{}
	for {
		line, err := reader.Read()
		// Break if at end of file
		if err == io.EOF {
			break
		} else if err != nil {
			return results, err
		}

		if len(line)<columnInd{
			return results,errors.New("Column id is out of range")
		}else {
			num, err := parseFloat(line[columnInd])
			if err == nil {
				results = append(results, num)
			}
		}


	}
	return results, nil
}
func findWordInd(str []string, key string) (int, error) {

	// Find column matching the keyword.
	// Assume this to be on the top row and occurring only once

	for rowIndex, word := range str {
		if word == key {
			return rowIndex, nil

		}

	}
	return -1, errors.New("String not found in slice of strings")

}

func parseFloat(str string) (float64, error) {
	s, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Could not parse float: %s", err))

	}
	return s, nil
}
