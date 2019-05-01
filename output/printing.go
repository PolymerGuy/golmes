package output

import (
	"fmt"
	"gonum.org/v1/gonum/optimize"
)

func PrettyPrint(res *optimize.Result) {
	fmt.Println("Argument values", res.Location.X)
	fmt.Println("Function value at argument values", res.Location.F)
	fmt.Println("Gradient", res.Location.Gradient)
	fmt.Println("Hessian", res.Location.Hessian)
	fmt.Println("Stats", res.Stats)
	fmt.Println("Major iterations", res.Stats.MajorIterations)
	fmt.Println("Function evaluations", res.Stats.FuncEvaluations)
	fmt.Println("Gradient evaluations", res.Stats.GradEvaluations)
	fmt.Println("Run time", res.Stats.Runtime)
	fmt.Println("Convergence status", res.Status)
}
