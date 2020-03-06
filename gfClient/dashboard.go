package gfClient

import (
	"encoding/json"
	"log"
	"net/url"
	"strings"
)

// Panel represents a Grafana dashboard panel position
type GridPos struct {
	H float64 `json:"h"`
	W float64 `json:"w"`
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Panel represents a Grafana dashboard panel
type Panel struct {
	Id      int
	Type    string
	Title   string
	GridPos GridPos
}

func (p PanelType) string() string {
	return [...]string{
		"singlestat",
		"text",
		"graph",
		"table",
	}[p]
}

type PanelType int

const (
	SingleStat PanelType = iota
	Text
	Graph
	Table
)

func (p Panel) Is(t PanelType) bool {
	if p.Type == t.string() {
		return true
	}
	return false
}

// Row represents a container for Panels
type Row struct {
	Id        int
	Showtitle bool
	Title     string
	Panels    []Panel
}

// Dashboard represents a Grafana dashboard
// This is both used to unmarshal the dashbaord JSON into
// and then enriched (sanitize fields for TeX consumption and add VarialbeValues)
type Dashboard struct {
	Title          string
	Description    string
	VariableValues string //Not present in the Grafana JSON structure. Enriched data passed used by the Tex templating
	Rows           []Row
	Panels         []Panel
}

type dashContainer struct {
	Dashboard Dashboard
	Meta      struct {
		Slug string
	}
}

// NewDashboard creates Dashboard from Grafana's internal JSON dashboard definition
func NewDashboard(dashJSON []byte, variables url.Values) Dashboard {
	var dash dashContainer
	err := json.Unmarshal(dashJSON, &dash)
	if err != nil {
		panic(err)
	}
	d := dash.NewDashboard(variables)
	log.Printf("Populated dashboard datastructure: %+v\n", d)
	return d
}

func (dc dashContainer) NewDashboard(variables url.Values) Dashboard {
	var dash Dashboard
	dash.Title = sanitizeLaTexInput(dc.Dashboard.Title)
	dash.Description = sanitizeLaTexInput(dc.Dashboard.Description)
	dash.VariableValues = sanitizeLaTexInput(getVariablesValues(variables))

	return populatePanelsFromV5JSON(dash, dc)
}

func populatePanelsFromV5JSON(dash Dashboard, dc dashContainer) Dashboard {
	for _, p := range dc.Dashboard.Panels {
		if p.Type == "row" {
			continue
		}
		p.Title = sanitizeLaTexInput(p.Title)
		dash.Panels = append(dash.Panels, p)
	}
	return dash
}

func getVariablesValues(variables url.Values) string {
	values := []string{}
	for _, v := range variables {
		values = append(values, strings.Join(v, ", "))
	}
	return strings.Join(values, ", ")
}

func sanitizeLaTexInput(input string) string {
	input = strings.Replace(input, "\\", "\\textbackslash ", -1)
	input = strings.Replace(input, "&", "\\&", -1)
	input = strings.Replace(input, "%", "\\%", -1)
	input = strings.Replace(input, "$", "\\$", -1)
	input = strings.Replace(input, "#", "\\#", -1)
	input = strings.Replace(input, "_", "\\_", -1)
	input = strings.Replace(input, "{", "\\{", -1)
	input = strings.Replace(input, "}", "\\}", -1)
	input = strings.Replace(input, "~", "\\textasciitilde ", -1)
	input = strings.Replace(input, "^", "\\textasciicircum ", -1)
	return input
}
