package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// API struct
type API struct {
}

//GetAllAudits call
func (a *API) GetAll(w http.ResponseWriter, r *http.Request) {
	g := getAllAggregation()
	b, e := json.Marshal(g)
	logger.Panice(e, "Unable to marshal response : %s", e)
	writeJsonResponse(w, string(b))
}

func getAllAggregation() *Aggregation {
	a := &Aggregation{}
	e := filepath.Walk(config.AuditDir, ParseAllAudits(a))
	logger.Panice(e, "An error occurred walking \"%s\" : %s", config.AuditDir, e)
	return a
}

//ParseAllAudits function
func ParseAllAudits(a *Aggregation) filepath.WalkFunc {
	return func(p string, i os.FileInfo, e error) error {
		if i.IsDir() {
			return nil
		}
		if filepath.Ext(p) == ".json" {
			f := parseAuditFile(p)
			f.Filename = strings.Replace(i.Name(), ".json", "", 1)
			a.Files = append(a.Files, f)
			a.SummaryTests.Count += f.Report.Summary.Count
			a.SummaryTests.Failed += f.Report.Summary.Failed
			a.SummaryTests.Passed += f.Report.Summary.Passed
			a.SummaryTests.Skipped += f.Report.Summary.Skipped
			a.SummaryTests.Duration += f.Report.Summary.Duration
			a.SummaryHosts.Count++
			switch {
			case f.Report.Summary.Failed == 0:
				a.SummaryHosts.Passed++
			case f.Report.Summary.Failed <= config.Tolerance:
				a.SummaryHosts.Skipped++
			default:
				a.SummaryHosts.Failed++
			}
		}
		return nil
	}
}

func parseAuditFile(p string) *AuditFile {
	r := &AuditFile{}
	f, e := os.Open(p)
	if e != nil {
		logger.Errorf("Unable to open file %s : %s", p, e)
	} else {
		e = json.NewDecoder(f).Decode(r)
		if e != nil {
			logger.Errorf("Unable to parse file %s : %s", p, e)
		}
	}
	return r
}

func writeJsonResponse(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, s)
}

//GetInventory entrypoint
func (a *API) GetInventory(w http.ResponseWriter, r *http.Request) {
	g := getAllAggregation()
	i := &Inventory{SummaryTests: g.SummaryTests, SummaryHosts: g.SummaryHosts}
	for _, f := range g.Files {
		j := convertFileToItem(f)
		switch {
		case f.Report.Summary.Failed == 0:
			j.Status = Passed
		case f.Report.Summary.Failed <= config.Tolerance:
			j.Status = Warning
		default:
			j.Status = Failed
		}
		i.Items = append(i.Items, j)
	}
	b, e := json.Marshal(i)
	logger.Panice(e, "Unable to marshal response : %s", e)
	writeJsonResponse(w, string(b))
}

func convertFileToItem(f *AuditFile) *InventoryItem {
	p := regexp.MustCompile(`test_(.+).py::test_(.+)\[.*\]$`)
	j := &InventoryItem{}
	j.Hostname = f.Filename
	j.Passed = f.Report.Summary.Passed
	j.Failed = f.Report.Summary.Failed
	j.Skipped = f.Report.Summary.Skipped
	j.Date = f.Report.CreatedAt
	for _, t := range f.Report.Tests {
		m := p.FindAllStringSubmatch(t.RawName, -1)
		t.Group = m[0][1]
		t.Name = m[0][2]
		j.Tests = append(j.Tests, t)
		if t.Group == "info" {
			if t.Name == "type" {
				j.OS = strings.Replace(t.Call.Output, "\n", " ", -1)
			} else if t.Name == "distrib" {
				std := strings.Split(strings.Replace(t.Call.Output, "None", "", -1), "\n")
				j.Distribution = std[0]
				j.Version = std[len(std)-2]
			}
		}
	}
	return j
}

func (a *API) GetHost(w http.ResponseWriter, r *http.Request) {
	p, ok := r.URL.Query()["h"]
	if !ok || len(p[0]) < 1 {
		logger.Error("Url Param 'h' is missing")
		return
	}
	h := p[0]
	f := filepath.Join(config.AuditDir, fmt.Sprintf("%s%s", h, ".json"))
	g := parseAuditFile(f)
	i := convertFileToItem(g)
	b, e := json.Marshal(i)
	logger.Panice(e, "Unable to marshal response : %s", e)
	writeJsonResponse(w, string(b))
}
