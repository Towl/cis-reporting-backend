package main

//Status string
type Status string

//List of Status
const (
	Passed  Status = "passed"
	Failed  Status = "failed"
	Skipped Status = "skipped"
	Warning Status = "warning"
)

//Aggregation struct
type Aggregation struct {
	Files        []*AuditFile `json:"audit_files"`
	SummaryTests Summary      `json:"summary_tests"`
	SummaryHosts Summary      `json:"summary_hosts"`
}

//AuditFile struct
type AuditFile struct {
	Report   Report `json:"report"`
	Filename string `json:"filename"`
}

//Report struct
type Report struct {
	Environment Environment `json:"environment"`
	Summary     Summary     `json:"summary"`
	Tests       []Test      `json:"tests"`
	CreatedAt   string      `json:"created_at"`
}

//Environment struct
type Environment struct {
	Python   string `json:"Python"`
	Platform string `json:"Platform"`
}

//Test struct
type Test struct {
	RawName  string  `json:"name"`
	Name     string  `json:"display_name"`
	Group    string  `json:"group"`
	Duration float64 `json:"duration"`
	RunIndex int     `json:"run_index"`
	Setup    Process `json:"setup"`
	Call     Process `json:"call"`
	Teardown Process `json:"teardown"`
	Outcome  Status  `json:"outcome"`
}

//Process struct
type Process struct {
	Name     string  `json:"name"`
	Duration float64 `json:"duration"`
	Outcome  Status  `json:"outcome"`
	Message  string  `json:"longrepr"`
	Output   string  `json:"stdout"`
}

//Summary struct
type Summary struct {
	Failed   int     `json:"failed"`
	Passed   int     `json:"passed"`
	Skipped  int     `json:"skipped"`
	Count    int     `json:"num_tests"`
	Duration float64 `json:"duration"`
}

//Inventory struct
type Inventory struct {
	Items        []*InventoryItem `json:"items"`
	SummaryTests Summary          `json:"summary_tests"`
	SummaryHosts Summary          `json:"summary_hosts"`
}

//InventoryItem struct
type InventoryItem struct {
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Distribution string `json:"distribution"`
	Version      string `json:"version"`
	Passed       int    `json:"passed"`
	Skipped      int    `json:"skipped"`
	Failed       int    `json:"failed"`
	Date         string `json:"date"`
	Status       Status `json:"status"`
	Tests        []Test `json:"tests"`
}
