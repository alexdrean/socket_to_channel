package socket_to_channel

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

func BindSocketToChannel(socket net.Conn, incoming chan <-string, outgoing <-chan string) {
	end := make(chan byte)
	go func() {
		reader := bufio.NewReader(socket)
		for {
			line, err := reader.ReadString(byte('\n'))
			if line != "" {
				incoming <- strings.TrimSuffix(line, "\n")
			}
			if err != nil {
				if err != io.EOF {
					end <- 1
					_ = fmt.Errorf("unexpected error while reading %+v: %s", socket, err.Error())
				}
				return
			}
		}
	}()
	go func() {
		for s := range outgoing {
			_, err := socket.Write([]byte(s))
			if err != nil {
				end <- 2
				_ = fmt.Errorf("unexpected error while writing %+v: %s", socket, err.Error())
				return
			}
		}
	}()
	<-end
	return
}

func DialToChannel(addr net.Addr, retryDelay time.Duration) (<-chan string, chan <-string) {

	incoming, outgoing := make(chan string), make(chan string)
	go func() {
		for {
			conn, _ := net.Dial(addr.Network(), addr.String())
			if conn != nil {
				BindSocketToChannel(conn, incoming, outgoing)
			} else {
				time.Sleep(retryDelay)
			}
		}
	}()
	return incoming, outgoing
}
