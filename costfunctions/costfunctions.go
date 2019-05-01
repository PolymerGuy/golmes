package costfunctions

type CostFunction interface {
	Eval(args []float64) float64
}
