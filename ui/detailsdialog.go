//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Display cover and summary for movies & series, using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"EmbyExplorer_for_Windows/assets"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	dialogWidth    = 500
	dialogHeight   = 330
	coverMaxWidth  = "250"
	coverMaxHeight = "250"
)

var detailsDialog *walk.Dialog
var imageView *walk.ImageView
var textEdit *walk.TextEdit
var okButton *walk.PushButton
var detailsIsOpen = false

func showDetails() {
	var err error
	err = Dialog{
		AssignTo:      &detailsDialog,
		Title:         assets.CapDetails,
		Icon:          "/assets/app.ico",
		MinSize:       Size{Width: dialogWidth, Height: dialogHeight},
		FixedSize:     true,
		DefaultButton: &okButton,
		Layout:        VBox{},
		Children: []Widget{
			Composite{Layout: Grid{Columns: 2, MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					ImageView{
						AssignTo:   &imageView,
						MaxSize:    Size{Width: 250, Height: 250},
						Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						Alignment:  AlignHCenterVCenter,
						Mode:       ImageViewModeIdeal,
					},
					TextEdit{
						AssignTo:  &textEdit,
						MinSize:   Size{Width: 250, Height: 250},
						Alignment: AlignHCenterVCenter,
						ReadOnly:  true,
						HScroll:   false,
						VScroll:   true,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &okButton,
						Text:     assets.CapClose,
						Enabled:  true,
						OnClicked: func() {
							detailsIsOpen = false
							detailsDialog.Accept()
						},
					},
				},
			},
		},
	}.Create(mainWindow)
	if err == nil {
		detailsIsOpen = true
		detailsDialog.Show()
	}
}

func setDetails(imagePath string, overview string) {
	// cannot fix: text comes up selected, func SetTextSelection(0, -1) does nothing here
	_ = textEdit.SetText(overview)
	if imagePath != "" {
		image, err := walk.NewImageFromFileForDPI(imagePath, 96)
		if err == nil {
			_ = imageView.SetImage(image)
		}
	}
}
