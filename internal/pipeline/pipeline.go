package pipeline

import (
	"context"
	"fmt"
	"github.com/ermos/cli/spinner"
	"github.com/huetify/back/internal/dbm"
	"github.com/huetify/back/internal/pipes/bridge"
	"os"
	"time"
)

type Piper interface {
	fmt.Stringer
	Run(ctx context.Context, db *dbm.Instance) error
}

func Run (ctx context.Context) error {
	db, err := dbm.Conn(
		ctx,
		os.Getenv("HUETIFY_DB_NAME"),
		os.Getenv("HUETIFY_DB_DRIVER"),
		os.Getenv("HUETIFY_DB_USER"),
		os.Getenv("HUETIFY_DB_PASSWORD"),
		os.Getenv("HUETIFY_DB_HOST"),
		os.Getenv("HUETIFY_DB_PORT"),
	)
	if err != nil {
		return err
	}

	s := spinner.Init(spinner.SimpleLoading, time.Millisecond * 100)
	s.Start("Start pipeline")
	for _, pipe := range pipelines {
		s.Write(pipe.String())
		err := pipe.Run(ctx, db)
		if err != nil {
			s.Stop(true, "Pipeline process failed (%.2fs)", s.TimeElapsed().Seconds())
			return err
		}
	}
	s.Stop(false, "Pipeline process completed successfully (%.2fs)", s.TimeElapsed().Seconds())
	return nil
}

// Default Pipelines
var defaultPipes = []Piper{
	bridge.Pipe{},
}

// Pipelines
var pipelines = append(
	defaultPipes,
)