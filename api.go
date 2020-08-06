package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
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
			a.Summary.Count++
			a.Summary.Failed += f.Report.Summary.Failed
			a.Summary.Passed += f.Report.Summary.Passed
			a.Summary.Skipped += f.Report.Summary.Skipped
			a.Summary.Duration += f.Report.Summary.Duration
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
	i := &Inventory{Summary: g.Summary}
	for _, f := range g.Files {
		j := &InventoryItem{}
		j.Hostname = f.Filename
		j.Passed = f.Report.Summary.Passed
		j.Failed = f.Report.Summary.Failed
		j.Skipped = f.Report.Summary.Skipped
		j.Date = f.Report.CreatedAt
		i.Items = append(i.Items, j)
	}
	b, e := json.Marshal(i)
	logger.Panice(e, "Unable to marshal response : %s", e)
	writeJsonResponse(w, string(b))
}
