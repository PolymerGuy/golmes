package minimize

import (
	"fmt"
	"github.com/PolymerGuy/golmes/costfunctions"
	"github.com/PolymerGuy/golmes/maths"
	"github.com/PolymerGuy/golmes/yamlparser"
	"github.com/btracey/meshgrid"
	"github.com/polymerguy/gorbi"
	"gonum.org/v1/gonum/bound"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/optimize"
	"log"
	"math"
	"os"
	"time"
)

func FindFunctionMinima(optJob OptimizationJob) (*optimize.Result, error) {

	p := optimize.Problem{
		Func: optJob.CostFunc.Eval,
		Grad: nil,
	}

	var method = optJob.Method
	optJob.Settings.Recorder = optimize.NewPrinter()

	res, err := optimize.Minimize(p, optJob.InitialParameters, &optJob.Settings, method)
	if err != nil {
		log.Fatal(err)
	}
	return res, err
}

type OptimizationJob struct {
	CostFunc          costfunctions.CostFunction
	Method            optimize.Method
	InitialParameters []float64
	Settings          optimize.Settings
}

func CoarseSearch(optJob OptimizationJob, coarse yamlparser.CoarseSearchSettings) (*optimize.Result, error) {

	p := optimize.Problem{
		Func: optJob.CostFunc.Eval,
	}

	dim := len(coarse.Seeds)

	nPts := 1
	for _, seed := range coarse.Seeds {
		nPts *= seed
	}

	bounds := []bound.Bound{}
	for i, _ := range coarse.Seeds {
		bounds = append(bounds, bound.Bound{coarse.Bounds[i*2], coarse.Bounds[i*2+1]})
	}

	d := makeUniformGrid(bounds, coarse.Seeds)

	method := optimize.GuessAndCheck{Rander: &d}

	optJob.Settings.MajorIterations = nPts
	optJob.Settings.FuncEvaluations = nPts
	optJob.Settings.Converger = &optimize.FunctionConverge{Absolute: 1e-3,
		Iterations: nPts}
	printer := &optimize.Printer{
		Writer:          os.Stdout,
		HeadingInterval: 30,
		ValueInterval:   100 * time.Millisecond,
	}
	optJob.Settings.Recorder = printer
	optJob.Settings.Concurrent = 4

	initX := make([]float64, dim)

	res, err := optimize.Minimize(p, initX, &optJob.Settings, &method)

	if err != nil {
		log.Fatal(err)
	}
	return res, err

}

func CoarseSearchSurf(optJob OptimizationJob, coarse yamlparser.CoarseSearchSettings)[]float64{

	p := optimize.Problem{
		Func: optJob.CostFunc.Eval,
	}


	nPts := 1
	for _, seed := range coarse.Seeds {
		nPts *= seed
	}

	bounds := []bound.Bound{}
	for i, _ := range coarse.Seeds {
		bounds = append(bounds, bound.Bound{coarse.Bounds[i*2], coarse.Bounds[i*2+1]})
	}


	d := makeUniformGrid(bounds, coarse.Seeds)

	args := [][]float64{}
	vals :=[]float64{}

	for range maths.Linspace(0.,.1,nPts){
		point := d.Rand([]float64{0.,0.0})

		args = append(args, point)
		vals = append(vals,p.Func(point))

	}


	fineSeeds := []int{}
	for _,seed := range coarse.Seeds{
		fineSeeds = append(fineSeeds,seed*50)
	}

	dFine := makeUniformGrid(bounds, fineSeeds)
	nPtsFine := 1
	for _, seed := range fineSeeds {
		nPtsFine *= seed
	}


	argsFine := [][]float64{}
	for range maths.Linspace(0.,.1,nPtsFine){
		point := dFine.Rand([]float64{0.,0.0})

		argsFine = append(argsFine, point)

	}

	fmt.Println("Checkling points",len(argsFine))

	rbi := gorbi.NewRBF(args,vals)
	rbiVals := rbi.At(argsFine)
	best :=floats.Min(rbiVals)
	bestArg := floats.MinIdx(rbiVals)

	fmt.Println("Stuff!")
	fmt.Println("Evaluated N pts: ",fineSeeds)
	fmt.Println("Best values was: ",floats.Min(vals))
	fmt.Println("at argument: ",argsFine[floats.MinIdx(vals)])
	fmt.Println("Best values is: ",best)
	fmt.Println("at argument: ",argsFine[bestArg])
	fmt.Println("Checked the best value and it is:",p.Func(argsFine[bestArg]))
	fmt.Println("Starting fine search")
	fmt.Println("")

	boundsFine := []bound.Bound{}
	scale := 0.3
	cent :=argsFine[bestArg]

	for i, _ := range coarse.Seeds {
		span := math.Abs(coarse.Bounds[i*2] - coarse.Bounds[i*2+1])
		min := cent[i] - span*scale/2.
		max := cent[i] + span*scale/2.

		boundsFine = append(boundsFine, bound.Bound{min, max})
	}

	fmt.Println("Bounds:",boundsFine)

	dFines := makeUniformGrid(boundsFine, coarse.Seeds)
	for range maths.Linspace(0.,.1,nPts){
		point := dFines.Rand([]float64{0.,0.0})
		args = append(args, point)
		vals = append(vals,p.Func(point))

	}

	bestFine :=floats.Min(vals)
	bestArgFine := floats.MinIdx(vals)
	bestFinePos := args[bestArgFine]

	fmt.Println("Second search gave:")
	fmt.Println(bestFine)
	fmt.Println(bestFinePos)



//	newPoint := argsFine[bestArg]
//	for range nLooks{
//		newValue := p.Func(newPoint)
//		args = append(args,newPoint)
//		vals = append(vals,newValue)
//
//		rbi := gorbi.NewRBF(args,vals)
//
//		rbiVals := rbi.At(argsFine)
//		bestArg := floats.MinIdx(rbiVals)
//		newPoint =  argsFine[bestArg]
//		fmt.Println("Best value is: ",floats.Min(rbiVals))
//		fmt.Println("With arg: ",newPoint)
//
//	}



	return argsFine[bestArg]



}


type UniformGrid struct {
	Bounds []bound.Bound
	Seeds  []int

	curpoint int
	grid     [][]float64
}

func (n *UniformGrid) Rand(x []float64) []float64 {
	pts := []float64{}

	if n.grid == nil {
		n.makeGrid()
	}

	if x == nil {
		return x
	}


	if n.curpoint >= len(n.grid) {
		return nil
	}

	pts = n.grid[n.curpoint]
	n.curpoint += 1
	for i, _ := range x {
		x[i] = pts[i]
	}

	return x
}

func (n *UniformGrid) makeGrid() {

	gridPts := [][]float64{}
	for i, bound := range n.Bounds {
		gridPts = append(gridPts, maths.Linspace(bound.Min, bound.Max, n.Seeds[i]))

	}

	grid := gridPts
	pts := meshgrid.Multiple(grid)
	fmt.Println("Grid: ",gridPts)

	n.grid = pts
}

func makeUniformGrid(bounds []bound.Bound, seeds []int) UniformGrid {
	return UniformGrid{Bounds: bounds,
		Seeds:    seeds,
		curpoint: 0}
}
