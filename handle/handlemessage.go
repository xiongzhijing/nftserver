package handle

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"hash/crc32"
	"io"
	"log"
	"net"
	"nftserver/nftapi"
	"nftserver/protocol"
)

const (
	MAXBUFSIZ = 16 * 1024
)

func HandleConn(conn net.Conn) {

	log.Println(conn.RemoteAddr().String() + " connected!")
	rdchan := make(chan []byte, 100)
	go handleMessage(rdchan, conn)

	defer func() {
		log.Println("[INFO] close connection", conn.RemoteAddr())
		conn.Close()
		close(rdchan)
	}()

	buffer := make([]byte, MAXBUFSIZ)
	var s = 0
	var packLen uint32 = 0
	for {

		n, err := conn.Read(buffer[s:])
		log.Println("[INFO] received ", n, " bytes")
		log.Println("\n" + hex.Dump(buffer[s:s+n]))
		if err == io.EOF && s != 0 {
			/* 将最后接收的数据直接传到channel */
			rdchan <- buffer[:s+n]
			break
		}
		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			//conn.Close()
			break
		}
		s += n
		if s >= 8 {
			if string(buffer[:4]) == "V!01" ||
				string(buffer[:4]) == "V!02" {
				packLen = binary.BigEndian.Uint32(buffer[4:8])
				if packLen+8+4 > MAXBUFSIZ {
					log.Println(conn.RemoteAddr(), "[ERROR] message too big")
					//conn.Close()
					break
				}
				//p += 8
			}
		}

		if s >= 8+4+int(packLen) && packLen != 0 {
			rdchan <- buffer[:8+4+int(packLen)]
			//将多出来的字节往前移
			copy(buffer[:], buffer[12+int(packLen):s])
			s = s - 12 - int(packLen)

		}
	}

}

func handleMessage(rdchan <-chan []byte, conn net.Conn) {

	for {
		buf, ok := <-rdchan
		if !ok {
			break
		}
		//log.Println("\n" + hex.Dump(buf))
		p := &protocol.Protocol{}
		if err := p.Unpack(bytes.NewBuffer(buf)); err != nil {
			log.Println("[ERROR] unpack message error", err)
			continue
		}
		if string(p.Ver[:]) == "V!01" /* JSON */ ||
			string(p.Ver[:]) == "V!02" /* XML */ {

			req := &protocol.NFTParamReq{}
			var (
				unmarshal func(data []byte, v interface{}) error
				marshal   func(v interface{}) ([]byte, error)
			)
			if string(p.Ver[:]) == "V!01" {
				unmarshal = json.Unmarshal
				marshal = json.Marshal

			} else if string(p.Ver[:]) == "V!02" {
				unmarshal = xml.Unmarshal
				marshal = xml.Marshal
			}
			if err := req.Unpack(p.Data, unmarshal); err != nil {
				log.Println("[ERROR] unpack request message error", err)
				continue
			}

			ret := protocol.NFTParamRet{}
			switch req.Func {
			case "ntransfer":
				ret = nftapi.Ntransfer(req)
			case "getTsfResponse":
				ret = nftapi.GetTsfResponse(req)
			case "getTsfState":
				ret = nftapi.GetTsfState(req)
			case "ncancelTransfer":
				ret = nftapi.NCancelTransfer(req)
			default:
				log.Println("[ERROR] invalid function")
				continue
			}
			if string(p.Ver[:]) == "V!01" ||
				string(p.Ver[:]) == "V!02" {
				str, err := marshal(ret)
				if err != nil {
					log.Println("[ERROR] cannot marshal response object to json string", ret)
					continue
				}
				log.Println("[INFO] " + string(str))
				go SendMessage(conn, p.Ver[:], str)
			}
		} else {
			log.Println("[ERROR] cannot identify message type")
		}
	}
}

func SendMessage(conn net.Conn, ver []byte, data []byte) {
	p := protocol.Protocol{}
	copy(p.Ver[:], ver)
	p.Data = data

	p.Crc32 = crc32.ChecksumIEEE(p.Data)
	p.Length = uint32(len(p.Data))

	err := binary.Write(conn, binary.BigEndian, p.Ver)
	if err != nil {
		log.Println("[ERROR] cannot write ver" + err.Error())
	}

	err = binary.Write(conn, binary.BigEndian, &p.Length)
	if err != nil {
		log.Println("[ERROR] cannot write length" + err.Error())
	}
	err = binary.Write(conn, binary.BigEndian, &p.Crc32)
	if err != nil {
		log.Println("[ERROR] cannot write crc32" + err.Error())
	}
	err = binary.Write(conn, binary.BigEndian, p.Data)
	if err != nil {
		log.Println("[ERROR] cannot write data" + err.Error())
	}

}

func ListenAndServe(addr string) {

	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("[ERROR] cannot listen address ", addr)
	}

	defer l.Close()

	log.Println("[INFO] Server start OK! ", l.Addr())
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go HandleConn(conn)
	}
}
