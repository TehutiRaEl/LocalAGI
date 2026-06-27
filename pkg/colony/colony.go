package colony

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"time"

	fiber "github.com/gofiber/fiber/v2"
)

var colonyInfo = fiber.Map{
	"colony_id":   "localagi",
	"colony_name": "LocalAGI",
	"role":        "colony",
	"archetype":   "body",
	"layer":       3,
	"entity":      "BODY (The Swarm)",
	"guilds":      []string{"swarm", "workflow", "agent"},
	"hive":        "sovereign-hive",
	"queen":       "https://github.com/TehutiRaEl/-sovereign-hive-meta",
	"version":     "1.0.0",
}

func RegisterFiberRoutes(app *fiber.App) {
	app.Get("/colony/info", handleInfo)
	app.Get("/colony/health", handleHealth)
	app.Get("/colony/agents", handleAgents)
	app.Post("/colony/events", handleEvents)
	app.Get("/colony/manifest", handleManifest)
}

func handleInfo(c *fiber.Ctx) error {
	return c.JSON(colonyInfo)
}

func handleHealth(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":    "ok",
		"colony_id": "localagi",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

func handleAgents(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"colony_id": "localagi",
		"agents":    []string{"task-executor", "skill-manager", "knowledge-agent", "integration-agent"},
	})
}

func handleEvents(c *fiber.Ctx) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(c.Body(), &payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid JSON"})
	}
	h := fnv.New32a()
	b, _ := json.Marshal(payload)
	h.Write(b)
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"accepted": true,
		"event_id": fmt.Sprintf("%08x", h.Sum32()),
	})
}

func handleManifest(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"colony_id": "localagi",
		"endpoints": fiber.Map{
			"info":     "/colony/info",
			"health":   "/colony/health",
			"agents":   "/colony/agents",
			"events":   "/colony/events",
			"manifest": "/colony/manifest",
		},
		"constitution": "https://raw.githubusercontent.com/TehutiRaEl/-sovereign-hive-meta/main/soul.md",
	})
}
