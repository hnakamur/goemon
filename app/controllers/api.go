package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/robfig/revel"
	"github.com/ziutek/rrd"
	"github.com/hnakamur/goemon/util/dir"
)

import Ti "github.com/hnakamur/goemon/util/time"

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

func (c Api) DataTsv(path, cf, duration, end string, step int) revel.Result {
	dataDir, found := revel.Config.String("rrd.data_dir")
	if !found {
		return c.RenderError(errors.New("config rrd.data_dir not found"))
	}
	rrdpath := fmt.Sprintf("%s/%s", dataDir, path)

	if cf == "" {
		cf = "AVERAGE"
	}

	var endTime time.Time
	revel.TRACE.Printf("DataTsv end=%s", end)
	if end == "" {
		endTime = time.Now()
	} else {
		endInt, err := strconv.ParseInt(end, 10, 64)
		if err != nil {
			return c.RenderError(err)
		}
		endTime = time.Unix(endInt, 0)
	}
	revel.TRACE.Printf("DataTsv endTime=%s", endTime)

	if duration == "" {
		duration = "1h"
	}
	graphDuration, err := Ti.ParseDuration(duration)
    if err != nil {
		return c.RenderError(err)
    }
	revel.TRACE.Printf("DataTsv graphDuration=%s", graphDuration)

	startTime := endTime.Add(-graphDuration)
	revel.TRACE.Printf("DataTsv startTime=%s", startTime)

	if step == 0 {
		step = 60;
		revel.TRACE.Printf("DataTsv step was zero, now %d", step)
	}
	stepDuration := time.Duration(step) * time.Second

    data, err := rrd.Fetch(rrdpath, cf, startTime, endTime, stepDuration)
    if err != nil {
		revel.TRACE.Printf("Fetch fails %s", err)
		return c.RenderError(err)
    }
	defer data.FreeValues()

	return RrdFetchTsvResult{data, step}
}

type RrdFilesResult struct {
	infos []os.FileInfo
}

func (r RrdFilesResult) Apply(req *revel.Request, resp *revel.Response) {
	revel.TRACE.Printf("RrdFilesResult Apply start\n")
	resp.WriteHeader(http.StatusOK, "text/plain")

	for _, fi := range r.infos {
		line := fmt.Sprintf("%s\n", fi.Name())
		resp.Out.Write([]byte(line))
	}
	revel.TRACE.Printf("RrdFilesResult Apply end\n")
}

func (c Api) Files() revel.Result {
	revel.TRACE.Printf("Api.Files start\n")
	defer revel.TRACE.Printf("Api.Files end\n")

	topDir, found := revel.Config.String("rrd.data_dir")
	if !found {
		return c.RenderError(errors.New("config rrd.data_dir not found"))
	}

	infos, err := dir.FindRrdFiles(topDir)
	if err != nil {
		return c.RenderError(err)
	}

	return RrdFilesResult{infos}
}
