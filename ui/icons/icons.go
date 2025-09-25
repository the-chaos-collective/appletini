package icons

import (
	_ "embed"

	"fyne.io/fyne/v2"
)

//go:embed assets/tray1.png
var iconDefault []byte

var ResIconDefault = &fyne.StaticResource{
	StaticName:    "tray1.png",
	StaticContent: iconDefault,
}

//go:embed assets/tray2.png
var iconReviewable []byte

var ResIconReviewable = &fyne.StaticResource{
	StaticName:    "tray2.png",
	StaticContent: iconReviewable,
}

//go:embed assets/tray3.png
var iconMergeable []byte

var ResIconMergeable = &fyne.StaticResource{
	StaticName:    "tray3.png",
	StaticContent: iconMergeable,
}

//go:embed assets/tray4.png
var iconBoth []byte

var ResIconBoth = &fyne.StaticResource{
	StaticName:    "tray4.png",
	StaticContent: iconBoth,
}

//go:embed assets/tray1_dark.png
var iconDefaultDark []byte

var ResIconDefaultDark = &fyne.StaticResource{
	StaticName:    "tray1_dark.png",
	StaticContent: iconDefaultDark,
}

//go:embed assets/tray2_dark.png
var iconReviewableDark []byte

var ResIconReviewableDark = &fyne.StaticResource{
	StaticName:    "tray2_dark.png",
	StaticContent: iconReviewableDark,
}

//go:embed assets/tray3_dark.png
var iconMergeableDark []byte

var ResIconMergeableDark = &fyne.StaticResource{
	StaticName:    "tray3_dark.png",
	StaticContent: iconMergeableDark,
}

//go:embed assets/tray4_dark.png
var iconBothDark []byte

var ResIconBothDark = &fyne.StaticResource{
	StaticName:    "tray4_dark.png",
	StaticContent: iconBothDark,
}

//go:embed assets/warning.png
var iconWarning []byte

var ResIconWarning = &fyne.StaticResource{
	StaticName:    "warning.png",
	StaticContent: iconWarning,
}
