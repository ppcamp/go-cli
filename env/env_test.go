package env_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ppcamp/go-lib/env"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	assert := require.New(t)

	expectedInt := 3
	expectedStr := "some_a"
	expectedFloat := 1.31
	expectedDef := "some_value"

	// update this environ
	envVars := map[string]string{
		"SOME_INT":    fmt.Sprintf("%d", expectedInt),
		"SOME_FLOAT":  fmt.Sprintf("%f", expectedFloat),
		"SOME_STRING": expectedStr,
	}
	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			t.Fail()
		}
	}

	var resultInt int
	var resultFloat float64
	var resultStr, resultDef string

	// try to get those parses
	flags := []env.Flag{
		&env.BaseFlag[int]{Value: &resultInt, EnvName: "SOME_INT"},
		&env.BaseFlag[float64]{Value: &resultFloat, EnvName: "SOME_FLOAT"},
		&env.BaseFlag[string]{Value: &resultStr, EnvName: "SOME_STRING"},
		&env.BaseFlag[string]{Value: &resultDef, EnvName: "NOT_IN_ENV", Default: expectedDef},
	}

	// check if occurred some error during parse
	err := env.Parse(flags)
	assert.Nil(err)

	// check if parsed successfully
	expected := []any{expectedInt, expectedFloat, expectedStr, expectedDef}
	results := []any{resultInt, resultFloat, resultStr, resultDef}
	assert.ElementsMatch(expected, results)

	// check the error scenario

	flags = []env.Flag{
		&env.BaseFlag[string]{Value: &resultDef, EnvName: "NOT_IN_ENV", Required: true},
	}

	err = env.Parse(flags)
	assert.NotNil(err)
	assert.ErrorIs(err, env.ErrFlagRequired)
}
