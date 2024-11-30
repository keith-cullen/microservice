package store

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/keith-cullen/microservice/config"
	"github.com/keith-cullen/microservice/store/ent"
	"github.com/keith-cullen/microservice/store/ent/thing"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	mu     sync.Mutex
	Client *ent.Client
}

const (
	DatabaseDriverName = "sqlite3"
)

func Open() (*Store, error) {
	client, err := ent.Open(DatabaseDriverName, config.Get(config.DatabaseFileKey))
	if err != nil {
		return nil, fmt.Errorf("failed to open store: %w", err)
	}
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed to open store: %w", err)
	}
	log.Print("store open")
	return &Store{
		Client: client,
	}, nil
}

func (store *Store) Close() {
	store.Client.Close()
	log.Print("store closed")
}

func (store *Store) GetThing(ctx context.Context, name string) (string, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	_, err := store.Client.Thing.
		Query().
		Where(thing.Name(name)).
		Only(ctx)
	if err != nil {
		err = fmt.Errorf("failed to get thing: %w", err)
		log.Print(err)
		return "", err
	}
	log.Printf("got thing: %q", name)
	return name, nil
}

func (store *Store) SetThing(ctx context.Context, name string) error {
	store.mu.Lock()
	defer store.mu.Unlock()
	t, err := store.Client.Thing.
		Query().
		Where(thing.Name(name)).
		Only(ctx)
	switch err.(type) {
	case *ent.NotSingularError:
		err = fmt.Errorf("failed to set thing: %w", err)
		log.Print(err)
	case *ent.NotFoundError:
		err = store.createThing(ctx, name)
	default:
		err = store.updateThing(ctx, t)
	}
	return err
}

// unsafe - store.mu must be locked when this method is called
func (store *Store) createThing(ctx context.Context, name string) error {
	_, err := store.Client.Thing.
		Create().
		SetName(name).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to create thing: %w", err)
		log.Print(err)
		return err
	}
	log.Printf("created thing: %q", name)
	return nil
}

// unsafe - store.mu must be locked when this method is called
func (store *Store) updateThing(ctx context.Context, t *ent.Thing) error {
	_, err := store.Client.Thing.
		UpdateOne(t).
		SetName(t.Name).
		Save(ctx)
	if err != nil {
		err = fmt.Errorf("failed to update thing: %w", err)
		log.Print(err)
		return err
	}
	log.Printf("updated thing: %q", t.Name)
	return nil
}
