package minimize

import (
	"github.com/PolymerGuy/golmes/costfunctions"
	"github.com/PolymerGuy/golmes/maths"
	"github.com/PolymerGuy/golmes/yamlparser"
	"github.com/btracey/meshgrid"
	"gonum.org/v1/gonum/bound"
	"gonum.org/v1/gonum/diff/fd"
	"gonum.org/v1/gonum/optimize"
	"log"
	"os"
	"time"
)

func FindFunctionMinima(optJob OptimizationJob) (*optimize.Result, error) {

	grad := gradWrapper{f: optJob.CostFunc.Eval}

	p := optimize.Problem{
		Func: optJob.CostFunc.Eval,
		Grad: grad.Gradient,
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

type gradWrapper struct {
	f func([]float64) float64
}

func (g gradWrapper) Gradient(grad []float64, x []float64) {
	// Compute the first derivative of f at 0 using the default settings.
	settings := fd.Settings{Step: 0.001}
	fd.Gradient(grad, g.f, x, &settings)
}

func CoarseSearch(optJob OptimizationJob, coarse yamlparser.CoarseSearchSettings) (*optimize.Result, error) {
	log.Println("Initiating coarse search")

	p := optimize.Problem{
		Func: optJob.CostFunc.Eval,
	}

	dim := len(coarse.Seeds)

	log.Println("BOUNDS:", coarse.Bounds)
	log.Println("Seeds:", coarse.Seeds)

	nPts := 1
	for _, seed := range coarse.Seeds {
		nPts *= seed
	}

	bounds := []bound.Bound{}
	for i, _ := range coarse.Seeds {
		bounds = append(bounds, bound.Bound{coarse.Bounds[i*2], coarse.Bounds[i*2+1]})
	}

	d := UniformGrid{Bounds: bounds, Seeds: coarse.Seeds}

	method := optimize.GuessAndCheck{Rander: &d}

	optJob.Settings.MajorIterations = nPts
	optJob.Settings.FuncEvaluations = nPts
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

	if n.grid == nil {
		n.makeGrid()
	}

	if x == nil {
		return x
	}

	pts = n.grid[n.curpoint]
	n.curpoint += 1
	//x=n.grid[n.curpoint]
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
}
