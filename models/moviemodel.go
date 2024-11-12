//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Model for collection type "movies", using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package models

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sort"
)

var MoviesTableDescription = TableDescription{
	12,
	"Name,OriginalTitle,MediaSources,Path,Genres,ProductionYear,People,Studios," +
		"Width,Height,Container," + "Overview,RunTimeTicks,Type_", //no spaces here!
	[]ColumnDescription{
		{"Title", "A", 70},
		{"Original Title", "B", 70},
		{"Year", "C", 10},
		{"Time", "D", 10},
		{"Actors", "E", 100},
		{"Director", "F", 50},
		{"Studio", "G", 30},
		{"Genre", "H", 70},
		{"Ext.", "I", 10},
		{"Codec", "J", 20},
		{"Resolution", "K", 15},
		{"Path", "L", 80},
	},
}

type MovieData struct {
	Name           string
	OriginalTitle  string
	ProductionYear string
	Runtime        string
	Actors         string
	Directors      string
	Studios        string
	Genres         string
	Container      string
	Codecs         string
	Resolution     string
	Path           string
	Overview       string
	MovieId        string
}

type MovieModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*MovieData
}

func NewMovieModel() *MovieModel {
	m := new(MovieModel)
	m.ResetRows()
	return m
}

func (m *MovieModel) SetItems(movieItems []*MovieData) {
	m.items = movieItems
}

func (m *MovieModel) GetItem(index int) *MovieData {
	return m.items[index]
}

func (m *MovieModel) GetItems() []*MovieData {
	return m.items
}

func (m *MovieModel) ResetRows() {
	for i := range m.items {
		m.items[i] = &MovieData{}
	}
	m.PublishRowsReset()
	_ = m.Sort(m.sortColumn, m.sortOrder)
}

func (m *MovieModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}
		switch m.sortColumn {
		case 0:
			return c(a.Name < b.Name)
		case 1:
			return c(a.OriginalTitle < b.OriginalTitle)
		case 2:
			return c(a.ProductionYear < b.ProductionYear)
		case 3:
			return c(a.Runtime < b.Runtime)
		case 4:
			return c(a.Actors < b.Actors)
		case 5:
			return c(a.Directors < b.Directors)
		case 6:
			return c(a.Studios < b.Studios)
		case 7:
			return c(a.Genres < b.Genres)
		case 8:
			return c(a.Container < b.Container)
		case 9:
			return c(a.Codecs < b.Codecs)
		case 10:
			return c(a.Resolution < b.Resolution)
		case 11:
			return c(a.Path < b.Path)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

func (m *MovieModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.OriginalTitle
	case 2:
		return item.ProductionYear
	case 3:
		return item.Runtime
	case 4:
		return item.Actors
	case 5:
		return item.Directors
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

func (m *MovieModel) ItemValue(item *MovieData, col int) string {
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.OriginalTitle
	case 2:
		return item.ProductionYear
	case 3:
		return item.Runtime
	case 4:
		return item.Actors
	case 5:
		return item.Directors
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

func (m *MovieModel) RowCount() int {
	return len(m.items)
}

func GetMovieColumns() []TableViewColumn {
	var cols []TableViewColumn
	for _, c := range MoviesTableDescription.Columns {
		col := TableViewColumn{
			Title: c.Caption,
			Width: int(c.XLSColumnWidth * 4),
		}
		cols = append(cols, col)
	}
	return cols
}
