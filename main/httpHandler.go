package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"grafanareports/genReports"
	"grafanareports/gfClient"
)

// ServeReportHandler interface facilitates testsing the reportServing http handler
type ServeReportHandler struct {
	newGrafanaClient func(url string, apiToken string, variables url.Values, sslCheck bool, gridLayout bool) gfClient.Client
	newReport        func(g gfClient.Client, dashName string, time gfClient.TimeRange, texTemplate string, gridLayout bool) genReports.Report
}

// RegisterHandlers registers all http.Handler's with their associated routes to the router
// Two different serve report handlers are used to provide support for both Grafana v4 (and older) and v5 APIs
func RegisterHandlers(router *mux.Router, reportServerV5 ServeReportHandler) {
	router.Handle("/api/v5/report/{dashId}", reportServerV5)
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "This is grafana-reporter. \nThe API endpoints are documented here: https://github.com/IzakMarais/reporter#endpoint.")
	})

}

func (h ServeReportHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Print("Reporter called")
	g := h.newGrafanaClient(*proto+*ip, apiToken(req), dashVariables(req), *sslCheck, *gridLayout)
	rep := h.newReport(g, dashID(req), time(req), texTemplate(req), *gridLayout)

	file, err := rep.Generate()
	if err != nil {
		log.Println("Error generating report:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	defer rep.Clean()
	defer file.Close()
	addFilenameHeader(w, rep.Title())

	_, err = io.Copy(w, file)
	if err != nil {
		log.Println("Error copying data to response:", err)
		http.Error(w, err.Error(), 500)
		return
	}
	log.Println("Report generated correctly")
}

func time(r *http.Request) gfClient.TimeRange {
	params := r.URL.Query()
	t := gfClient.NewTimeRange(params.Get("from"), params.Get("to"))
	log.Println("Called with time range:", t)
	return t
}

func apiToken(r *http.Request) string {
	apiToken := r.URL.Query().Get("apitoken")
	log.Println("Called with api Token:", apiToken)
	return apiToken
}

func dashVariables(r *http.Request) url.Values {
	output := url.Values{}
	for k, v := range r.URL.Query() {
		if strings.HasPrefix(k, "var-") {
			log.Println("Called with variable:", k, v)
			for _, singleV := range v {
				output.Add(k, singleV)
			}
		}
	}
	if len(output) == 0 {
		log.Println("Called without variable")
	}
	return output
}

func dashID(r *http.Request) string {
	vars := mux.Vars(r)
	d := vars["dashId"]
	log.Println("Called with dashboard:", d)
	return d
}

func texTemplate(r *http.Request) string {
	fName := r.URL.Query().Get("template")
	if fName == "" {
		return ""
	}
	file := filepath.Join(*templateDir, fName+".tex")
	log.Println("Called with template:", file)

	customTemplate, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Error reading template file: %q", err)
		return ""
	}

	return string(customTemplate)
}

func addFilenameHeader(w http.ResponseWriter, title string) {
	//sanitize title. Http headers should be ASCII
	filename := strconv.QuoteToASCII(title)
	filename = strings.TrimLeft(filename, "\"")
	filename = strings.TrimRight(filename, "\"")
	filename += ".pdf"
	log.Println("Extracted filename from dashboard title: ", filename)
	header := fmt.Sprintf("inline; filename=\"%s\"", filename)
	w.Header().Add("Content-Disposition", header)
}
