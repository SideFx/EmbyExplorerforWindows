//----------------------------------------------------------------------------------------------------------------------
// (w) 2024 by Jan Buchholz
// Models - common types
//----------------------------------------------------------------------------------------------------------------------

package models

type ColumnDescription struct {
	Caption        string
	XLSColumn      string
	XLSColumnWidth float64
}

type TableDescription struct {
	NoOfColumns int
	APIFields   string
	Columns     []ColumnDescription
}
