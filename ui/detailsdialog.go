//----------------------------------------------------------------------------------------------------------------------
// (w) 2025 by Jan Buchholz
// Details dialog for movies & tv shows, using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"EmbyExplorer_for_Windows/assets"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var detailsDialog *walk.Dialog
var coverView *walk.ImageView
var txtView *walk.TextEdit
var detailsDialogIsOpen = false

const (
	wndWidth       = 560
	wndHeight      = 200
	coverMaxWidth  = "200"
	coverMaxHeight = "200"
)

func createDetailsDialog() {
	err := Dialog{
		AssignTo:        &detailsDialog,
		Title:           assets.CapDetails,
		Icon:            "/assets/app.ico",
		MinSize:         Size{Width: wndWidth, Height: wndHeight},
		MaxSize:         Size{Width: wndWidth, Height: wndHeight},
		Size:            Size{Width: wndWidth, Height: wndHeight},
		Layout:          VBox{},
		DoubleBuffering: true,
		Children: []Widget{
			Composite{Layout: Grid{Columns: 2, MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					ImageView{
						AssignTo:   &coverView,
						MaxSize:    Size{Width: 200, Height: 200},
						Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						Alignment:  AlignHCenterVCenter,
						Mode:       ImageViewModeIdeal,
					},
					TextEdit{
						AssignTo:        &txtView,
						MinSize:         Size{Width: 300, Height: 190},
						Alignment:       AlignHCenterVCenter,
						Background:      SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						DoubleBuffering: true,
						ReadOnly:        true,
						Enabled:         true,
						HScroll:         false,
						VScroll:         true,
					},
				},
			},
		},
	}.Create(mainWindow)
	if err == nil {
		detailsDialogIsOpen = false
		detailsDialog.Closing().Attach(onDetailsDialogClosing)
	}
}

func onDetailsDialogClosing(canceled *bool, _ walk.CloseReason) {
	detailsDialogIsOpen = false
	*canceled = !globalShutdown
	if *canceled {
		detailsDialog.Hide()
	}
}

func setDetails(imagePath string, overview string) {
	_ = txtView.SetText(overview)
	if imagePath != "" {
		image, err := walk.NewImageFromFileForDPI(imagePath, 96)
		if err == nil {
			_ = coverView.SetImage(image)
		}
	}
	detailsDialog.Show()
	txtView.SetTextSelection(-1, -1)
	_ = txtView.SetFocus()
	detailsDialogIsOpen = true
}
