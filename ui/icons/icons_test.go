package icons_test

import (
	"bytes"
	"testing"

	"appletini/ui/icons"

	"fyne.io/fyne/v2"
)

var pngMagic = []byte("\x89PNG\r\n\x1a\n")

func TestIconsAreEmbeddedPNGs(t *testing.T) {
	testCases := []struct {
		name     string
		resource *fyne.StaticResource
	}{
		{name: "default", resource: icons.ResIconDefault},
		{name: "reviewable", resource: icons.ResIconReviewable},
		{name: "mergeable", resource: icons.ResIconMergeable},
		{name: "both", resource: icons.ResIconBoth},
		{name: "default dark", resource: icons.ResIconDefaultDark},
		{name: "reviewable dark", resource: icons.ResIconReviewableDark},
		{name: "mergeable dark", resource: icons.ResIconMergeableDark},
		{name: "both dark", resource: icons.ResIconBothDark},
		{name: "default error", resource: icons.ResIconDefaultErr},
		{name: "reviewable error", resource: icons.ResIconReviewableErr},
		{name: "mergeable error", resource: icons.ResIconMergeableErr},
		{name: "both error", resource: icons.ResIconBothErr},
		{name: "default error dark", resource: icons.ResIconDefaultErrDark},
		{name: "reviewable error dark", resource: icons.ResIconReviewableErrDark},
		{name: "mergeable error dark", resource: icons.ResIconMergeableErrDark},
		{name: "both error dark", resource: icons.ResIconBothErrDark},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !bytes.HasPrefix(tc.resource.Content(), pngMagic) {
				t.Fatalf("%s is not an embedded PNG (%d bytes)", tc.resource.Name(), len(tc.resource.Content()))
			}
		})
	}
}
