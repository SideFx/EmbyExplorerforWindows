//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Model for collection type "tvshows", using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package models

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var SeriesTableDescription = TableDescription{
	NoOfColumns: 12,
	APIFields: "Name,MediaSources,Path,Genres,ProductionYear,People,Studios," +
		"Width,Height,Container,RunTimeTicks," +
		"Overview,SeriesId,SeasonId,Id,ParentId,IndexNumber,Type_", //no spaces here!
	Columns: []ColumnDescription{
		{"Series", "A", 50},
		{"Episode", "B", 50},
		{"Season", "C", 30},
		{"Year", "D", 10},
		{"Time", "E", 10},
		{"Actors", "F", 100},
		{"Studio", "G", 30},
		{"Genre", "H", 70},
		{"Ext.", "I", 10},
		{"Codec", "J", 20},
		{"Resolution", "K", 15},
		{"Path", "L", 80},
	},
}

type SeriesData struct {
	Name           string
	Episode        string
	Season         string
	ProductionYear string
	Runtime        string
	Actors         string
	Studios        string
	Genres         string
	Container      string
	Codecs         string
	Resolution     string
	Path           string
	Overview       string
	SeriesId       string
	SeasonId       string
	EpisodeId      string
	Type_          string
	SortIndex      int32
}

type SeriesModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*SeriesData
}

func NewSeriesModel() *SeriesModel {
	m := new(SeriesModel)
	m.ResetRows()
	return m
}

func (m *SeriesModel) SetItems(seriesItems []*SeriesData) {
	m.items = seriesItems
}

func (m *SeriesModel) GetItem(index int) *SeriesData {
	return m.items[index]
}

func (m *SeriesModel) GetItems() []*SeriesData {
	return m.items
}

func (m *SeriesModel) ResetRows() {
	for i := range m.items {
		m.items[i] = &SeriesData{}
	}
	m.PublishRowsReset()
	_ = m.Sort(m.sortColumn, m.sortOrder)
}

func (m *SeriesModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	return m.SorterBase.Sort(col, order) //no sort
}

func (m *SeriesModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.Episode
	case 2:
		return item.Season
	case 3:
		return item.ProductionYear
	case 4:
		return item.Runtime
	case 5:
		return item.Actors
	case 6:
		return item.Studios
	case 7:
		return item.Genres
	case 8:
		return item.Container
	case 9:
		return item.Codecs
	case 10:
		return item.Resolution
	case 11:
		return item.Path
	}
	panic("unexpected col")
}

func (m *SeriesModel) ItemValue(item *SeriesData, col int) string {
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.Episode
	case 2:
		return item.Season
	case 3:
		return item.ProductionYear
	case 4:
		return item.Runtime
	case 5:
		return item.Actors
	case 6:
		return item.Studios
	case 7:
		return item.Genres
	case 8:
		return item.Container
	case 9:
		return item.Codecs
	case 10:
		return item.Resolution
	case 11:
		return item.Path
	}
	panic("unexpected col")
}

func (m *SeriesModel) RowCount() int {
	return len(m.items)
}

func GetSeriesColumns() []TableViewColumn {
	var cols []TableViewColumn
	for _, c := range SeriesTableDescription.Columns {
		col := TableViewColumn{
			Title: c.Caption,
			Width: int(c.XLSColumnWidth * 4),
		}
		cols = append(cols, col)
	}
	return cols
}
