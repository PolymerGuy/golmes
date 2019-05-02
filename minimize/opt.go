package minimize

import (
	"github.com/PolymerGuy/golmes/costfunctions"
	"github.com/PolymerGuy/golmes/maths"
	"github.com/PolymerGuy/golmes/yamlparser"
	"github.com/btracey/meshgrid"
	"gonum.org/v1/gonum/bound"
	"gonum.org/v1/gonum/optimize"
	"log"
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

type UniformGrid struct {
	Bounds []bound.Bound
	Seeds  []int

	curpoint int
	grid     [][]float64
}

func (n *UniformGrid) Rand(x []float64) []float64 {
	pts := []float64{}
	log.Println("calling rand")

	if n.grid == nil {
		n.makeGrid()
	}

	if x == nil {
		return x
	}

	log.Println("Point", n.curpoint)

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

	n.grid = pts
	log.Println("Made the grid: ", pts)
}

func makeUniformGrid(bounds []bound.Bound, seeds []int) UniformGrid {
	return UniformGrid{Bounds: bounds,
		Seeds:    seeds,
		curpoint: 0}
}
