//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Emby connection settings, using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"EmbyExplorer_for_Windows/assets"
	"EmbyExplorer_for_Windows/settings"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

const defaultPort = "8096"

var protocolBox *walk.CheckBox
var editServer, editPort, editUser, editPassword *walk.LineEdit
var cancelBtn, okBtn *walk.PushButton

// cannot fix: LineEdit for port cannot be sized (MaxSize does nothing)

func preferencesDialog() {
	const (
		dialogWidth  = 300
		dialogHeight = 250
	)
	var dialog *walk.Dialog
	var err error
	err = Dialog{
		AssignTo:      &dialog,
		Title:         assets.CapEmbyPreferences,
		Icon:          "/assets/app.ico",
		DefaultButton: &okBtn,
		CancelButton:  &cancelBtn,
		MinSize:       Size{Width: dialogWidth, Height: dialogHeight},
		FixedSize:     true,
		Layout:        VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{Text: assets.CapProtocol},
					CheckBox{
						AssignTo: &protocolBox,
					},
					Label{Text: assets.CapServer},
					LineEdit{
						AssignTo: &editServer,
						OnTextChanged: func() {
							checkComplete()
						},
					},
					Label{Text: assets.CapPort},
					LineEdit{
						AssignTo:  &editPort,
						MaxLength: 4,
						Text:      defaultPort,
						OnTextChanged: func() {
							checkComplete()
						},
					},
					Label{Text: assets.CapUser},
					LineEdit{
						AssignTo: &editUser,
						OnTextChanged: func() {
							checkComplete()
						},
					},
					Label{Text: assets.CapPassword},
					LineEdit{
						AssignTo:     &editPassword,
						PasswordMode: true,
						OnTextChanged: func() {
							checkComplete()
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &okBtn,
						Text:     assets.CapOk,
						Enabled:  false,
						OnClicked: func() {
							setSettings()
							dialog.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelBtn,
						Text:      assets.CapCancel,
						Enabled:   true,
						OnClicked: func() { dialog.Cancel() },
					},
				},
			},
		},
	}.Create(mainWindow)
	if err == nil {
		h, s, p, u, x := settings.GetConnectionSettings()
		protocolBox.SetChecked(h)
		_ = editServer.SetText(s)
		_ = editPort.SetText(p)
		_ = editUser.SetText(u)
		_ = editPassword.SetText(x)
		dialog.Show()
	}
}

func checkComplete() {
	okBtn.SetEnabled(editServer.Text() != "" &&
		editPort.Text() != "" &&
		editUser.Text() != "" &&
		editPassword.Text() != "")
}

func setSettings() {
	settings.SetPreferencesDetail(
		protocolBox.Checked(), editServer.Text(), editPort.Text(),
		editUser.Text(), editPassword.Text())
}
