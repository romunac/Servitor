//go:build linux

package collect

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Metric struct { Name string; Value float64; Unit string; Source string }

func Read() []Metric {
	var metrics []Metric
	if value, err := cpuPercent(); err == nil { metrics = append(metrics, Metric{"cpu_percent", value, "percent", "proc/stat"}) }
	if value, err := memoryPercent(); err == nil { metrics = append(metrics, Metric{"memory_percent", value, "percent", "proc/meminfo"}) }
	metrics = append(metrics, temperatures()...)
	return metrics
}

func cpuPercent() (float64, error) {
	read := func() (uint64, uint64, error) {
		f, err := os.Open("/proc/stat"); if err != nil { return 0, 0, err }; defer f.Close()
		s := bufio.NewScanner(f); if !s.Scan() { return 0, 0, s.Err() }
		fields := strings.Fields(s.Text()); if len(fields) < 5 { return 0, 0, os.ErrInvalid }
		var total, idle uint64
		for i, field := range fields[1:] { n, e := strconv.ParseUint(field,10,64); if e != nil { return 0,0,e }; total += n; if i == 3 || i == 4 { idle += n } }
		return total, idle, nil
	}
	t1, i1, err := read(); if err != nil { return 0, err }; time.Sleep(100*time.Millisecond); t2, i2, err := read(); if err != nil { return 0, err }
	dt, di := float64(t2-t1), float64(i2-i1); if dt <= 0 { return 0, os.ErrInvalid }; return 100*(dt-di)/dt, nil
}

func memoryPercent() (float64,error) {
	f,err:=os.Open("/proc/meminfo"); if err!=nil{return 0,err}; defer f.Close()
	values:=map[string]float64{}; s:=bufio.NewScanner(f)
	for s.Scan(){p:=strings.Fields(s.Text()); if len(p)>=2 {n,e:=strconv.ParseFloat(p[1],64);if e==nil{values[strings.TrimSuffix(p[0],":")]=n}}}
	total,available:=values["MemTotal"],values["MemAvailable"]; if total==0{return 0,os.ErrInvalid}; return 100*(total-available)/total,nil
}

func temperatures() []Metric {
	var out []Metric
	paths,_:=filepath.Glob("/sys/class/hwmon/hwmon*/temp*_input")
	zones,_:=filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	paths=append(paths,zones...)
	for _,path:=range paths {
		b,err:=os.ReadFile(path); if err!=nil{continue}; raw,err:=strconv.ParseFloat(strings.TrimSpace(string(b)),64);if err!=nil{continue}
		c:=raw/1000; if c < -40 || c > 200 {continue}
		out=append(out,Metric{"temperature_c",c,"celsius",filepath.Base(filepath.Dir(path))})
	}
	return out
}
