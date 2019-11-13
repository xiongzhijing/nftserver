package protocol

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"hash/crc32"
	"io"
	"log"
)

type Protocol struct {
	Ver    [4]byte
	Length uint32
	Crc32  uint32
	Data   []byte //json data
}

type NFTParamReq struct {
	Func       string `json:"func"`
	FilePath   string `json:"filepath"`
	RouteMode  string `json:"routemode"`
	RouteID    string `json:"routeid"`
	Direction  string `json:"direction"`
	AppType    string `json:"apptype"`
	FileMode   string `json:"filemode"`
	UserMsg    string `json:"usermsg"`
	TransferID string `json:"transferid"`
	User       string `json:"user"`
	Password   string `json:"password"`
}

type NFTParamRet struct {
	RetCode       int32  `json:"retcode"`
	TransferID    string `json:"transferid"`
	TransferState string `json:"transferState"`
	RetMsg        string `json:"retmsg"`
}

var (
	errChecksum = errors.New("protocol checksum error")
	errVer      = errors.New("protocol ver error")
	errLength   = errors.New("protocol length error")
	errData     = errors.New("protocol data error")
)

func (p *Protocol) Unpack(reader io.Reader) error {
	var err error

	if _, err = reader.Read(p.Ver[:4]); err != nil {
		log.Println(err)
		return errVer
	}
	if err = binary.Read(reader, binary.BigEndian, &p.Length); err != nil {
		log.Println(err)
		return errLength
	}

	if err = binary.Read(reader, binary.BigEndian, &p.Crc32); err != nil {
		log.Println(err)
		return errChecksum
	}

	p.Data = make([]byte, p.Length)
	if _, err = reader.Read(p.Data); err != nil {
		log.Println(err)
		return errData
	}

	if c := crc32.ChecksumIEEE(p.Data); c != p.Crc32 {
		log.Printf("[ERROR] %X!=%X", c, p.Crc32)
		return errChecksum
	}

	return nil
}

/* 将请求报文转换为结构 */
func (r *NFTParamReq) Unpack(buf []byte, unmarshal func(data []byte, v interface{}) error) error {

	if err := unmarshal(buf, r); err != nil {
		log.Println("[ERROR] cannot unmarshal error", err)
		log.Println("\n" + hex.Dump(buf))
		return err
	}
	return nil

}
