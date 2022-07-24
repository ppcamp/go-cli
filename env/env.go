package env

import (
	"errors"
	"fmt"
	basestrings "strings"
	"syscall"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/ppcamp/go-xtendlib/interfaces"
)

var (
	ErrFlagRequired error = errors.New("the flag is required")
)

func fromEnv(envVar string) (string, bool) {
	envVar = basestrings.TrimSpace(envVar)
	return syscall.Getenv(envVar)
}

// isEmpty check if the value has the same value as an unitialized variable
func isEmpty[T interfaces.Ordered](value T) bool {
	var r T
	return r == value
}

type Flag interface {
	Apply() error
	fmt.Stringer
	Name() string
	CurrentValue() any
	IsRequired() bool
	DefaultValue() any
}

type Flags []Flag

// Parse the passed flags
func Parse(flags Flags) error {
	for _, v := range flags {
		if err := v.Apply(); err != nil {
			return fmt.Errorf("fail to parse %w", err)
		}
	}
	return nil
}

func (s Flags) String() string {
	t := table.NewWriter()
	t.SetTitle("\nBase Flag\n")
	t.AppendHeader(table.Row{"Var Name", "Default Value", "Required", "Current Value"})
	style := table.StyleLight
	style.Format.Header = text.FormatDefault
	t.SetStyle(style)

	for _, flag := range s {
		t.AppendRows([]table.Row{
			{flag.Name(), flag.DefaultValue(), flag.IsRequired(), flag.CurrentValue()},
		})
	}

	return t.Render()
}
