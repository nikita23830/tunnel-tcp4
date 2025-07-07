package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type TunnelConfig struct {
	Name    string `json:"name"`
	Port    int    `json:"port"`
	To      string `json:"to"`
	HAProxy bool   `json:"haproxy"`
}

func main() {
	configs, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Cant load config: %v", err)
	}

	for _, cfg := range configs {
		cfg := cfg
		go startTunnel(cfg)
	}

	select {}
}

func loadConfig(path string) ([]TunnelConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var configs []TunnelConfig
	if err := json.NewDecoder(file).Decode(&configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func startTunnel(cfg TunnelConfig) {
	addr := fmt.Sprintf(":%d", cfg.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[%s] Failed to start listening on %s: %v", cfg.Name, addr, err)
	}
	log.Printf("[%s] The tunnel is running: %s -> %s (HAProxy: %v)", cfg.Name, addr, cfg.To, cfg.HAProxy)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[%s] Connection acceptance error: %v", cfg.Name, err)
			continue
		}
		go handleConnection(conn, cfg)
	}
}

func handleConnection(src net.Conn, cfg TunnelConfig) {
	defer src.Close()

	dst, err := net.Dial("tcp", cfg.To)
	if err != nil {
		log.Printf("[%s] Server connection error: %v", cfg.Name, err)
		return
	}
	defer dst.Close()

	if cfg.HAProxy {
		clientAddr := src.RemoteAddr().(*net.TCPAddr)
		dstAddr := dst.RemoteAddr().(*net.TCPAddr)

		proxyLine := fmt.Sprintf("PROXY TCP4 %s %s %d %d\r\n",
			clientAddr.IP.String(),
			dstAddr.IP.String(),
			clientAddr.Port,
			dstAddr.Port,
		)

		_, err = dst.Write([]byte(proxyLine))
		if err != nil {
			log.Printf("[%s] Error when sending the PROXY line: %v", cfg.Name, err)
			return
		}
	}

	go io.Copy(dst, src)
	io.Copy(src, dst)
}
