package functions

import (
	"log"
)

type function struct {
	//TODO: Better naming
	input  []float64
	expr   Expression
	output []float64
}

type Expression struct {
	expr func(args... float64) []float64
	nArgs int
}


type Function interface {
	Args([]float64)
	Evaluate()
	Value() []float64
}


func NewFunction(expr Expression) Function{
	//TODO: Understand why i need to return a pointer...
	return &function{expr:expr}
}


func (f *function) Args(input []float64){
	//TODO: Do some checks
	if len(input)!=f.expr.nArgs{
		log.Panicf("Number of arguments does not match Function")
	}
	f.input = input

}


func (f *function) Value() []float64{
	//TODO: Rework output type
	return f.output
}

func (f *function) Evaluate(){
	f.output = f.expr.expr(f.input...)

}


type FunctionReader struct {
	F Function
}

func (wf *FunctionReader) Eval(args []float64) float64{
	wf.F.Args(args)
	wf.F.Evaluate()
	return wf.F.Value()[0]
}





