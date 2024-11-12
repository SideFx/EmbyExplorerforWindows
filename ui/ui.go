//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// UI, using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package ui

import (
	"EmbyExplorer_for_Windows/api"
	"EmbyExplorer_for_Windows/assets"
	"EmbyExplorer_for_Windows/export"
	"EmbyExplorer_for_Windows/models"
	"EmbyExplorer_for_Windows/settings"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const (
	wndMinWidth  = 680
	wndMinHeight = 400
)

const (
	tmpFile1 = "_embyTmpFile1_"
	tmpFile2 = "_embyTmpFile2_"
)

var mainWindow *walk.MainWindow
var tvMovies, tvSeries, tvVideos *walk.TableView
var colsMovies, colsSeries, colsVideos []TableViewColumn
var movieModel *models.MovieModel
var seriesModel *models.SeriesModel
var videoModel *models.VideoModel
var prefAction, authAction, fetchAction, detailsAction, exportAction,
	quitAction, aboutAction *walk.Action
var tab *walk.TabWidget
var moviesPage, seriesPage, videosPage *walk.TabPage
var handleTabChange bool
var activePage int
var movieCollId, seriesCollId, videoCollId string
var movieSelection, seriesSelection = -1, -1
var tempFolder string
var lastExportFolder string

func init() {
	tempFolder = os.TempDir()
}

func CreateUi() error {
	colsMovies = models.GetMovieColumns()
	colsSeries = models.GetSeriesColumns()
	colsVideos = models.GetVideosColumns()
	movieModel = models.NewMovieModel()
	seriesModel = models.NewSeriesModel()
	videoModel = models.NewVideoModel()
	activePage = 0
	if err := (MainWindow{
		AssignTo: &mainWindow,
		Title:    assets.AppName + " " + assets.Version,
		Icon:     "/assets/app.ico",
		MinSize:  Size{Width: wndMinWidth, Height: wndMinHeight},
		MenuItems: []MenuItem{
			Menu{
				Text: assets.CapFile,
				Items: []MenuItem{
					Action{
						AssignTo:    &quitAction,
						Text:        assets.CapQuit,
						Enabled:     true,
						Visible:     true,
						Shortcut:    Shortcut{Modifiers: walk.ModControl, Key: walk.KeyQ},
						OnTriggered: quitActionTriggered,
					},
				},
			},
			Menu{
				Text: assets.CapHelp,
				Items: []MenuItem{
					Action{
						AssignTo:    &aboutAction,
						Text:        assets.CapAbout,
						Enabled:     true,
						Visible:     true,
						OnTriggered: aboutActionTriggered,
					},
				},
			},
		},
		ToolBar: ToolBar{
			ButtonStyle: ToolBarButtonImageBeforeText,
			Items: []MenuItem{
				Action{
					AssignTo:    &prefAction,
					Text:        assets.CapPreferences,
					Image:       "/assets/preferences.png",
					Enabled:     true,
					Visible:     true,
					OnTriggered: prefActionTriggered,
				},
				Action{
					AssignTo:    &authAction,
					Text:        assets.CapAuthenticate,
					Image:       "/assets/authenticate.png",
					Enabled:     true,
					Visible:     true,
					OnTriggered: authActionTriggered,
				},
				Separator{},
				Action{
					AssignTo:    &fetchAction,
					Text:        assets.CapFetch,
					Image:       "/assets/fetch.png",
					Enabled:     false,
					Visible:     true,
					OnTriggered: fetchActionTriggered,
				},
				Action{
					AssignTo:    &detailsAction,
					Text:        assets.CapDetails,
					Image:       "/assets/details.png",
					Enabled:     false,
					Visible:     true,
					OnTriggered: detailsActionTriggered,
				},
				Separator{},
				Action{
					AssignTo:    &exportAction,
					Text:        assets.CapExport,
					Image:       "/assets/export.png",
					Enabled:     false,
					Visible:     true,
					OnTriggered: exportActionTriggered,
				},
			},
		},
		Layout: VBox{MarginsZero: false},
		Children: []Widget{
			TabWidget{
				AssignTo:              &tab,
				Visible:               false,
				Enabled:               true,
				OnCurrentIndexChanged: onTabChanged,
				Pages: []TabPage{
					{
						AssignTo: &moviesPage,
						Enabled:  false,
						Title:    assets.CapNotAvailable,
						Content: TableView{
							AssignTo:         &tvMovies,
							Columns:          colsMovies,
							MultiSelection:   false,
							CheckBoxes:       false,
							ColumnsOrderable: false,
							Visible:          true,
							Enabled:          true,
							StyleCell: func(style *walk.CellStyle) {
								style.BackgroundColor = walk.RGB(255, 255, 255)
								style.TextColor = walk.RGB(0, 0, 0)
							},
							OnCurrentIndexChanged: func() {
								moviesSelectionChanged()
							},
						},
					},
					{
						AssignTo: &seriesPage,
						Enabled:  false,
						Title:    assets.CapNotAvailable,
						Content: TableView{
							AssignTo:         &tvSeries,
							Columns:          colsSeries,
							MultiSelection:   false,
							CheckBoxes:       false,
							ColumnsOrderable: false,
							Visible:          true,
							Enabled:          true,
							StyleCell: func(style *walk.CellStyle) {
								style.BackgroundColor = walk.RGB(255, 255, 255)
								style.TextColor = walk.RGB(0, 0, 0)
							},
							OnCurrentIndexChanged: func() {
								seriesSelectionChanged()
							},
						},
					},
					{
						AssignTo: &videosPage,
						Enabled:  false,
						Title:    assets.CapNotAvailable,
						Content: TableView{
							AssignTo:         &tvVideos,
							Columns:          colsVideos,
							MultiSelection:   false,
							CheckBoxes:       false,
							ColumnsOrderable: false,
							Visible:          true,
							Enabled:          true,
							StyleCell: func(style *walk.CellStyle) {
								style.BackgroundColor = walk.RGB(255, 255, 255)
								style.TextColor = walk.RGB(0, 0, 0)
							},
							OnCurrentIndexChanged: func() {
								videosSelectionChanged()
							},
						},
					},
				},
			},
		},
	}.Create()); err != nil {
		return err
	}
	closingEvent := mainWindow.AsFormBase().Closing()
	closingEvent.Attach(onClosing)
	e := settings.LoadPreferences()
	if e == nil {
		lastExportFolder = settings.GetLastExportFolder()
		bounds := settings.GetWindowBounds()
		if bounds.Width < wndMinWidth {
			bounds.Width = wndMinWidth
		}
		if bounds.Height < wndMinHeight {
			bounds.Height = wndMinHeight
		}
		_ = mainWindow.SetBounds(bounds)
		if settings.Valid() {
			h, s, p, u, x := settings.GetConnectionSettings()
			api.InitApiPreferences(h, s, p, u, x)
			_ = authAction.SetEnabled(true)
		}
	}
	// initialize tab content
	handleTabChange = false
	_ = tab.SetCurrentIndex(0)
	_ = tab.SetCurrentIndex(1)
	_ = tab.SetCurrentIndex(2)
	_ = tab.SetCurrentIndex(activePage)
	handleTabChange = true
	mainWindow.Run()
	return nil
}

func quitActionTriggered() {
	_ = mainWindow.Close()
}

func prefActionTriggered() {
	preferencesDialog()
	if settings.Valid() {
		h, s, p, u, x := settings.GetConnectionSettings()
		api.InitApiPreferences(h, s, p, u, x)
	}
}

func authActionTriggered() {
	var err error
	var views []api.UserView
	err = api.AuthenticateUserInt()
	if err != nil {
		dialogToDisplaySystemError(assets.ErrAuthFailed, err)
		return
	}
	views, err = api.UserGetViewsInt()
	if err != nil {
		dialogToDisplaySystemError(assets.ErrFetchViewsFailed, err)
		return
	}
	for _, view := range views {
		switch view.CollectionType {
		case api.CollectionMovies:
			moviesPage.SetEnabled(true)
			_ = moviesPage.SetTitle(view.Name)
			movieCollId = view.Id
			break
		case api.CollectionTVShows:
			seriesPage.SetEnabled(true)
			_ = seriesPage.SetTitle(view.Name)
			seriesCollId = view.Id
			break
		case api.CollectionHomeVideos:
			videosPage.SetEnabled(true)
			_ = videosPage.SetTitle(view.Name)
			videoCollId = view.Id
			break
		}
	}
	_ = fetchAction.SetEnabled(true)
	tab.SetVisible(true)
}

func fetchActionTriggered() {
	var dto []api.BaseItemDto
	var err error
	switch activePage {
	case 0:
		if moviesPage.Enabled() && movieCollId != "" {
			dto, err = api.UserGetItemsInt(movieCollId, api.CollectionMovies)
			if err != nil {
				dialogToDisplaySystemError(assets.ErrFetchItemsFailed, err)
				return
			}
			movieData := api.GetMovieDisplayData(dto)
			movieModel.SetItems(movieData)
			_ = tvMovies.SetModel(movieModel)
			if movieModel.RowCount() > 0 {
				_ = tvMovies.SetCurrentIndex(0)
				movieSelection = 0
			}
		}
		_ = detailsAction.SetEnabled(true)
		_ = exportAction.SetEnabled(true)
		break
	case 1:
		if seriesPage.Enabled() && seriesCollId != "" {
			dto, err = api.UserGetItemsInt(seriesCollId, api.CollectionTVShows)
			if err != nil {
				dialogToDisplaySystemError(assets.ErrFetchItemsFailed, err)
				return
			}
			seriesData := api.GetSeriesDisplayData(dto)
			seriesModel.SetItems(seriesData)
			_ = tvSeries.SetModel(seriesModel)
			if seriesModel.RowCount() > 0 {
				_ = tvSeries.SetCurrentIndex(0)
				seriesSelection = 0
			}
		}
		_ = detailsAction.SetEnabled(true)
		_ = exportAction.SetEnabled(true)
		break
	case 2:
		if videosPage.Enabled() && videoCollId != "" {
			dto, err = api.UserGetItemsInt(videoCollId, api.CollectionHomeVideos)
			if err != nil {
				dialogToDisplaySystemError(assets.ErrFetchItemsFailed, err)
				return
			}
			videoData := api.GetVideoDisplayData(dto)
			videoModel.SetItems(videoData)
			_ = tvVideos.SetModel(videoModel)
			if videoModel.RowCount() > 0 {
				_ = tvVideos.SetCurrentIndex(0)
			}
		}
		_ = detailsAction.SetEnabled(false)
		_ = exportAction.SetEnabled(true)
		break
	}
}

func detailsActionTriggered() {
	if activePage == 0 {
		if movieSelection >= 0 {
			showDetails()
			movieGetDetail()
		}
	} else {
		if seriesSelection >= 0 {
			showDetails()
			seriesGetDetail()
		}
	}
}

func exportActionTriggered() {
	buildAndExport()
}

func aboutActionTriggered() {
	aboutDialog()
}

func moviesSelectionChanged() {
	movieSelection = tvMovies.CurrentIndex()
	movieGetDetail()
}

func seriesSelectionChanged() {
	seriesSelection = tvSeries.CurrentIndex()
	seriesGetDetail()
}

func videosSelectionChanged() {
	// nothing to do, no details
}

func onTabChanged() {
	if handleTabChange {
		activePage = tab.CurrentIndex()
		_ = detailsAction.SetEnabled(activePage != 2) // no details for home videos
	}
}

func onClosing(canceled *bool, _ walk.CloseReason) {
	*canceled = false
	settings.SetWindowBounds(mainWindow.Bounds())
	settings.SetLastExportFolder(lastExportFolder)
	settings.SavePreferences()
}

func movieGetDetail() {
	var err error
	var rawImage []byte
	if movieSelection >= 0 && detailsIsOpen {
		data := movieModel.GetItem(movieSelection)
		rawImage, err = api.GetPrimaryImageForItemInt(data.MovieId, api.ImageFormatPng, coverMaxWidth, coverMaxHeight)
		if err == nil {
			p := filepath.Join(tempFolder, tmpFile1)
			err = os.WriteFile(p, rawImage, 0644)
			if err != nil {
				p = ""
			}
			setDetails(p, data.Overview)
		}
	}
}

func seriesGetDetail() {
	var err error
	var rawImage []byte
	if seriesSelection >= 0 && detailsIsOpen {
		data := seriesModel.GetItem(seriesSelection)
		rawImage, err = api.GetPrimaryImageForItemInt(data.SeasonId, api.ImageFormatPng, coverMaxWidth, coverMaxHeight)
		if err == nil {
			p := filepath.Join(tempFolder, tmpFile2)
			err = os.WriteFile(p, rawImage, 0644)
			if err != nil {
				p = ""
			}
			setDetails(p, data.Overview)
		}
	}
}

func buildAndExport() {
	var i, j int
	var sheet string
	var exp = make([]export.Payload, 0)
	var hdr = make([]export.HeaderData, 0)
	var e export.Payload
	var c export.HeaderData
	j = 1 // xlsx start row
	switch activePage {
	case 0:
		for i = 0; i < models.MoviesTableDescription.NoOfColumns; i++ {
			c.XLSCell = models.MoviesTableDescription.Columns[i].XLSColumn + strconv.Itoa(j)
			c.Name = models.MoviesTableDescription.Columns[i].Caption
			c.Column = models.MoviesTableDescription.Columns[i].XLSColumn
			c.Width = models.MoviesTableDescription.Columns[i].XLSColumnWidth
			hdr = append(hdr, c)
		}
		for _, m := range movieModel.GetItems() {
			j++
			for i = 0; i < models.MoviesTableDescription.NoOfColumns; i++ {
				e.XLSCell = models.MoviesTableDescription.Columns[i].XLSColumn + strconv.Itoa(j)
				e.Data = movieModel.ItemValue(m, i)
				exp = append(exp, e)
			}
		}
		sheet = assets.CapMovies
		break
	case 1:
		for i = 0; i < models.SeriesTableDescription.NoOfColumns; i++ {
			c.XLSCell = models.SeriesTableDescription.Columns[i].XLSColumn + strconv.Itoa(j)
			c.Name = models.SeriesTableDescription.Columns[i].Caption
			c.Column = models.SeriesTableDescription.Columns[i].XLSColumn
			c.Width = models.SeriesTableDescription.Columns[i].XLSColumnWidth
			hdr = append(hdr, c)
		}
		for _, t := range seriesModel.GetItems() {
			j++
			for i = 0; i < models.SeriesTableDescription.NoOfColumns; i++ {
				e.XLSCell = models.SeriesTableDescription.Columns[i].XLSColumn + strconv.Itoa(j)
				e.Data = seriesModel.ItemValue(t, i)
				exp = append(exp, e)
			}
		}
		sheet = assets.CapSeries
		break
	case 2:
		for i = 0; i < models.VideosTableDescription.NoOfColumns; i++ {
			c.XLSCell = models.VideosTableDescription.Columns[i].XLSColumn + strconv.Itoa(j)
			c.Name = models.VideosTableDescription.Columns[i].Caption
			c.Column = models.VideosTableDescription.Columns[i].XLSColumn
			c.Width = models.VideosTableDescription.Columns[i].XLSColumnWidth
			hdr = append(hdr, c)
		}
		for _, h := range videoModel.GetItems() {
			j++
			for i = 0; i < models.VideosTableDescription.NoOfColumns; i++ {
				e.XLSCell = models.VideosTableDescription.Columns[i].XLSColumn + strconv.Itoa(j)
				e.Data = videoModel.ItemValue(h, i)
				exp = append(exp, e)
			}
		}
		sheet = assets.CapVideos
		break
	default:
		return
	}
	date := time.Now().Format("2006-01-02")
	folder := settings.GetLastExportFolder()
	if folder == "" {
		folder, _ = os.UserHomeDir()
	}
	preferredFileName := assets.CapEmby + " " + sheet + " " + date + "." + assets.FileExtension
	dlg := new(walk.FileDialog)
	dlg.FilePath = filepath.Join(lastExportFolder, preferredFileName)
	dlg.Filter = createFilter()
	dlg.Title = assets.CapSave
	ok, err := dlg.ShowSave(mainWindow)
	if err == nil && ok {
		p := dlg.FilePath
		lastExportFolder, _ = path.Split(p)
		settings.SetLastExportFolder(lastExportFolder)
		err = export.XlsxExport(exp, hdr, p, sheet)
	}
	if err != nil {
		dialogToDisplaySystemError(assets.CapError, err)
	}
}

func createFilter() string {
	return assets.CapExcelFiles + " " + "(*." + assets.FileExtension + ")" +
		"|" + "*." + assets.FileExtension
}
