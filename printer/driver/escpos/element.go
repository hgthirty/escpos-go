package escpos

import "encoding/xml"

type Root struct {
	XMLName xml.Name `xml:"root"` //p标签
	Lines   []Line   `xml:"p"`
}

/***
 * 一行
 */
type Line struct {
	XMLName xml.Name `xml:"p"` //p标签
	Cells   []Cell   `xml:"cell"`
	Loop    Loop
	Img     Img
	Qrcode  Qrcode
	Align   string `xml:"align,attr"` // 读取align属性
	Size    uint8  `xml:"size,attr"`  // strong
	Strong  uint8  `xml:"strong,attr"`
	Text    string `xml:",cdata"`
}

func NewLine(size uint8, text, align string) *Line {
	line := new(Line)
	line.Align = align
	line.Text = text
	line.Size = size
	return line
}

func (l *Line) AppendCell(cell Cell) {
	//tex ：= sync.Mutex{}
	l.Cells = append(l.Cells, cell)
}

/***
 * 单元格
 */
type Cell struct {
	XMLName xml.Name `xml:"cell"`       //p标签
	Width   float32  `xml:"width,attr"` // width
	Align   string   `xml:"align,attr"` // width
	Text    string   `xml:",cdata"`
}

func NewCell(width float32, text, align string) *Cell {
	instance := new(Cell)
	instance.Width = width
	instance.Align = align
	instance.Text = text
	return instance
}

/***
 * 循环
 */
type Loop struct {
	XMLName xml.Name `xml:"loop"` //loop标签
	Text    string   `xml:",cdata"`
}

/***
 * 图片
 */
type Img struct {
	XMLName xml.Name `xml:"img"` //img标签
	Text    string   `xml:",cdata"`
}

/***
 * 二维码
 */
type Qrcode struct {
	XMLName xml.Name `xml:"qrcode"` //qrcode标签
	Text    string   `xml:",cdata"`
}
