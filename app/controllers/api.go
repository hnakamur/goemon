package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"
	"github.com/robfig/revel"
	"github.com/ziutek/rrd"
)

type RrdFetchTsvResult struct {
	data rrd.FetchResult
	origStep int
}

func (r RrdFetchTsvResult) Apply(req *revel.Request, resp *revel.Response) {
	revel.TRACE.Printf("Apply start\n")
	resp.WriteHeader(http.StatusOK, "text/tab-separated-values")

	resp.Out.Write([]byte("date\tvalue\n"))
	data := r.data
	row := 0
	for ti := data.Start.Add(data.Step);
		ti.Before(data.End) || ti.Equal(data.End);
		ti = ti.Add(data.Step) {
		v := data.ValueAt(0, row)
		line := fmt.Sprintf("%d\t%e\n", ti.Unix(), v)
		resp.Out.Write([]byte(line))
		row++
	}
	revel.TRACE.Printf("Apply exit\n")
/*
	origStep := r.origStep * 1000
	jStep := origStep / int(step)
	if jStep < 1 {
		jStep = 1
	}

	// TODO: Calculate max, average and such if jStep > 1
	revel.TRACE.Printf("jStep=%d, step=%d, origStep=%d\n", jStep, step, origStep)
	for j := 0; j < data.RowLen; j += jStep {
		t := start + int64(j + 1) * step
		v := data.ValueAt(0, j)
		if j > 0 && math.IsNaN(v) {
			break
		}
		line := fmt.Sprintf("%d\t%e\n", t, v)
		resp.Out.Write([]byte(line))
	}
*/
}

type Api struct {
	*revel.Controller
}

func (c Api) DataTsv(path, cf string, start, end time.Time, step int) revel.Result {
	format, found := revel.Config.String("rrd.path_format")
	if !found {
		return c.RenderError(errors.New("config rrd.path_format not found"))
	}
	rrdpath := fmt.Sprintf(format, path)

	if cf == "" {
		cf = "AVERAGE";
	}

	revel.TRACE.Printf("DataTsv start=%s, end=%s, step=%d", start, end, step)
	if end.IsZero() {
		end = time.Now()
		revel.TRACE.Printf("DataTsv end was zero, now %s", end)
	}
	if start.IsZero() {
		start = end.Add(-36 * time.Hour)
		revel.TRACE.Printf("DataTsv start was zero, now %s", start)
	}

	if step == 0 {
		step = 60;
		revel.TRACE.Printf("DataTsv step was zero, now %d", step)
	}
	stepDuration := time.Duration(step) * time.Second

    data, err := rrd.Fetch(rrdpath, cf, start, end, stepDuration)
    if err != nil {
		revel.TRACE.Printf("Fetch fails %s", err)
		return c.RenderError(err)
    }
	defer data.FreeValues()

	return RrdFetchTsvResult{data, step}
}
