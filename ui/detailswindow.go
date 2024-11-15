package ui

import (
	"EmbyExplorer_for_Windows/assets"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var detailsWindow *walk.MainWindow
var imgView *walk.ImageView
var txtEdit *walk.TextEdit
var button *walk.PushButton
var detailsIsOpen = false

const (
	wndWidth       = 660
	wndHeight      = 350
	coverMaxWidth  = "250"
	coverMaxHeight = "250"
)

func showDetails() {
	if err := (MainWindow{
		AssignTo:        &detailsWindow,
		Title:           assets.CapDetails,
		Icon:            "/assets/app.ico",
		MinSize:         Size{Width: wndWidth, Height: wndHeight},
		MaxSize:         Size{Width: wndWidth, Height: wndHeight},
		Size:            Size{Width: wndWidth, Height: wndHeight},
		Layout:          VBox{},
		DoubleBuffering: true,
		OnSizeChanged: func() {
			_ = detailsWindow.SetHeight(wndHeight)
			_ = detailsWindow.SetWidth(wndWidth)
			_ = detailsWindow.Invalidate()
		},
		Children: []Widget{
			Composite{Layout: Grid{Columns: 2, MarginsZero: true, SpacingZero: true},
				Children: []Widget{
					ImageView{
						AssignTo:   &imgView,
						MaxSize:    Size{Width: 250, Height: 250},
						Background: SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						Alignment:  AlignHCenterVCenter,
						Mode:       ImageViewModeIdeal,
					},
					TextEdit{
						AssignTo:        &txtEdit,
						MinSize:         Size{Width: 400, Height: 250},
						Alignment:       AlignHCenterVCenter,
						Background:      SolidColorBrush{Color: walk.RGB(255, 255, 255)},
						DoubleBuffering: true,
						ReadOnly:        true,
						Enabled:         false,
						HScroll:         false,
						VScroll:         false,
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &button,
						Text:     assets.CapClose,
						Enabled:  true,
						OnClicked: func() {
							detailsIsOpen = false
							_ = detailsWindow.Close()
						},
					},
				},
			},
		},
	}.Create()); err == nil {
		detailsIsOpen = true
		detailsWindow.Show()
	}
}

func setDetails(imagePath string, overview string) {
	_ = txtEdit.SetText(overview)
	if imagePath != "" {
		image, err := walk.NewImageFromFileForDPI(imagePath, 96)
		if err == nil {
			_ = imgView.SetImage(image)
		}
	}
}
