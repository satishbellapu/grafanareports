package gfClient

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestV5Dashboard(t *testing.T) {
	Convey("When creating a new dashboard from Grafana v5 dashboard JSON", t, func() {
		const v5DashJSON = `
{"Dashboard":
	{
		"Panels":
			[{"Type":"singlestat", "Id":0},
			{"Type":"graph", "Id":1, "GridPos":{"H":6,"W":24,"X":0,"Y":0}},
			{"Type":"singlestat", "Id":2, "Title":"Panel3Title #"},
			{"Type":"text", "GridPos":{"H":6.5,"W":20.5,"X":0,"Y":0}, "Id":3},
			{"Type":"table", "Id":4},
			{"Type":"row", "Id":5}],
		"Title":"DashTitle #"
	},

"Meta":
	{"Slug":"testDash"}
}`
		dash := NewDashboard([]byte(v5DashJSON), url.Values{})

		Convey("Panel Is(type) should work for all panels", func() {
			So(dash.Panels[0].Is(SingleStat), ShouldBeTrue)
			So(dash.Panels[1].Is(Graph), ShouldBeTrue)
			So(dash.Panels[2].Is(SingleStat), ShouldBeTrue)
			So(dash.Panels[3].Is(Text), ShouldBeTrue)
			So(dash.Panels[4].Is(Table), ShouldBeTrue)
		})

		Convey("Panel titles should be parsed and sanitised", func() {
			So(dash.Panels[2].Title, ShouldEqual, "Panel3Title \\#")
		})

		Convey("Panels should contain all panels that have type != row", func() {
			So(dash.Panels, ShouldHaveLength, 5)
			So(dash.Panels[0].Id, ShouldEqual, 0)
			So(dash.Panels[1].Id, ShouldEqual, 1)
			So(dash.Panels[2].Id, ShouldEqual, 2)
		})

		Convey("The Title should be parsed", func() {
			So(dash.Title, ShouldEqual, "DashTitle \\#")
		})

		Convey("Panels should contain GridPos H & W", func() {
			So(dash.Panels[1].GridPos.H, ShouldEqual, 6)
			So(dash.Panels[1].GridPos.W, ShouldEqual, 24)
		})

		Convey("Panels GridPos should allow floatt", func() {
			So(dash.Panels[3].GridPos.H, ShouldEqual, 6.5)
			So(dash.Panels[3].GridPos.W, ShouldEqual, 20.5)
		})

	})
}

func TestVariableValues(t *testing.T) {
	Convey("When creating a dashboard and passing url varialbes in", t, func() {
		const v5DashJSON = `
{
	"Dashboard":
		{
		}
}`
		vars := url.Values{}
		vars.Add("var-one", "oneval")
		vars.Add("var-two", "twoval")
		dash := NewDashboard([]byte(v5DashJSON), vars)

		Convey("The dashboard should contain the variable values in a random order", func() {
			So(dash.VariableValues, ShouldContainSubstring, "oneval")
			So(dash.VariableValues, ShouldContainSubstring, "twoval")
		})
	})
}
