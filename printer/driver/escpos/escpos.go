package escpos

import (
	"io"
	"strings"
	"fmt"
	"strconv"
	"log"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"os"
	"math"
	"github.com/gw123/escpos-go/util"
)

// Text replacement map
var textReplaceMap = map[string]string{
	// horizontal tab
	"&#9;":  "\x09",
	"&#x9;": "\x09",

	// linefeed
	"&#10;": "\n",
	"&#xA;": "\n",

	// xml stuff
	"&apos;": "'",
	"&quot;": `"`,
	"&gt;":   ">",
	"&lt;":   "<",

	// ampersand must be last to avoid double decoding
	"&amp;": "&",
}

// replace Text from the above map
func textReplace(data string) string {
	for k, v := range textReplaceMap {
		data = strings.Replace(data, k, v, -1)
	}
	return data
}

type Escpos struct {
	// destination
	dst io.ReadWriter

	// font metrics
	width, height uint8

	// state toggles ESC[char]
	underline  uint8
	emphasize  uint8
	upsidedown uint8
	rotate     uint8
	lineWidth  uint8 //打印机一行可以打印多少字符
	// state toggles GS[char]
	reverse, smooth uint8
}

// reset toggles
func (e *Escpos) reset() {
	e.width = 1
	e.height = 1

	e.underline = 0
	e.emphasize = 0
	e.upsidedown = 0
	e.rotate = 0

	e.reverse = 0
	e.smooth = 0
	e.lineWidth = 32
}

// create Escpos printer
func NewEscpos(dst io.ReadWriter) (e *Escpos) {
	e = &Escpos{dst: dst}
	e.Init()
	return
}

// init/reset printer settings
func (e *Escpos) Init() {
	e.reset()
	e.WriteRaw( []byte{ ESC , 0x40} )
}

// end output
func (e *Escpos) End() {
	e.Write("\xFA")
}

// send cut
func (e *Escpos) Cut() {
	e.Write("\x1DVA0")
}

// send cut minus one point (partial cut)
func (e *Escpos) CutPartial() {
	e.WriteRaw([]byte{GS, 0x56, 1})
}

// send cash
func (e *Escpos) Cash() {
	e.Write("\x1B\x70\x00\x0A\xFF")
}

// send linefeed
func (e *Escpos) Linefeed() {
	e.Write("\n")
}

// send N formfeeds
func (e *Escpos) FormfeedN(n int) {
	e.Write(fmt.Sprintf("\x1Bd%c", n))
}

// send formfeed
func (e *Escpos) Formfeed() {
	e.FormfeedN(1)
}

// set font
func (e *Escpos) SetFont(font string) {
	f := 0
	switch font {
	case "A":
		f = 0
	case "B":
		f = 1
	case "C":
		f = 2
	default:
		log.Fatalf("Invalid font: '%s', defaulting to 'A'", font)
		f = 0
	}

	e.Write(fmt.Sprintf("\x1BM%c", f))
}

func (e *Escpos) SendFontSize() {
	e.Write(fmt.Sprintf("\x1D!%c", ((e.width-1)<<4)|(e.height-1)))
}

// set font size
func (e *Escpos) SetFontSize(width, height uint8) {
	if width > 0 && height > 0 && width <= 8 && height <= 8 {
		e.width = width
		e.height = height
		e.SendFontSize()
	} else {
		log.Fatalf("Invalid font size passed: %d x %d", width, height)
	}
}

// send underline
func (e *Escpos) SendUnderline() {
	e.Write(fmt.Sprintf("\x1B-%c", e.underline))
}

// send emphasize / doublestrike
func (e *Escpos) SendEmphasize() {
	e.Write(fmt.Sprintf("\x1BG%c", e.emphasize))
}

// send upsidedown
func (e *Escpos) SendUpsidedown() {
	e.Write(fmt.Sprintf("\x1B{%c", e.upsidedown))
}

// send rotate
func (e *Escpos) SendRotate() {
	e.Write(fmt.Sprintf("\x1BR%c", e.rotate))
}

// send reverse
func (e *Escpos) SendReverse() {
	e.Write(fmt.Sprintf("\x1DB%c", e.reverse))
}

// send smooth
func (e *Escpos) SendSmooth() {
	e.Write(fmt.Sprintf("\x1Db%c", e.smooth))
}

// send move x
func (e *Escpos) SendMoveX(x uint16) {
	e.Write(string([]byte{0x1b, 0x24, byte(x % 256), byte(x / 256)}))
}

// send move y
func (e *Escpos) SendMoveY(y uint16) {
	e.Write(string([]byte{0x1d, 0x24, byte(y % 256), byte(y / 256)}))
}

// set underline
func (e *Escpos) SetUnderline(v uint8) {
	e.underline = v
	e.SendUnderline()
}

// set emphasize
func (e *Escpos) SetEmphasize(u uint8) {
	e.emphasize = u
	e.SendEmphasize()
}

// set upsidedown
func (e *Escpos) SetUpsidedown(v uint8) {
	e.upsidedown = v
	e.SendUpsidedown()
}

// set rotate
func (e *Escpos) SetRotate(v uint8) {
	e.rotate = v
	e.SendRotate()
}

// set reverse
func (e *Escpos) SetReverse(v uint8) {
	e.reverse = v
	e.SendReverse()
}

// set smooth
func (e *Escpos) SetSmooth(v uint8) {
	e.smooth = v
	e.SendSmooth()
}

// pulse (open the drawer)
func (e *Escpos) Pulse() {
	// with t=2 -- meaning 2*2msec
	e.Write("\x1Bp\x02")
}

// set alignment
func (e *Escpos) SetAlign(align string) {
	a := 0
	switch align {
	case "left":
		a = 0
	case "center":
		a = 1
	case "right":
		a = 2
	default:
		log.Fatalf("Invalid alignment: %s", align)
	}
	e.Write(fmt.Sprintf("\x1Ba%c", a))
}

// set language -- ESC R
func (e *Escpos) SetLang(lang string) {
	l := 0

	switch lang {
	case "en":
		l = 0
	case "fr":
		l = 1
	case "de":
		l = 2
	case "uk":
		l = 3
	case "da":
		l = 4
	case "sv":
		l = 5
	case "it":
		l = 6
	case "es":
		l = 7
	case "ja":
		l = 8
	case "no":
		l = 9
	default:
		log.Fatalf("Invalid language: %s", lang)
	}
	e.Write(fmt.Sprintf("\x1BR%c", l))
}

// feed the printer
func (e *Escpos) Feed(params map[string]string) {
	// handle lines (form feed X lines)
	if l, ok := params["line"]; ok {
		if i, err := strconv.Atoi(l); err == nil {
			e.FormfeedN(i)
		} else {
			log.Fatalf("Invalid line number %s", l)
		}
	}

	// handle units (dots)
	if u, ok := params["unit"]; ok {
		if i, err := strconv.Atoi(u); err == nil {
			e.SendMoveY(uint16(i))
		} else {
			log.Fatalf("Invalid unit number %s", u)
		}
	}

	// send linefeed
	e.Linefeed()

	// reset variables
	e.reset()

	// reset printer
	e.SendEmphasize()
	e.SendRotate()
	e.SendSmooth()
	e.SendReverse()
	e.SendUnderline()
	e.SendUpsidedown()
	e.SendFontSize()
	e.SendUnderline()
}

// feed and cut based on parameters
func (e *Escpos) FeedAndCut(params map[string]string) {
	if t, ok := params["type"]; ok && t == "feed" {
		e.Formfeed()
	}

	e.Cut()
}

/***
 * 写入byte数据
 */
func (e *Escpos) WriteRaw(data []byte) (n int, err error) {
	if len(data) > 0 {
		log.Printf("Writing %d bytes\n", len(data))
		e.dst.Write(data)
	} else {
		log.Printf("Wrote NO bytes\n")
	}

	return 0, nil
}

/***
 * 读取
 */
func (e *Escpos) ReadRaw(data []byte) (n int, err error) {
	return e.dst.Read(data)
}

/***
 * 写入字符串
 */
func (e *Escpos) Write(data string) (int, error) {
	return e.WriteRaw([]byte(data))
}

/***
 * 写入带有中午的字符串 （data 的编码为utf-8 如果data编码为gbk直接使用Write）
 */
func (e *Escpos) WriteGbk(data string) (int, error) {
	str := util.ConvertToGbk(data)
	return e.WriteRaw([]byte(str))
}

/***
 * 写入一行 支持左右对齐中间自动补齐空格
 */
func (e *Escpos) WriteLRLine(left, right string) (int, error) {
	left_str := util.ConvertToGbk(left)
	right_str := util.ConvertToGbk(right)
	sum := len(left_str) + len(right_str)

	fmt.Printf("sum : %d \n", sum)
	paddingNum := (int(e.lineWidth) - sum*int(e.width)) / (int(e.width))
	if paddingNum < 2 {
		paddingNum = 2
	}
	for i := 0; i < paddingNum; i++ {
		left_str += " "
	}
	return e.WriteRaw([]byte(left_str + right_str))
}

/***
 * 使用同一个字符写入一行
 */
func (e *Escpos) WriteALine(byte2 byte) (int, error) {
	e.SetFontSize(1, 1)
	str := ""
	for i := 0; i < int(e.lineWidth); i++ {
		str += string([]byte{byte2})
	}
	return e.Write(str)
}

/***
 * 对于xml中的一个p标签
 */
func (e *Escpos) WriteLine(line Line) {
	str := ""
	var fontSize uint8
	if line.Size < 1 {
		fontSize = 1
	} else {
		fontSize = line.Size
	}
	e.SetFontSize(fontSize, fontSize)
	var align string
	if line.Align == "" {
		align = "left"
	} else {
		align = line.Align
	}
	e.SetAlign(align)

	if line.Loop.Text != "" {
		e.WriteALine([]byte(line.Loop.Text)[0])
	} else if line.Img.Text != "" {
		e.WriteRemoteImage(line.Img.Text)
	} else if line.Qrcode.Text != "" {
		e.WriteRemoteImage(line.Qrcode.Text)
	} else if len(line.Cells) > 0 {
		e.WriteCells(line.Cells)
	} else {
		data := textReplace(line.Text)
		str += util.ConvertToGbk(data)
		str += "\n"
		e.Write(str)
	}
}

/***
 * 将解析后的xml打印出来
 */
func (e *Escpos) WriteXml(root Root) {
	for _, line := range root.Lines {
		e.WriteLine(line)
	}
	e.Linefeed()
}

/***
 * 写入远程图片
 */
func (e *Escpos) WriteRemoteImage(url string) {
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	e.WriteImage(response.Body)
}

/***
 * 写入本地图片
 */
func (e *Escpos) WriteLocalImage(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}
	e.WriteImage(file)
}

/***
 * 写入图片
 */
func (e *Escpos) WriteImage(r io.Reader) {
	img := getImageInfo(r)
	e.WriteRaw([]byte{0x1D, 0x76, 0x30, 0x00, img.Xl, img.Xh, img.Yl, img.Yh})
	e.WriteRaw(img.Data)
	e.Linefeed()
}

/***
 * 写入一行单元格
 */
func (e *Escpos) WriteCells(cells []Cell) {
	str := ""
	for _, cell := range cells {
		width := int(math.Floor(float64(float32(e.lineWidth) * cell.Width)))
		padding := width - int(util.GetStrLength(cell.Text))
		padding = padding / int(e.width)
		//fmt.Printf("%d %d %d\n", padding, int(GetStrLength(cell.Text)), width)
		if cell.Align == "left" {
			str += util.ConvertToGbk(cell.Text)
			for i := 0; i < padding; i++ {
				str += " "
			}
		} else if cell.Align == "center" {
			for i := 0; i < padding/2; i++ {
				str += " "
			}
			str += util.ConvertToGbk(cell.Text)
			for i := 0; i < padding/2; i++ {
				str += " "
			}
		} else {
			for i := 0; i < padding; i++ {
				str += " "
			}
			str += util.ConvertToGbk(cell.Text)
		}
		str = textReplace(str)
	}
	fmt.Printf("%d\n", len(str))
	str += "\n"
	e.Write(str)
}

/***
 * 读取打印机状态
 */
func (e *Escpos) ReadStatus(n byte) (byte, error) {
	e.WriteRaw([]byte{DLE, EOT, n})
	data := make([]byte, 1)
	_, err := e.ReadRaw(data)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

type PinterImageInfo struct {
	Img  image.Image
	Data []byte
	Xl   uint8
	Xh   uint8
	Yl   uint8
	Yh   uint8
}

/***
 * 图片解析
 */
func getImageInfo(r io.Reader) (imageInfo PinterImageInfo) {
	imageHanel, t, err := image.Decode(r)
	if err != nil {
		fmt.Println(t, err)
	}

	rect := imageHanel.Bounds()

	width := (rect.Size().X + 7) / 8
	height := rect.Size().Y
	imageInfo.Xl = uint8(width % 256)
	imageInfo.Xh = uint8(width / 256)
	imageInfo.Yl = uint8(height % 256)
	imageInfo.Yh = uint8(height / 256)
	realWidth := rect.Size().X

	newImgData := make([]byte, width*height)
	var part = [8]byte{}
	index := 0
	for y := 0; y < height; y ++ {
		for x := 0; x < realWidth; x += 8 {
			for k := 0; k < 8; k++ {
				if x+k >= realWidth {
					part[k] = 0
				} else {
					rgba := imageHanel.At(x+k, y)
					r, g, b, _ := rgba.RGBA()
					value := float32(r>>8)*0.3 + float32(g>>8)*0.59 + float32(b>>8)*0.11
					if value > 160 {
						value = 0
					} else {
						value = 1
					}
					part[k] = byte(value)
				}
			}

			newImgData[index] = part[0]<<7 | part[1]<<6 + part[2]<<5 + part[3]<<4 + part[4]<<3 + part[5]<<2 + part[6]<<1 + part[7];
			index += 1
		}
	}
	imageInfo.Data = newImgData
	return
}
