//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Model for collection type "homevideos", using walk GUI toolkit
// https://github.com/lxn/walk
//----------------------------------------------------------------------------------------------------------------------

package models

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"sort"
)

var VideosTableDescription = TableDescription{
	NoOfColumns: 7,
	APIFields: "Name,MediaSources,Path,Width,Height,Container," +
		"RunTimeTicks,ParentId,Type_", //no spaces here!
	Columns: []ColumnDescription{
		{"Title", "A", 100},
		{"Folder", "B", 30},
		{"Time", "C", 10},
		{"Ext.", "D", 10},
		{"Codec", "E", 20},
		{"Resolution", "F", 15},
		{"Path", "G", 150},
	},
}

type VideoData struct {
	Name       string
	Folder     string
	Runtime    string
	Container  string
	Codecs     string
	Resolution string
	Path       string
	FolderId   string
	ParentId   string
}

type VideoModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*VideoData
}

func NewVideoModel() *VideoModel {
	m := new(VideoModel)
	m.ResetRows()
	return m
}

func (m *VideoModel) SetItems(videoItems []*VideoData) {
	m.items = videoItems
}

func (m *VideoModel) GetItems() []*VideoData {
	return m.items
}

func (m *VideoModel) ResetRows() {
	for i := range m.items {
		m.items[i] = &VideoData{}
	}
	m.PublishRowsReset()
	_ = m.Sort(m.sortColumn, m.sortOrder)
}

func (m *VideoModel) Sort(col int, order walk.SortOrder) error {
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
			return c(a.Folder < b.Folder)
		case 2:
			return c(a.Runtime < b.Runtime)
		case 3:
			return c(a.Container < b.Container)
		case 4:
			return c(a.Codecs < b.Codecs)
		case 5:
			return c(a.Resolution < b.Resolution)
		case 6:
			return c(a.Path < b.Path)
		}
		panic("unreachable")
	})
	return m.SorterBase.Sort(col, order)
}

func (m *VideoModel) Value(row, col int) interface{} {
	item := m.items[row]
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.Folder
	case 2:
		return item.Runtime
	case 3:
		return item.Container
	case 4:
		return item.Codecs
	case 5:
		return item.Resolution
	case 6:
		return item.Path
	}
	panic("unexpected col")
}

func (m *VideoModel) ItemValue(item *VideoData, col int) string {
	switch col {
	case 0:
		return item.Name
	case 1:
		return item.Folder
	case 2:
		return item.Runtime
	case 3:
		return item.Container
	case 4:
		return item.Codecs
	case 5:
		return item.Resolution
	case 6:
		return item.Path
	}
	panic("unexpected col")
}

func (m *VideoModel) RowCount() int {
	return len(m.items)
}

func GetVideosColumns() []TableViewColumn {
	var cols []TableViewColumn
	for _, c := range VideosTableDescription.Columns {
		col := TableViewColumn{
			Title: c.Caption,
			Width: int(c.XLSColumnWidth * 4),
		}
		cols = append(cols, col)
	}
	return cols
}
