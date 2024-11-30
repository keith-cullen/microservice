package store

import (
	"context"
	"fmt"
	"log"

	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store/ent"
	"github.com/keith-cullen/microservice/store/ent/thing"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	Client *ent.Client
}

func Open() (*Store, error) {
	client, err := ent.Open(
		config.Get(config.StoreDriverNameKey),
		config.Get(config.StoreDatabaseFileKey))
	if err != nil {
		return nil, fmt.Errorf("failed to open store: %w", err)
	}
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed to create schema resources: %w", err)
	}
	log.Print("store open")
	return &Store{Client: client}, nil
}

func (store *Store) Close() {
	store.Client.Close()
	log.Print("store closed")
}

func (store *Store) GetThing(ctx context.Context, name string) (string, error) {
	log.Printf("get thing: %s", name)
	u, err := store.Client.Thing.
		Query().
		Where(thing.Name(name)).
		// 'Only' fails if no thing found, or more than 1 thing returned
		Only(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get thing: %w", err)
	}
	return u.Name, nil
}

func (store *Store) SetThing(ctx context.Context, name string) error {
	log.Printf("set thing: %s", name)
	_, err := store.Client.Thing.
		Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to set thing: %w", err)
	}
	return nil
}
