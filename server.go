package main

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// StartListening for incomming request
func StartListening() {
	http.HandleFunc("/health", GenerateHandler("^/health$", HealthHandler))
	http.HandleFunc("/static/", GenerateHandler("^/(static/(js/|css/|media/)[a-zA-Z0-9._]*)$", FileHandler))
	http.HandleFunc("/audits/", GenerateHandler("^/(static/[a-zA-Z0-9._-]*)$", FileHandler))
	http.HandleFunc("/api/", GenerateHandler("^/api/(get/(all|inventory))$", APIHandler))
	http.HandleFunc("/", GenerateHandler("^/(.*)$", FileHandler))
	a := fmt.Sprintf("%s:%s", config.Host, config.Port)
	logger.Infof("Start listening \"%s\"...", a)
	logger.Fatale(http.ListenAndServe(a, nil), "Server crashed !")
}

// GenerateHandler handler
func GenerateHandler(p string, f func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Tracef("Catch request %s %s", p, reflect.TypeOf(f).Name())
		v := regexp.MustCompile(p)
		m := v.FindStringSubmatch(r.URL.Path)
		if m == nil {
			logger.Warningf("Invalid path \"%s\" doesn't match pattern \"%s\"", r.URL.Path, p)
			DefaultHandler(w, r, "")
			return
		}
		logger.Tracef("Pattern %s matched : %q", p, m)
		defer func(w http.ResponseWriter, r *http.Request) {
			if e := recover(); e != nil {
				logger.Recoverf("Recover from handling request : %s", e)
				DefaultHandler(w, r, "")
			}
		}(w, r)
		if len(m) > 1 {
			f(w, r, m[1])
		} else {
			f(w, r, m[0])
		}
	}
}

// HealthHandler handler
func HealthHandler(w http.ResponseWriter, r *http.Request, n string) {
	SetNoCacheHeaders(w)
	logger.Trace("Health check !")
	fmt.Fprint(w, "OK")
}

// FileHandler handler
func FileHandler(w http.ResponseWriter, r *http.Request, n string) {
	switch r.Method {
	case "GET":
		m := fmt.Sprintf("%s/%s", config.WorkingDir, n)
		if _, e := os.Stat(m); os.IsNotExist(e) {
			logger.Debugf("File not found \"%s\"", m)
			logger.Debug("Serve index instead")
			DefaultHandler(w, r, n)
		} else {
			logger.Debugf("Serve file \"%s\"", m)
			http.ServeFile(w, r, m)
		}
	default:
		logger.Warningf("Wrong request method \"%s\" !", r.Method)
		DefaultHandler(w, r, n)
	}
}

// DefaultHandler handler
func DefaultHandler(w http.ResponseWriter, r *http.Request, n string) {
	logger.Debugf("Request Default received : %s", n)
	switch r.Method {
	case "GET":
		SetNoCacheHeaders(w)
		m := fmt.Sprintf("%s/index.html", config.WorkingDir)
		http.ServeFile(w, r, m)
	default:
		http.Error(w, "", http.StatusBadRequest)
	}
}

// APIHandler handler
func APIHandler(w http.ResponseWriter, r *http.Request, n string) {
	SetNoCacheHeaders(w)
	logger.Debugf("Request API received : %s", n)
	m := FormatEndpointMethod(n)
	a := &API{}
	ExecEndpoint(a, m, w, r)
}

// FormatEndpointMethod to get method name from url
func FormatEndpointMethod(n string) string {
	p := strings.Split(n, "/")
	m := ""
	for _, s := range p {
		m = fmt.Sprintf("%s%s", m, strings.Title(s))
	}
	return m
}

// ExecEndpoint to run method from url
func ExecEndpoint(i interface{}, m string, w http.ResponseWriter, r *http.Request) {
	a := reflect.ValueOf(i)
	logger.Debugf("Method call : %s", strings.Title(m))
	f := a.MethodByName(strings.Title(m))
	if f.IsZero() {
		DefaultHandler(w, r, m)
	} else {
		q := []reflect.Value{
			reflect.ValueOf(w),
			reflect.ValueOf(r),
		}
		f.Call(q)
	}
}

// SetNoCacheHeaders to prevent browser caching
func SetNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}
