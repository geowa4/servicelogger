package internalservicelog

import "github.com/geowa4/servicelogger/pkg/teaspoon"

func Program() (string, bool, error) {
	tm, err := teaspoon.Program(initialModel())
	if err != nil {
		return "", false, err
	}
	if m, ok := tm.(*model); ok {
		return m.Markdown(), m.confirmation, nil
	}
	return "", false, err
}
