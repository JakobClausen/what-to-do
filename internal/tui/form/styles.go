package form

import (
	"github.com/charmbracelet/huh"
)

func GetStyledForm(data *CommandForm) *huh.Form {
	form := data.form

	theme := huh.ThemeCharm()

	form = form.WithTheme(theme)
	return form
}
