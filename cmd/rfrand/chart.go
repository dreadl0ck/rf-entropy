package main

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/wcharczuk/go-chart/v2"
)

type measurement struct {
	value float64
	time time.Time
}

var (
	// stats
	inputRates []*measurement
	inputEntropy []*measurement
	
	outputRates []*measurement
	outputEntropy []*measurement
)

func makeChart(name string, data []*measurement, xAxisName, yAxisName string, unitBytes bool) {

	var (
		xValues   []float64
		yValues   []float64
	)

	// collect samples
	for _, m := range data {
		xValues = append(xValues, float64(m.time.UnixNano()))
		yValues = append(yValues, m.value)
	}
	
	// set defaults
	chart.DefaultBackgroundColor = chart.ColorWhite
	chart.DefaultCanvasColor = chart.ColorWhite

	// create chart instance
	graph := chart.Chart{
		Title:      name,
		TitleStyle: chart.Shown(),
		Background: chart.Style{
			Padding: chart.Box{
				Top: 100,
			},
		},
		XAxis: chart.XAxis{
			Name:      xAxisName,
			NameStyle: chart.Shown(),
			Style:     chart.Shown(),
		},
		YAxis: chart.YAxis{
			Name:      yAxisName,
			NameStyle: chart.Shown(),
			Style:     chart.Shown(),
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValueFormatter: func(v interface{}) string {
					return time.Unix(0, int64((v.(float64)))).Format("15:04")
				},
				XValues: xValues,
				YValues: yValues,
				YValueFormatter: func(v interface{}) string {
					if unitBytes {
						return humanize.Bytes(uint64(v.(float64)))
					} else {
						return strconv.FormatFloat(v.(float64), 'f', 4, 64)
					}
				},
			},
		},
	}

	// save the chart to disk
	saveChart(name, graph)
}

// saveChart saves a chart.Graph to the file system
func saveChart(name string, graph chart.Chart) {
	
	// init buffer
	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Println(name, ": failed to render chart:", err)
	}

	// create file
	f, err := os.Create(name)
	if err != nil {
		log.Println(name, ": failed to create chart file:", err)
	}

	// write buf
	_, err = f.Write(buffer.Bytes())
	if err != nil {
		log.Println(name, ": failed to write chart:", err)
	}

	// close handle
	err = f.Close()
	if err != nil {
		log.Println(name, ": failed to close chart file:", err)
	}
}