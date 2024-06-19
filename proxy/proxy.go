package proxy

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/mikuta0407/tcpproxy/config"
)

func Proxy(configFile string) {

	config, err := config.LoadConfigFile(configFile)
	if err != nil {
		log.Fatalln(err)
	}

	if len(config.Proxies) == 0 {
		log.Fatal("no proxies setting")
	}

	for _, proxy := range config.Proxies {
		fmt.Printf("Name: %s\nSource: %s\nDestination: %s\n\n", proxy.Name, proxy.Source, proxy.Destination)
		go listen(proxy.Source, proxy.Destination)
	}

	select {}

}

func listen(src, dst string) {
	listener, err := net.Listen("tcp", src)
	if err != nil {
		panic("connection error:" + err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Accept Error:", err)
			continue
		}
		copyConn(conn, dst)
	}
}

func copyConn(src net.Conn, dstAddr string) {
	dst, err := net.Dial("tcp", dstAddr)
	if err != nil {
		panic("Dial Error:" + err.Error())
	}

	done := make(chan struct{})

	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(dst, src)
		done <- struct{}{}
	}()

	go func() {
		defer src.Close()
		defer dst.Close()
		io.Copy(src, dst)
		done <- struct{}{}
	}()

	<-done
	<-done
}
