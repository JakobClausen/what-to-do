package tui

import (
	column "what-to-do/internal/tui/column"
)

func InitialModel(width int) column.ColumnModel {

	return column.CreateColumnModel(width)
}
