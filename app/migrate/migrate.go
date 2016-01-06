package migrate

import (
	"fmt"
	"github.com/codegangsta/cli"
	"gopkg.in/inconshreveable/log15.v2"
	"strings"
)

var (
	registeredTask map[string]Task

	ErrMigrateArgsError = fmt.Errorf("migrate args need be 'pugo migrate <type> <source>'")
	ErrMigrateUnknown   = fmt.Errorf("migrate type unknown")
)

func init() {
	registeredTask = map[string]Task{
		TypeRSS: new(RSSTask),
	}
}

type (
	Task interface {
		Is(conf string) bool
		New(ctx *cli.Context) (Task, error)
		Type() string
		Dir() string
		Do() error
	}
)

func Detect(ctx *cli.Context) (Task, error) {
	src := ctx.String("src")
	if len(src) == 0 || !strings.Contains(src, "://") {
		return nil, ErrMigrateArgsError
	}
	for _, task := range registeredTask {
		if task.Is(src) {
			log15.Info("Migrate.Detect.[" + task.Type() + "]")
			return task.New(ctx)
		}
	}
	return nil, ErrMigrateUnknown
}

// Register new migrate task
func Register(task Task) {
	registeredTask[task.Type()] = task
}
