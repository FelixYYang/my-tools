// Package communication tcp实现
package communication

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
)

func NewTcpClient(address string, config TcpConfig) *TcpClient {
	if config.MaxSendMsgNum == 0 {
		config.MaxSendMsgNum = 16
	}
	if config.MaxRecvMsgNum == 0 {
		config.MaxRecvMsgNum = 16
	}
	return &TcpClient{
		address: address,
		config:  config,
	}
}

type TcpConfig struct {
	TlsEnabled    bool
	TLsConfig     *tls.Config
	MaxRecvMsgNum uint32
	MaxSendMsgNum uint32
}

type TcpClient struct {
	address  string
	config   TcpConfig
	conn     net.Conn
	handler  *tcpMessageHandler
	cancel   func()
	status   uint8
	runError error
}

func (t *TcpClient) Connect() error {
	t.runError = t.Close()
	if err := t.dial(); err != nil {
		return err
	}
	connHandler := TcpMessageHandler(t.conn,
		tcpConnHandlerOption{
			maxRecvMsgNum: t.config.MaxRecvMsgNum,
			maxSendMsgNum: t.config.MaxSendMsgNum,
		})
	t.handler = connHandler
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel
	go func() {
		err := t.handler.Work(ctx)
		t.runError = err
	}()
	return nil
}

func (t *TcpClient) Recv() <-chan *Message {
	if t.handler != nil {
		return t.handler.Recv()
	}
	return nil
}

func (t *TcpClient) Close() error {
	if t.cancel != nil {
		t.cancel()
	}
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *TcpClient) GetLastErr() error {
	return t.runError
}

func (t *TcpClient) Send(message *Message) error {
	if t.handler != nil {
		return t.handler.Send(message)
	}
	return nil
}

func (t *TcpClient) dial() error {
	config := t.config
	//解析地址
	addr, e := net.ResolveTCPAddr("tcp", t.address)
	if e != nil {
		return e
	}
	var conn net.Conn
	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		return DialTcpError{
			Addr: *addr,
			Err:  err,
		}
	}
	if config.TlsEnabled {
		conn = tls.Client(conn, config.TLsConfig)
	}
	t.conn = conn
	return nil
}

type DialTcpError struct {
	Addr net.TCPAddr
	Err  error
}

func (d DialTcpError) Error() string {
	return fmt.Sprintf("dial ip fail :%s", d.Addr.String())
}
