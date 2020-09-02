package main // import "github.com/simon-engledew/pg-go-rpc/src"

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgproto3/v2"
	"log"
	"net"
	"strings"
)

type Backend struct {
	backend   *pgproto3.Backend
	conn      net.Conn
}

func NewBackend(conn net.Conn) *Backend {
	backend := pgproto3.NewBackend(pgproto3.NewChunkReader(conn), conn)

	connHandler := &Backend{
		backend:   backend,
		conn:      conn,
	}

	return connHandler
}

func (b *Backend) Run() error {
	defer b.Close()

	err := b.handleStartup()
	if err != nil {
		return err
	}

	for {
		msg, err := b.backend.Receive()
		if err != nil {
			return fmt.Errorf("error receiving message: %w", err)
		}

		switch msg.(type) {
		case *pgproto3.Query:
			log.Print(msg)
			query := msg.(*pgproto3.Query)

			log.Print(query.String)

			buf := (&pgproto3.RowDescription{Fields: []pgproto3.FieldDescription{
				{
					Name:                 []byte("rpc"),
					TableOID:             0,
					TableAttributeNumber: 0,
					DataTypeOID:          25,
					DataTypeSize:         -1,
					TypeModifier:         -1,
					Format:               0,
				},
			}}).Encode(nil)

			if strings.HasPrefix(query.String, "RPC ") {
				var request interface{}

				if err := json.Unmarshal([]byte(strings.TrimPrefix(query.String, "RPC ")), &request); err != nil {
					panic(err)
				}

				response, err := json.Marshal(request)

				if err != nil {
					panic(err)
				}

				buf = (&pgproto3.DataRow{Values: [][]byte{response}}).Encode(buf)
			}

			buf = (&pgproto3.CommandComplete{CommandTag: []byte(query.String)}).Encode(buf)
			buf = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
			_, err = b.conn.Write(buf)
		case *pgproto3.Terminate:
			return nil
		default:
			return fmt.Errorf("received message other than Query from client: %#v", msg)
		}
	}
}

func (b *Backend) handleStartup() error {
	startupMessage, err := b.backend.ReceiveStartupMessage()
	if err != nil {
		return fmt.Errorf("error receiving startup message: %w", err)
	}

	switch startupMessage.(type) {
	case *pgproto3.StartupMessage:
		buf := (&pgproto3.AuthenticationOk{}).Encode(nil)
		buf = (&pgproto3.ReadyForQuery{TxStatus: 'I'}).Encode(buf)
		_, err = b.conn.Write(buf)
		if err != nil {
			return fmt.Errorf("error sending ready for query: %w", err)
		}
	case *pgproto3.SSLRequest:
		_, err = b.conn.Write([]byte("N"))
		if err != nil {
			return fmt.Errorf("error sending deny SSL request: %w", err)
		}
		return b.handleStartup()
	default:
		return fmt.Errorf("unknown startup message: %#v", startupMessage)
	}

	return nil
}

func (b *Backend) Close() error {
	return b.conn.Close()
}

func main() {
	ln, err := net.Listen("tcp","0.0.0.0:15432")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Listening on", ln.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Accepted connection from", conn.RemoteAddr())

		b := NewBackend(conn)
		go func() {
			err := b.Run()
			if err != nil {
				log.Println(err)
			}
			log.Println("Closed connection from", conn.RemoteAddr())
		}()
	}
}