package event

import (
	"context"
	"fmt"

	"time"

	"github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
)

const eventType = "cache_set_event"

type Client struct {
	sender *nats.Sender
	events cloudevents.Client
}

type Data struct {
	Key string `json:"key"`
}

func New(addr, subject string) (*Client, error) {
	// create cloudevents nats sender
	// other protocol implementations: https://github.com/cloudevents/sdk-go/tree/main/protocol
	sender, err := nats.NewSender(addr, subject, nats.NatsOptions())
	if err != nil {
		return nil, err
	}

	// create cloudevents client
	eventsClient, err := cloudevents.NewClient(sender)
	if err != nil {
		return nil, err
	}

	return &Client{
		sender: sender,
		events: eventsClient,
	}, nil
}

func (c *Client) Send(ctx context.Context, key string) error {
	e, err := newEvent(key)
	if err != nil {
		return err
	}

	res := c.events.Send(ctx, *e)
	if cloudevents.IsUndelivered(res) {
		return fmt.Errorf("failed to send event for key: %s, reason: %v", key, res)
	}

	return nil
}

func (c *Client) CLose(ctx context.Context) error {
	return c.sender.Close(ctx)
}

func newEvent(key string) (*event.Event, error) {
	e := cloudevents.NewEvent()
	e.SetID(uuid.NewString()) // required field
	e.SetSource("cache")      // required field
	e.SetType(eventType)      // required field
	e.SetTime(time.Now())

	err := e.SetData(event.ApplicationJSON, &Data{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return &e, nil
}
