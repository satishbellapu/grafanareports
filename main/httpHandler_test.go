package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"grafanareports/genReports"
	"grafanareports/gfClient"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type mockReport struct {
}

func (m mockReport) Generate() (pdf io.ReadCloser, err error) {
	return ioutil.NopCloser(bytes.NewReader(nil)), nil
}

func (m mockReport) Title() string {
	return "title"
}

func (m mockReport) Clean() {
}

func TestV5ServeReportHandler(t *testing.T) {
	Convey("When the v5 report server handler is called", t, func() {
		//mock new grafana client function to capture and validate its input parameters
		var clAPIToken string
		var clVars url.Values
		newGrafanaClient := func(url string, apiToken string, variables url.Values, sslCheck bool, gridLayout bool) gfClient.Client {
			clAPIToken = apiToken
			clVars = variables
			return gfClient.NewV5Client(url, apiToken, variables, true, false)
		}
		//mock new report function to capture and validate its input parameters
		var repDashName string
		newReport := func(g gfClient.Client, dashName string, _ gfClient.TimeRange, _ string, _ bool) genReports.Report {
			repDashName = dashName
			return &mockReport{}
		}

		router := mux.NewRouter()
		RegisterHandlers(router, ServeReportHandler{newGrafanaClient, newReport})
		rec := httptest.NewRecorder()

		Convey("It should extract dashboard ID from the URL and forward it to the new reporter ", func() {
			req, _ := http.NewRequest("GET", "/api/v5/report/testDash", nil)
			router.ServeHTTP(rec, req)
			So(repDashName, ShouldEqual, "testDash")
		})

		Convey("It should extract the apiToken from the URL and forward it to the new Grafana Client ", func() {
			req, _ := http.NewRequest("GET", "/api/v5/report/testDash?apitoken=1234", nil)
			router.ServeHTTP(rec, req)
			So(clAPIToken, ShouldEqual, "1234")
		})

		Convey("It should extract the grafana variables and forward them to the new Grafana Client ", func() {
			req, _ := http.NewRequest("GET", "/api/v5/report/testDash?var-test=testValue", nil)
			router.ServeHTTP(rec, req)
			expected := url.Values{}
			expected.Add("var-test", "testValue")
			So(clVars, ShouldResemble, expected)

			Convey("Variables should not contain other query parameters ", func() {
				req, _ := http.NewRequest("GET", "/api/v5/report/testDash?var-test=testValue&apitoken=1234", nil)
				router.ServeHTTP(rec, req)
				expected := url.Values{}
				expected.Add("var-test", "testValue") //apitoken not expected here
				So(clVars, ShouldResemble, expected)
			})
		})
	})
}
