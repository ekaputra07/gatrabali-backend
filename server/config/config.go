package config

import "os"

// ServicePort ...
var ServicePort = os.Getenv("SERVICE_PORT")

// MinifluxHost ...
var MinifluxHost = os.Getenv("MINIFLUX_HOST")

// MinifluxUser ...
var MinifluxUser = os.Getenv("MINIFLUX_USER")

// MinifluxPass ...
var MinifluxPass = os.Getenv("MINIFLUX_PASS")

// GCPProject ...
var GCPProject = os.Getenv("GCP_PROJECT")

// PushNotificationTopic ...
var PushNotificationTopic = os.Getenv("PUSH_NOTIFICATION_TOPIC")

// PubSubAPIKey ...
var PubSubAPIKey = os.Getenv("PUBSUB_API_KEY")

func init() {
	if ServicePort == "" {
		ServicePort = "8080"
	}
}
