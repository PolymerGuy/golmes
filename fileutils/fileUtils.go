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
	"path"
	"strconv"
	"strings"
)

func MakeInputFile(templateFileName string, workdir string, values []float64, placeholders []string) (string, error) {
	filename := path.Base(templateFileName)

	inputFileName := strings.TrimSuffix(filename, ".txt") + "_iter.txt"

	inputFilePath := path.Join(workdir, inputFileName)

	valueStrings := floatsToStrings(values)

	err := makeInputFileFromTemplate(templateFileName, inputFilePath, placeholders, valueStrings)
	if err != nil {
		return inputFilePath, errors.New("Could not make the file: " + inputFilePath + "\n because: " + err.Error())
	}

	return inputFilePath, nil

}

func makeInputFileFromTemplate(templateFilePath string, saveAs string, placeholders []string, substitutes []string) error {
	input, err := ioutil.ReadFile(templateFilePath)
	if err != nil {
		return errors.New("Could not read the template file: " + err.Error())
	}

	output, err := replaceStrings(input, placeholders, substitutes)
	if err != nil {
		return errors.New("Could not replace the keywords in the template file: " + err.Error())
	}

	err = ioutil.WriteFile(saveAs, output, 0644)
	if err != nil {
		return errors.New("Could not write the input file to disc: " + err.Error())

	}
	return nil
}
func replaceStrings(input []byte, targets []string, substitutes []string) ([]byte, error) {
	//TODO: Strange composition of input types?
	output := input

	for i, target := range targets {
		if found := bytes.Contains(input, []byte(target)); found != true {
			return nil, errors.New("Did not find the key:" + target)
		}
		output = bytes.Replace(output, []byte(target), []byte(substitutes[i]), -1)
	}
	return output, nil

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

		if len(line) < columnInd {
			return results, errors.New("Column id is out of range")
		} else {
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
