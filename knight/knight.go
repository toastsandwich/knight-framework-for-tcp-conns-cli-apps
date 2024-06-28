package knight

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
)

var ErrPatternAlreadyFound = errors.New("knight: pattern already appointed to something")

type Open interface {
	OpenPoint
}

type OpenPoint func(req *Request, res *Response)

type Knight struct {
	rwlock  *sync.RWMutex
	Faction string
	Addr    string
	conns   sync.Map
	routes  map[string]OpenPoint //map of openpoints like url: localhost:8080/home
	errorch chan error
}

func Suitup(addr, faction string) *Knight {
	if faction != "tcp" && faction != "udp" {
		log.Fatal("knight should belong to type -> tcp or udp")
	}
	return &Knight{
		Addr:    addr,
		Faction: faction,

		rwlock:  &sync.RWMutex{},
		errorch: make(chan error),
		routes:  make(map[string]OpenPoint),
	}
}

func (k *Knight) Serve() error {
	log.Println("knight listening on ", k.Addr, "[", k.Faction, "]")
	ln, err := net.Listen(k.Faction, k.Addr)
	if err != nil {
		return fmt.Errorf("knight failed to listen on %s: %v", k.Addr, err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("knight failed to connect: ", err)
			continue
		}
		log.Println("[+] connection established with ", conn.RemoteAddr().String())

		// On successful connection store it to conns for further usage
		k.conns.Store(conn.RemoteAddr().String(), conn)

		go k.handleConn(conn)
	}
}

func (k *Knight) readRequest(conn net.Conn) (*Request, error) {
	k.rwlock.RLock()
	defer k.rwlock.RUnlock()
	buf := make([]byte, 1024)

	// get from
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	from := string(buf[:n])

	// get hit
	n, err = conn.Read(buf)
	if err != nil {
		return nil, err
	}
	hit := string(buf[:n])

	// logic for params

	return &Request{
		From: from,
		Hit:  hit,
	}, nil
}

func (k *Knight) handleConn(conn net.Conn) {
	defer func() {
		k.conns.Delete(conn.RemoteAddr().String())
		conn.Close()
	}()
	req, err := k.readRequest(conn)
	if err != nil {
		k.errorch <- err
	}
	var res *Response
	k.routes[req.Hit](req, res)
}

func (k *Knight) HandlePoint(pattern string, openPoint OpenPoint) error {
	if _, ok := k.routes[pattern]; ok {
		return ErrPatternAlreadyFound
	}
	k.routes[pattern] = openPoint
	return nil
}
