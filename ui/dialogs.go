//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Error dialogs, using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"EmbyExplorer_for_Windows/assets"
	"github.com/lxn/walk"
)

func dialogToDisplaySystemError(primary string, detail error) {
	walk.MsgBox(mainWindow, assets.CapError, primary+"\n"+detail.Error(),
		walk.MsgBoxOK+walk.MsgBoxIconError+walk.MsgBoxApplModal)
}

/*
func dialogToDisplayErrorMessage(primary string, detail string) {
	walk.MsgBox(mainWindow, assets.CapError, primary+"\n"+detail,
		walk.MsgBoxOK+walk.MsgBoxIconError+walk.MsgBoxApplModal)
}
*/
