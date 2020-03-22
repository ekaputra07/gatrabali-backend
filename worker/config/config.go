package config

import "os"

// ServicePort ...
var ServicePort = os.Getenv("PORT")

// MinifluxHost ...
var MinifluxHost = os.Getenv("MINIFLUX_HOST")

// MinifluxUser ...
var MinifluxUser = os.Getenv("MINIFLUX_USER")

// MinifluxPass ...
var MinifluxPass = os.Getenv("MINIFLUX_PASS")

// GCPProject ...
var GCPProject = os.Getenv("GCP_PROJECT")

func init() {
	if ServicePort == "" {
		ServicePort = "8080"
	}
}
