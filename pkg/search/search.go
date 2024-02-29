package search

import (
	"fmt"
	"github.com/geowa4/servicelogger/pkg/teaspoon"
	"github.com/geowa4/servicelogger/pkg/templates"
)

func Program() (*templates.Template, error) {
	tm, err := teaspoon.Program(NewModel())
	if err != nil {
		return nil, err
	}
	m, ok := tm.(*Model)
	if !ok {
		return nil, fmt.Errorf("received unexpected model type from program: %v\n", err)
	}
	return m.templateSelection.ToTemplate(), nil
}
