package minimize

import (
	"fmt"
	"github.com/PolymerGuy/golmes/costfunctions"
	"github.com/PolymerGuy/golmes/maths"
	"github.com/PolymerGuy/golmes/yamlparser"
	"github.com/PolymerGuy/gorbi"
	"github.com/btracey/meshgrid"
	"gonum.org/v1/gonum/bound"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/optimize"
	"log"
	"math"
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

func CoarseSearchSurf(optJob OptimizationJob, coarse yamlparser.CoarseSearchSettings) ([]float64,error) {
	args := [][]float64{}
	vals := []float64{}
	bestArg := []float64{}

	nDims := len(coarse.Seeds)


	nScales := 3
	refinementFactor := 0.1
	fineSearchUpsampling := 50
	d := makeUniformGrid(coarse.Bounds, coarse.Seeds)

	for i:=0;i<nScales;i++ {

		// Evaluate all points on the grid and store args and values
		for range maths.Linspace(0., .1, coarse.NPts) {
			dummyArg := make([]float64, nDims)
			point := d.Rand(dummyArg)
			args = append(args, point)
			vals = append(vals, optJob.CostFunc.Eval(point))

		}

		// Make a finer grid on which the interpolation wil be evaluated
		fineSeeds := []int{}
		for _, seed := range coarse.Seeds {
			fineSeeds = append(fineSeeds, seed*fineSearchUpsampling)
		}
		dFine := makeUniformGrid(coarse.Bounds, fineSeeds)
		nPtsFine := coarse.NPts * int(math.Pow(float64(fineSearchUpsampling),float64(nDims)))

		argsFine := [][]float64{}
		for range maths.Linspace(0., .1, nPtsFine) {
			dummyArg := make([]float64, nDims)
			point := dFine.Rand(dummyArg)
			argsFine = append(argsFine, point)

		}

		fmt.Println("Checkling points", len(argsFine))
		fmt.Println("Checkling points", nPtsFine)

		rbi,err := gorbi.NewRBF(args, vals)
		if err != nil{
			log.Fatal(err)
		}

		rbiVals := rbi.At(argsFine)
		best := floats.Min(rbiVals)
		bestArg = argsFine[floats.MinIdx(rbiVals)]

		fmt.Println("Beste values is:",best)
		fmt.Println("found at:",bestArg)



		boundsFine := []float64{}
		cent := bestArg

		for i, _ := range coarse.Seeds {
			span := math.Abs(coarse.Bounds[i*2] - coarse.Bounds[i*2+1])
			min := cent[i] - span*refinementFactor/2.
			max := cent[i] + span*refinementFactor/2.

			boundsFine = append(boundsFine, min, max)
		}

		fmt.Println("Bounds:", boundsFine)

		d = makeUniformGrid(boundsFine, coarse.Seeds)






	}

	return bestArg,nil

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
	//fmt.Println("Grid: ",gridPts)

	n.grid = pts
}

func makeUniformGrid(bounds []float64, seeds []int) UniformGrid {
	nPairs := len(bounds) / 2

	boundss := []bound.Bound{}
	for i := 0; i < nPairs; i++ {
		boundss = append(boundss, bound.Bound{bounds[i*2], bounds[i*2+1]})
	}

	return UniformGrid{Bounds: boundss,
		Seeds:    seeds,
		curpoint: 0}
}
