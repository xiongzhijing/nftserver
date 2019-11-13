package main

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"hash/crc32"
	"log"
	"net"
	"nftserver/protocol"
	"os"
	"time"
)

func sender(conn net.Conn) {
	for i := 0; i < 10000; i++ {
		p := protocol.Protocol{Ver: [4]byte{'V', '!', '0', '1'}}
		req := protocol.NFTParamReq{}
		req.FilePath = "aaa.lst"
		req.RouteMode = "1"
		req.RouteID = "G1"
		req.Direction = "R"
		req.FileMode = "0"
		req.UserMsg = "xxxx"
		req.TransferID = ""
		req.User = "S1"
		req.Password = ""
		req.Func = "ntransfer"
		req.AppType = "1"
		if i%2 == 0 {
			p.Data, _ = json.Marshal(req)

		} else {
			p.Ver[3] = '2'
			p.Data, _ = xml.Marshal(req)

		}
		log.Println(string(p.Data))
		//p.Data = []byte("hello")
		//hash := md5.Sum(p.Data)
		//copy(p.Digest[:32], hex.EncodeToString(hash[:]))
		p.Crc32 = crc32.ChecksumIEEE(p.Data)
		p.Length = uint32(len(p.Data))

		binary.Write(conn, binary.BigEndian, p.Ver)
		binary.Write(conn, binary.BigEndian, &p.Length)
		binary.Write(conn, binary.BigEndian, &p.Crc32)
		binary.Write(conn, binary.BigEndian, p.Data)
		time.Sleep(time.Millisecond * 2000)
		b := make([]byte, 1024)
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		n, err := conn.Read(b)
		log.Println("received bytes ", n, err)
		log.Println("\n" + hex.Dump(b[:n]))

	}
}

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	server := "128.1.104.76:9989"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	func() {
		for i := 0; i < 300; i++ {
			conn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
				os.Exit(1)
			}
			fmt.Println("connect success", i)

			go sender(conn)
		}
	}()

	for {

		time.Sleep(time.Second)
	}
}
