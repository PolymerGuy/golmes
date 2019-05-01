package main

import (
	"github.com/PolymerGuy/golmes/minimize"
	"github.com/PolymerGuy/golmes/output"
	"gonum.org/v1/gonum/optimize"
	"log"
)

func main(){

	booth := NewFunction(NewBoothFunction())
	boothReader := FunctionReader{F: booth}


	optJob := minimize.OptimizationJob{
		CostFunc:&boothReader,
		Method:&optimize.GradientDescent{},
		InitialParameters:[]float64{1,1},
		Settings:optimize.Settings{FuncEvaluations:10},
	}



	res, err := minimize.FindFunctionMinima(optJob)

	if err != nil{
		log.Panic(err)
	}
	output.PrettyPrint(res)
}

