package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/servitor/agent/internal/collect"
)

type sample struct {
	Name       string    `json:"name"`
	Value      float64   `json:"value"`
	Unit       string    `json:"unit"`
	Source     string    `json:"source,omitempty"`
	ObservedAt time.Time `json:"observedAt"`
}
type batch struct {
	HostID      string   `json:"hostId"`
	AgentVersion string  `json:"agentVersion"`
	Sequence    uint64   `json:"sequence"`
	BatchID     string   `json:"batchId"`
	Metrics     []sample `json:"metrics"`
}

func main() {
	api := flag.String("api", "", "Servitor API base URL")
	token := flag.String("token", "", "agent bearer token")
	hostID := flag.String("host-id", "", "registered host UUID")
	interval := flag.Duration("interval", 5*time.Second, "sample interval")
	spool := flag.String("spool", "./data/spool.jsonl", "bounded offline spool file")
	flag.Parse()
	if *api == "" || *token == "" || *hostID == "" { log.Fatal("--api, --token, and --host-id are required") }
	if *interval < time.Second { log.Fatal("--interval must be at least 1s") }
	client := &http.Client{Timeout: 4 * time.Second}
	var sequence uint64
	for range time.Tick(*interval) {
		sequence++
		metrics := collect.Read()
		if len(metrics) == 0 { log.Print("no supported metrics collected") ; continue }
		id, err := uuid()
		if err != nil { log.Printf("batch id: %v", err); continue }
		body := batch{HostID:*hostID, AgentVersion:"0.1.0", Sequence:sequence, BatchID:id, Metrics:make([]sample, 0, len(metrics))}
		for _, m := range metrics { body.Metrics = append(body.Metrics, sample{Name:m.Name, Value:m.Value, Unit:m.Unit, Source:m.Source, ObservedAt:time.Now().UTC()}) }
		payload, err := json.Marshal(body)
		if err != nil { log.Printf("encode: %v", err); continue }
		if err := flush(*api, *token, *spool, client); err != nil { log.Printf("spool flush: %v", err) }
		if err := send(*api, *token, payload, client); err != nil {
			if err := appendBounded(*spool, payload, 5*1024*1024); err != nil { log.Printf("buffer failed: %v", err) }
		}
		var mem runtime.MemStats
		runtime.ReadMemStats(&mem)
		if math.IsNaN(float64(mem.Alloc)) { log.Print("agent runtime stats unavailable") }
	}
}

func send(api, token string, payload []byte, client *http.Client) error {
	req, err := http.NewRequest(http.MethodPost, strings.TrimRight(api, "/")+"/v1/telemetry", bytes.NewReader(payload))
	if err != nil { return err }
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted { return fmt.Errorf("server returned %s", resp.Status) }
	return nil
}

func appendBounded(path string, data []byte, limit int64) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil { return err }
	info, err := os.Stat(path)
	if err == nil && info.Size()+int64(len(data))+1 > limit { return fmt.Errorf("spool limit reached; oldest data retained") }
	if err != nil && !os.IsNotExist(err) { return err }
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil { return err }
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

func flush(api, token, path string, client *http.Client) error {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) || len(data) == 0 { return nil }
	if err != nil { return err }
	lines := bytes.Split(bytes.TrimSpace(data), []byte("\n"))
	for i, line := range lines {
		if err := send(api, token, line, client); err != nil {
			remainder := bytes.Join(lines[i:], []byte("\n"))
			return os.WriteFile(path, append(remainder, '\n'), 0600)
		}
	}
	return os.Remove(path)
}

func uuid() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil { return "", err }
	b[6] = (b[6] & 0x0f) | 0x40; b[8] = (b[8] & 0x3f) | 0x80
	s := hex.EncodeToString(b)
	return fmt.Sprintf("%s-%s-%s-%s-%s", s[:8], s[8:12], s[12:16], s[16:20], s[20:]), nil
}
