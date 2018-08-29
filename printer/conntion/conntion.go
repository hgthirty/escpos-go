package conntion

const (
	Offline = 0x01
	Online  = 0x02
	Busy    = 0x03
)

type Status uint8

type Connetion interface {
	Write(content []byte) (int, error)
	Read([]byte) (int, error)
	Info() Status
}
