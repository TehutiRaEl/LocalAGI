package colony

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var startTime = time.Now()

var events []map[string]interface{}

type colonyInfo struct {
	ColonyID            string   `json:"colony_id"`
	ColonyName          string   `json:"colony_name"`
	Role                string   `json:"role"`
	Description         string   `json:"description"`
	Hive                string   `json:"hive"`
	Repo                string   `json:"repo"`
	Guilds              []string `json:"guilds"`
	Agents              []string `json:"agents"`
	Capabilities        []string `json:"capabilities"`
	ConstitutionVersion string   `json:"constitution_version"`
}

var identity = colonyInfo{
	ColonyID:            "localagi",
	ColonyName:          "LocalAGI",
	Role:                "inference",
	Description:         "Local AI inference and model management colony",
	Hive:                "sovereign-hive",
	Repo:                "https://github.com/tehutirael/localagi",
	Guilds:              []string{"inference", "models", "embeddings"},
	Agents:              []string{"inference-agent", "model-manager", "embedding-agent"},
	Capabilities:        []string{"local-inference", "model-management", "vector-embeddings"},
	ConstitutionVersion: "1.0.0",
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, identity)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]interface{}{
		"colony_id":      identity.ColonyID,
		"status":         "healthy",
		"uptime_seconds": int(time.Since(startTime).Seconds()),
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	})
}

func AgentsHandler(w http.ResponseWriter, r *http.Request) {
	agents := make([]map[string]interface{}, len(identity.Agents))
	for i, a := range identity.Agents {
		agents[i] = map[string]interface{}{
			"id":           a,
			"status":       "active",
			"capabilities": identity.Capabilities,
		}
	}
	writeJSON(w, map[string]interface{}{
		"colony_id": identity.ColonyID,
		"agents":    agents,
	})
}

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	body["ts"] = time.Now().UTC().Format(time.RFC3339)
	events = append(events, body)
	if len(events) > 100 {
		events = events[1:]
	}
	writeJSON(w, map[string]interface{}{"status": "accepted"})
}

func ManifestHandler(w http.ResponseWriter, r *http.Request) {
	var soulContent string
	if data, err := os.ReadFile("soul.md"); err == nil {
		soulContent = string(data)
	}
	writeJSON(w, map[string]interface{}{
		"colony":    identity,
		"soul_hash": hashString(soulContent),
		"endpoints": map[string]string{
			"info":     "/colony/info",
			"health":   "/colony/health",
			"agents":   "/colony/agents",
			"events":   "/colony/events",
			"manifest": "/colony/manifest",
		},
	})
}

func hashString(s string) string {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return fmt.Sprintf("%08x", h)
}

// RegisterRoutes wires colony endpoints onto an http.ServeMux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/colony/info", InfoHandler)
	mux.HandleFunc("/colony/health", HealthHandler)
	mux.HandleFunc("/colony/agents", AgentsHandler)
	mux.HandleFunc("/colony/events", EventsHandler)
	mux.HandleFunc("/colony/manifest", ManifestHandler)
}
