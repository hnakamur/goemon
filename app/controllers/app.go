package controllers

import "github.com/robfig/revel"

type Application struct {
	*revel.Controller
}

func (c Application) Index() revel.Result {
	return c.Render()
}

func (c Application) Chart() revel.Result {
	query := c.Request.URL.RawQuery
	return c.Render(query)
}
