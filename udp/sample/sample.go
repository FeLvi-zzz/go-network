package sample

import (
	"fmt"
	"io"

	"github.com/FeLvi-zzz/go-network/ipv4"
	"github.com/FeLvi-zzz/go-network/udp"
)

func Serve(sender *ipv4.Sender, addr []byte, port uint16) error {
	s := udp.NewService(sender)
	l := s.Listen(addr, port)

	for {
		conn := l.Accept()
		handle(conn, conn)
	}
}

func handle(r io.ReadCloser, w io.Writer) {
	defer r.Close()

	s, err := io.ReadAll(r)
	if err != nil {
		panic(err)
	}

	fmt.Printf("req: `%s`\n", s)
	fmt.Fprintf(w, "Hello, you said %s\n", s)
}

func RequestHoge(sender *ipv4.Sender, addr []byte, port uint16) {
	s := udp.NewService(sender)
	conn := s.Dial(addr, port, addr, port)
	conn.Write([]byte("hoge"))

	b := make([]byte, 100)
	conn.Read(b)

	fmt.Printf("%s", b)
}
