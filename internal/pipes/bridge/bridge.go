package bridge

import (
	"context"
	"github.com/huetify/back/internal/dbm"
	"github.com/huetify/back/internal/models/bridge"
)

type Pipe struct{}

func (Pipe) String() string {
	return "loading bridge store"
}

func (Pipe) Run(ctx context.Context, db *dbm.Instance) error {
	return bridge.SetBridgesStore(ctx, db)
}