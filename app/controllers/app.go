package controllers

import (
	"errors"
	"strconv"
	"strings"
	"time"
    "github.com/robfig/revel"
	"github.com/hnakamur/goemon/util/dir"
)

type Application struct {
	*revel.Controller
}

func (c Application) Index() revel.Result {
	return c.Render()
}

func (c Application) Chart(path, cf, duration string, end time.Time) revel.Result {
	revel.TRACE.Printf("chart start! end=%s %d\n", end, end.Unix());
	var endUnix string
	if end.IsZero() {
		revel.TRACE.Printf("chart end.IsZero()=true\n")
		endUnix = ""
	} else {
		endUnix = strconv.FormatInt(end.Unix(), 10)
	}
	revel.TRACE.Printf("chart endUnix=%s\n", endUnix)

	durationsStr, found := revel.Config.String("graph.durations")
	if !found {
		return c.RenderError(errors.New("config graph.durations not found"))
	}

	shiftDurationsStr, found := revel.Config.String("graph.shift_durations")
	if !found {
		return c.RenderError(errors.New("config graph.shift_durations not found"))
	}

	durations := strings.Split(durationsStr, ",")
	shiftDurations := strings.Split(shiftDurationsStr, ",")
	return c.Render(path, cf, duration, endUnix, durations, shiftDurations)
}

func (c Application) Files() revel.Result {
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

	return c.Render(infos)
}
