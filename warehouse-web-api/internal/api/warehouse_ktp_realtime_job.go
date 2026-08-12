package api

import (
	"encoding/json"
	"fmt"
	"time"

	"web-api/internal/api/services"
	"web-api/internal/realtime"
)

const warehouseKTPRealtimeInterval = 5 * time.Minute

type warehouseKTPRealtimeEvent struct {
	Type string `json:"type"`
	Time string `json:"time"`
}

func StartWarehouseKTPRealtimeJob() {
	ticker := time.NewTicker(warehouseKTPRealtimeInterval)
	defer ticker.Stop()

	lastSignature, err := getWarehouseKTPSignature()
	if err != nil {
		fmt.Println("Warehouse KTP realtime baseline failed:", err)
	}

	for range ticker.C {
		signature, err := getWarehouseKTPSignature()
		if err != nil {
			fmt.Println("Warehouse KTP realtime check failed:", err)
			continue
		}
		if lastSignature == "" {
			lastSignature = signature
			continue
		}
		if signature == lastSignature {
			continue
		}

		lastSignature = signature
		broadcastWarehouseKTPUpdated()
	}
}

func getWarehouseKTPSignature() (string, error) {
	return services.WarehouseKTP.GetRealtimeSignature()
}

func broadcastWarehouseKTPUpdated() {
	payload, err := json.Marshal(warehouseKTPRealtimeEvent{
		Type: "WAREHOUSE_KTP_RACKS_UPDATED",
		Time: time.Now().Format(time.RFC3339),
	})
	if err != nil {
		fmt.Println("Warehouse KTP realtime marshal failed:", err)
		return
	}

	realtime.AlertHub.Broadcast(payload)
}
