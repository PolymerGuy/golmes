package main

import (
	"os"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"strconv"
	"time"
)

func main() {

	// `os.ReadArgs` provides access to raw command-line
	// arguments. Note that the first value in this slice
	// is the path to the program, and `os.ReadArgs[1:]`
	// holds the arguments to the program.


	inputFile := os.Args[1]
	saveFileAs := strings.TrimSuffix(inputFile, ".txt") + "_res.txt"

	// You can get individual args with normal indexing.
	fmt.Println(inputFile)

	data := readFloatsFromFile(inputFile)
	fmt.Println(data)

	results,args := mysterious3ArgFunction(data)
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Results:",results)
	writeShitToFile(saveFileAs,results,args)


}
func writeShitToFile(saveAsFileName string,data []float64,args []float64) {

	keyword :=[]string{"strain,stress"}
	stringData := floatsToStrings(data)
	stringArgs := floatsToStrings(args)

	stringDatas :=zipStringWithPad(stringArgs,stringData,",")

	strData := append(keyword,stringDatas...)

	output := strings.Join(strData, "\n")
	fmt.Println(output)

	ioutil.WriteFile(saveAsFileName, []byte(output), 0644)




}
func zipStringWithPad(string1 []string, string2 []string, sep string) []string{
	output := []string{}
	for i,word:=range string1{
		output = append(output,strings.Join([]string{word,string2[i]},sep))
	}
	return output

}


func readFloatsFromFile(filename string,)[]float64 {
	results :=[]float64{}


	input, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	lines := strings.Split(string(input), "\n")


	for _, line := range lines {
		num, err := strconv.ParseFloat(line, 32)
		if err != nil {
			fmt.Println(err)
		} else {
			results = append(results, num)}
		}
	return results
}

func parseFloat(str string) float64 {
	s, err := strconv.ParseFloat(str, 32)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func mysterious3ArgFunction(args[]float64)([]float64,[]float64){
	// evals abs(args[0]x²+args[1}x) for x=[-2,-1,1,2]
	//xs :=[]float64{0.,1.2,1.7,1.9,2,2}
	xs := make([]float64, 100)
	for i := 0.0; i < 100.0; i=i+1.0 {
		xs[int(i)] = i*0.006
	}


	results :=[]float64{}
	for _, x:=range xs {
		//results = append(results, math.Abs(args[0]*math.Pow(x,2)+args[1]*x))
		results = append(results, args[2] + bilinear(x,args[0],args[1]))
	}
	return results,xs
	}

func floatsToStrings(floatList []float64) []string {
	stringList := []string{}
	for _, elm := range floatList {
		floatAsString := strconv.FormatFloat(elm, 'e', -1, 32)
		stringList = append(stringList, floatAsString)
	}
	return stringList
}

func bilinear(x float64,sigma_y float64,h float64) float64{
	E := 340.0
	trial := E*x
	if trial < sigma_y{
		return trial
	}else{
		return (x-(sigma_y/E))*h+sigma_y
	}
}
