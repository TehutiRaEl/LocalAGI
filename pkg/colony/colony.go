package colony

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	fiber "github.com/gofiber/fiber/v2"
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

// RegisterFiberRoutes mounts the Colony Standard Layer onto the Fiber app.
func RegisterFiberRoutes(app *fiber.App) {
	app.Get("/colony/info", handleInfo)
	app.Get("/colony/health", handleHealth)
	app.Get("/colony/agents", handleAgents)
	app.Post("/colony/events", handleEvents)
	app.Get("/colony/manifest", handleManifest)
}

func handleInfo(c *fiber.Ctx) error {
	return c.JSON(identity)
}

func handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"colony_id":      identity.ColonyID,
		"status":         "healthy",
		"uptime_seconds": int(time.Since(startTime).Seconds()),
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
	})
}

func handleAgents(c *fiber.Ctx) error {
	agents := make([]fiber.Map, len(identity.Agents))
	for i, a := range identity.Agents {
		agents[i] = fiber.Map{
			"id":           a,
			"status":       "active",
			"capabilities": identity.Capabilities,
		}
	}
	return c.JSON(fiber.Map{
		"colony_id": identity.ColonyID,
		"agents":    agents,
	})
}

func handleEvents(c *fiber.Ctx) error {
	var body map[string]interface{}
	if err := json.Unmarshal(c.Body(), &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "bad request"})
	}
	body["ts"] = time.Now().UTC().Format(time.RFC3339)
	events = append(events, body)
	if len(events) > 100 {
		events = events[1:]
	}
	return c.JSON(fiber.Map{"status": "accepted"})
}

func handleManifest(c *fiber.Ctx) error {
	var soulContent string
	if data, err := os.ReadFile("soul.md"); err == nil {
		soulContent = string(data)
	}
	return c.JSON(fiber.Map{
		"colony":    identity,
		"soul_hash": hashString(soulContent),
		"endpoints": fiber.Map{
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
