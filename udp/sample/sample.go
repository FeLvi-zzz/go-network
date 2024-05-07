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

func RequestHoge(sender *ipv4.Sender, raddr []byte, rport uint16, laddr []byte, lport uint16) error {
	s := udp.NewService(sender)
	conn := s.Dial(raddr, rport, laddr, lport)
	defer conn.Close()

	if _, err := conn.Write([]byte("hoge")); err != nil {
		return err
	}

	fmt.Println("send: hoge!!")

	b := make([]byte, 100)
	if _, err := conn.Read(b); err != nil && err != io.EOF {
		return err
	}

	fmt.Printf("recv: %s!!!!!!!!!!!!!!!!!!!!!!!!\n", b)

	return nil
}
