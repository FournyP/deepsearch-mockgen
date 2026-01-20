package tui

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

func RunTextInputModel(m textInputModel) (textInputModel, error) {
	p := tea.NewProgram(m)
	final, err := p.Run()
	if err != nil {
		return textInputModel{}, err
	}
	if fm, ok := final.(textInputModel); ok {
		if fm.Cancelled {
			return textInputModel{}, errors.New("cancelled")
		}
		return fm, nil
	}
	return textInputModel{}, nil
}
