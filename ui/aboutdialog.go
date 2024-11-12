//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// About dialog, using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"EmbyExplorer_for_Windows/assets"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const (
	aboutDialogWidth  = 300
	aboutDialogHeight = 350
)

func aboutDialog() {
	var dialog *walk.Dialog
	var okButton *walk.PushButton
	_, _ = Dialog{
		AssignTo:      &dialog,
		Title:         assets.CapAbout + " " + assets.AppName,
		Icon:          "/assets/app.ico",
		DefaultButton: &okButton,
		MinSize:       Size{Width: aboutDialogWidth, Height: aboutDialogHeight},
		FixedSize:     true,
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 1},
				Children: []Widget{
					TextLabel{
						Alignment:     AlignHCenterVNear,
						Text:          assets.TxtAboutEmbyExplorer,
						TextAlignment: AlignHCenterVNear,
					},
					ImageView{
						Alignment: AlignHCenterVNear,
						Image:     "/assets/gopher.png",
					},
					TextLabel{
						Alignment:     AlignHCenterVNear,
						Text:          assets.TxtCredits,
						TextAlignment: AlignHCenterVNear,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &okButton,
						Text:     assets.CapOk,
						Enabled:  true,
						OnClicked: func() {
							dialog.Accept()
						},
					},
				},
			},
		},
	}.Run(mainWindow)
}
