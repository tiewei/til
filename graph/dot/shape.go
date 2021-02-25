package dot

import "encoding/xml"

// newCustomShape creates a custom graph shape that is serializable to HTML.
//
// Its format is a table with 2 vertical cells: a header and a body.
//
//  +------------+  <-- bordercolor=accentColor
//  |  <header>  |  <-- bgcolor=accentColor;fontcolor=headerTxtColor
//  +------------+
//  |   <body>   |
//  +------------+
//
func newCustomShape(header, body, accentColor, headerTxtColor string) *customShape {
	headerCell := tableData{
		BgColor: accentColor,
		Font: font{
			Color: headerTxtColor,
			Value: header,
		},
	}

	bodyCell := tableData{
		Font: font{
			Value: body,
		},
	}

	return &customShape{
		Color: accentColor,

		CellBorder:  1,
		CellPadding: 4,

		TableRows: []tableRow{{
			TableData: headerCell,
		}, {
			TableData: bodyCell,
		}},
	}
}

// customShape is a custom DOT shape that can be serialized to a HTML table
// with multiple rows, each containing a single cell.
type customShape struct {
	XMLName xml.Name `xml:"table"`

	Border      int    `xml:"border,attr"`
	Color       string `xml:"color,attr"` // border color
	CellBorder  int    `xml:"cellborder,attr"`
	CellSpacing int    `xml:"cellspacing,attr"`
	CellPadding int    `xml:"cellpadding,attr"`

	TableRows []tableRow `xml:"tr"`
}

// tableRow represents a <TR> element (row) in the customShape HTML.
type tableRow struct {
	TableData tableData `xml:"td"`
}

// tableData represents a <TD> element (cell) in the customShape HTML.
type tableData struct {
	BgColor string `xml:"bgcolor,attr,omitempty"`
	Font    font   `xml:"font"`
}

// font is a <FONT> element in a HTML cell.
type font struct {
	Color string `xml:"color,attr,omitempty"`
	Value string `xml:",chardata"`
}
