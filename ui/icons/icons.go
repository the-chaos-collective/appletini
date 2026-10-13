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

// Error variants: same badges, logo tinted amber while polling is failing.

//go:embed assets/tray1_err.png
var iconDefaultErr []byte

var ResIconDefaultErr = &fyne.StaticResource{
	StaticName:    "tray1_err.png",
	StaticContent: iconDefaultErr,
}

//go:embed assets/tray2_err.png
var iconReviewableErr []byte

var ResIconReviewableErr = &fyne.StaticResource{
	StaticName:    "tray2_err.png",
	StaticContent: iconReviewableErr,
}

//go:embed assets/tray3_err.png
var iconMergeableErr []byte

var ResIconMergeableErr = &fyne.StaticResource{
	StaticName:    "tray3_err.png",
	StaticContent: iconMergeableErr,
}

//go:embed assets/tray4_err.png
var iconBothErr []byte

var ResIconBothErr = &fyne.StaticResource{
	StaticName:    "tray4_err.png",
	StaticContent: iconBothErr,
}

//go:embed assets/tray1_err_dark.png
var iconDefaultErrDark []byte

var ResIconDefaultErrDark = &fyne.StaticResource{
	StaticName:    "tray1_err_dark.png",
	StaticContent: iconDefaultErrDark,
}

//go:embed assets/tray2_err_dark.png
var iconReviewableErrDark []byte

var ResIconReviewableErrDark = &fyne.StaticResource{
	StaticName:    "tray2_err_dark.png",
	StaticContent: iconReviewableErrDark,
}

//go:embed assets/tray3_err_dark.png
var iconMergeableErrDark []byte

var ResIconMergeableErrDark = &fyne.StaticResource{
	StaticName:    "tray3_err_dark.png",
	StaticContent: iconMergeableErrDark,
}

//go:embed assets/tray4_err_dark.png
var iconBothErrDark []byte

var ResIconBothErrDark = &fyne.StaticResource{
	StaticName:    "tray4_err_dark.png",
	StaticContent: iconBothErrDark,
}
