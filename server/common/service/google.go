package service

import (
	"context"
	"sync"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

// Google is a struct that holds all Google's cloud services instance instance
// and has function to easily access GCF features such as Firestore, Messaging, etc.
type Google struct {
	firebaseApp   *firebase.App
	project       string
	firestoreOnce sync.Once
	messagingOnce sync.Once
	pubsubOnce    sync.Once

	Firestore *firestore.Client
	Messaging *messaging.Client
	Pubsub    *pubsub.Client
}

// InitFirestore initialize Firestore client
func (g *Google) InitFirestore(ctx context.Context) error {
	var err error
	g.firestoreOnce.Do(func() {
		g.Firestore, err = g.firebaseApp.Firestore(ctx)
	})
	return err
}

// InitMessaging returns FCM client instance (lazily loaded)
func (g *Google) InitMessaging(ctx context.Context) error {
	var err error
	g.messagingOnce.Do(func() {
		g.Messaging, err = g.firebaseApp.Messaging(ctx)
	})
	return err
}

// InitPubsub initialize PubSub client
func (g *Google) InitPubsub(ctx context.Context) error {
	var err error
	g.pubsubOnce.Do(func() {
		g.Pubsub, err = pubsub.NewClient(ctx, g.project)
	})
	return err
}

// PublishToTopic publish a message to specified topic
func (g *Google) PublishToTopic(ctx context.Context, topic string, msg *pubsub.Message) (string, error) {
	if err := g.InitPubsub(ctx); err != nil {
		return "", err
	}
	t := g.Pubsub.Topic(topic)
	return t.Publish(ctx, msg).Get(ctx)
}

// NewGoogle create new Firebase
func NewGoogle(cxt context.Context, project string) (*Google, error) {
	app, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	return &Google{firebaseApp: app, project: project}, nil
}
