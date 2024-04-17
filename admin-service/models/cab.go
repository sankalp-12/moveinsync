package models

import "time"

// Cab represents the data structure for a cab
type Cab struct {
	Location    GeoJSONPoint `json:"location" bson:"location"`
	Status      string       `json:"status" bson:"status"`
	LastUpdated time.Time    `json:"last_updated" bson:"last_updated"`
}

// GeoJSONPoint represents a GeoJSON Point
type GeoJSONPoint struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates"`
}
