package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrConnectionClosed = errors.New("connection closed by peer")
	ErrConnectionEmpty  = errors.New("connection is empty")
)

type TelnetClient interface {
	Connect() error
	Close() error
	Send() error
	Receive() error
}

type Client struct {
	address    string
	timeout    time.Duration
	connection net.Conn
	in         io.ReadCloser
	out        io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (client *Client) Connect() error {
	connection, err := net.DialTimeout("tcp", client.address, client.timeout)
	if err != nil {
		return err
	}

	client.connection = connection

	return nil
}

func (client *Client) Close() error {
	if client.connection == nil {
		return ErrConnectionEmpty
	}

	err := client.in.Close()
	if err != nil {
		return err
	}

	return client.connection.Close()
}

func (client *Client) Send() error {
	if client.connection == nil {
		return ErrConnectionEmpty
	}

	_, err := io.Copy(client.connection, client.in)
	if err != nil {
		return ErrConnectionClosed
	}

	return nil
}

func (client *Client) Receive() error {
	if client.connection == nil {
		return ErrConnectionEmpty
	}

	_, err := io.Copy(client.out, client.connection)
	if err != nil {
		return ErrConnectionClosed
	}

	return nil
}
