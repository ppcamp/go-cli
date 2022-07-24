package env

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
	"github.com/ppcamp/go-strings"
	"github.com/ppcamp/go-xtendlib/interfaces"
)

// BaseFlag can be used for the BaseFlagTypes only.
//
// If you don't pass a Default value, the variable will be mandatory
type BaseFlag[T interfaces.Ordered] struct {
	// Value is the address of some variable, which can be some pkg variable, for example
	Value *T

	// Default is the default value to assign to this variable
	Default T

	// EnvName is the name of the environment variable that will try to fetch this data
	EnvName string

	// Required is used to check if the value was found in the path.
	// This is necessary, since that an empty variable can be in the env
	Required bool
}

func (s *BaseFlag[T]) Apply() error {
	valueFromEnv, exist := fromEnv(s.EnvName)
	// check if the flag don't exist and if there's no default value
	if !exist && isEmpty(s.Default) && s.Required {
		return fmt.Errorf("flag %s is not defined: %w", s.EnvName, ErrFlagRequired)
	}

	// creates a pointer of the type T pointing to the response object and switch basing on the ptrs
	var response T
	switch p := any(&response).(type) {
	case *int:
		tmp, err := strings.ToInt[int](valueFromEnv)
		if err != nil {
			return fmt.Errorf("fail to parse flag %s error %w", s.EnvName, err)
		}
		*p = tmp

	case *int32:
		tmp, err := strings.ToInt[int32](valueFromEnv)
		if err != nil {
			return fmt.Errorf("fail to parse flag %s error %w", s.EnvName, err)
		}
		*p = tmp

	case *int64:
		tmp, err := strings.ToInt[int64](valueFromEnv)
		if err != nil {
			return fmt.Errorf("fail to parse flag %s error %w", s.EnvName, err)
		}
		*p = tmp

	case *float32:
		tmp, err := strings.ToFloat[float32](valueFromEnv)
		if err != nil {
			return fmt.Errorf("fail to parse flag %s error %w", s.EnvName, err)
		}
		*p = tmp

	case *float64:
		tmp, err := strings.ToFloat[float64](valueFromEnv)
		if err != nil {
			return fmt.Errorf("fail to parse flag %s error %w", s.EnvName, err)
		}
		*p = tmp

	case *string:
		*p = valueFromEnv

	default:
		return fmt.Errorf("type %T is not supported yet", p)
	}

	// update the value of the passed variable
	if isEmpty(response) {
		*s.Value = s.Default
	} else {
		*s.Value = response
	}

	return nil
}

func (s *BaseFlag[T]) String() string {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetTitle("\nBase Flag\n")
	t.AppendHeader(table.Row{"Var Name", "Default Value", "Required", "Current Value"})
	t.AppendRows([]table.Row{
		{s.EnvName, s.Default, s.Required, *s.Value},
	})

	style := table.StyleLight
	style.Format.Header = text.FormatDefault
	t.SetStyle(style)

	return t.Render()
}

func (s *BaseFlag[T]) Name() string {
	return s.EnvName
}

func (s *BaseFlag[T]) DefaultValue() any {
	return s.Default
}

func (s *BaseFlag[T]) CurrentValue() any {
	return *s.Value
}

func (s *BaseFlag[T]) IsRequired() bool {
	return s.Required
}
