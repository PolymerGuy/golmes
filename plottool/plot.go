package plottool

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"github.com/PolymerGuy/golmes/data"
)

func PlotSeries(series []plotter.XYs,filename string) {

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"



	for i,serie:=range series {
		l, pts,err := plotter.NewLinePoints(serie)
		l.Color = plotutil.Color(i)
		// If its a point
		if serie.Len() ==1{
			p.Add(pts)
		}else {
			p.Add(l)
		}


		if err != nil {
			panic(err)
		}
	}


	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, filename+".png"); err != nil {
		panic(err)
	}
}


// Makes XY series from an x slice and an y slice
func MakeXYs(xs []float64,ys []float64) plotter.XYs {
	pts := make(plotter.XYs, len(xs))
	for i := range pts {
		pts[i].X = xs[i]
		pts[i].Y = ys[i]
	}
	return pts
}


func PlotResults(comparator data.PairWithArgs,filename string){
	// Plotting, just for fun //

	data := comparator.GetFields()
	ref := MakeXYs(data[0].ReadArgs(),data[0].Read())
	cur := MakeXYs(data[1].ReadArgs(),data[1].Read())

	PlotSeries([]plotter.XYs{ref,cur},filename)
}
