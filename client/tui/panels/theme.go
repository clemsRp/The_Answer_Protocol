package panel

import (
	"github.com/gdamore/tcell/v2"
)

type Theme struct {
	// Global colors
	Background    tcell.Color
	TextPrimary   tcell.Color
	TextSecondary tcell.Color
	TextHighlight tcell.Color

	// Border and titles
	BorderInactive tcell.Color
	BorderActive   tcell.Color
	TitleInactive  tcell.Color
	TitleActive    tcell.Color

	// Lists
	ListSelectedBg  tcell.Color
	ListSelectedTxt tcell.Color

	// Popups and Buttons
	PopupBackgroundHexa string
	PopupBackground     tcell.Color
	ButtonRestBg        tcell.Color
	ButtonActiveBg      tcell.Color
	ButtonText          tcell.Color
}

var AppTheme = Theme{
	Background:    tcell.NewRGBColor(0, 0, 0),
	TextPrimary:   tcell.ColorWhite,
	TextSecondary: tcell.ColorGray,
	TextHighlight: tcell.ColorYellow,

	BorderInactive: tcell.ColorDimGray,
	BorderActive:   tcell.ColorYellow,
	TitleInactive:  tcell.ColorWhite,
	TitleActive:    tcell.ColorYellow,

	ListSelectedBg:  tcell.ColorBlue,
	ListSelectedTxt: tcell.ColorWhite,

	PopupBackgroundHexa: "#3a3838",
	PopupBackground:     tcell.GetColor("#3a3838"),
	ButtonRestBg:        tcell.GetColor("#474646"),
	ButtonActiveBg:      tcell.GetColor("#7e7979"),
	ButtonText:          tcell.ColorWhite,
}
