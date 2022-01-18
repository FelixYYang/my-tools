package communication

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"sync"
)

func TcpMessageHandler(conn net.Conn, option tcpConnHandlerOption) *tcpMessageHandler {
	return &tcpMessageHandler{
		conn:       conn,
		msgChan:    make(chan *Message, option.maxRecvMsgNum),
		sendChan:   make(chan *Message, option.maxSendMsgNum),
		errors:     make(chan error, 10),
		packMaxLen: 2 << 20,
	}
}

type tcpConnHandlerOption struct {
	maxRecvMsgNum uint32
	maxSendMsgNum uint32
}

type tcpMessageHandler struct {
	conn net.Conn

	// 消息缓存通道
	msgChan chan *Message
	// 消息发送缓存通道
	sendChan   chan *Message
	errors     chan error
	packMaxLen uint32
	statusLock sync.RWMutex
	// 状态 0:待启动，1:已启动，-1:已关闭
	status int8
}

func (hd *tcpMessageHandler) Send(message *Message) error {
	if hd.status == -1 {
		return fmt.Errorf("tcpMessageHandler isnot running,status: %d", hd.status)
	}
	hd.sendChan <- message
	return nil
}

func (hd *tcpMessageHandler) Recv() <-chan *Message {
	return hd.msgChan
}

func (hd *tcpMessageHandler) Work(ctx context.Context) (e error) {
	defer func() {
		if err := recover(); err != nil {
			e = fmt.Errorf("tcpMessageHandler::handle panic: %s", err)
		}
		hd.destruct()
	}()
	if hd.status != 0 {
		return fmt.Errorf("tcpMessageHandler isnot really,status: %d", hd.status)
	}

	go hd.handleSend()
	go hd.handleRead()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-hd.errors:
		return err
	}
}

func (hd *tcpMessageHandler) handleRead() {
	errChan := hd.errors
	defer func() {
		if err := recover(); err != nil {
			errChan <- fmt.Errorf("tcpMessageHandler::handleRead panic: %s", err)
		}
	}()
	conn := hd.conn
	msgChan := hd.msgChan
	// 临时缓冲区，用来存储被截断的数据
	var temBuf []byte
	// 读取缓存
	buf := make([]byte, 1024)

	for {
		if hd.status < 0 {
			return
		}
		//conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		n, err := conn.Read(buf)
		if err != nil {
			errChan <- err
			return
		}

		// var pError error
		temBuf, err = hd.unpack(append(temBuf, buf[:n]...), msgChan)
		if err != nil {
			errChan <- err
			return
		}
	}
}

func (hd *tcpMessageHandler) handleSend() {
	errChan := hd.errors
	defer func() {
		if err := recover(); err != nil {
			errChan <- fmt.Errorf("tcpMessageHandler::handleSend panic: %s", err)
		}
	}()
	for {
		msg, ok := <-hd.sendChan
		if !ok {
			return
		}
		if pack, err := hd.pack(msg); err != nil {
			errChan <- err
			return
		} else {
			if _, err := hd.conn.Write(pack); err != nil {
				errChan <- fmt.Errorf("conn write err:%s", err.Error())
				return
			}
		}
	}
}

func (hd *tcpMessageHandler) destruct() {
	hd.statusLock.Lock()
	defer hd.statusLock.Unlock()
	if hd.status == -1 {
		return
	}
	hd.status = -1
	close(hd.msgChan)
	close(hd.sendChan)
}

func (hd *tcpMessageHandler) checkStatusTodo(expect int8, callback func()) bool {
	hd.statusLock.RLock()
	defer hd.statusLock.RUnlock()
	if hd.status != expect {
		return false
	}
	callback()
	return true
}

// 拆包
func (hd *tcpMessageHandler) unpack(buffer []byte, msgChan chan<- *Message) (last []byte, err error) {
	length := uint32(len(buffer))
	var i uint32
	for i = 0; i < length; i++ {
		if length-i < 4 {
			break
		}
		// 获取数据长度
		messageLen := bytesToUint32(buffer[i : i+4])
		if messageLen > hd.packMaxLen {
			return buffer, fmt.Errorf("(%d)len of package of message greater than max (%d)len of package", messageLen, hd.packMaxLen)
		}
		// 数据长度不足
		if length < i+messageLen {
			break
		}
		body := buffer[i+12 : i+messageLen]
		copyBody := make([]byte, len(body))
		copy(copyBody, body)
		msgChan <- NewMessage(bytesToUint32(buffer[i+4:i+8]), bytesToUint32(buffer[i+8:i+12]), copyBody)
		i += messageLen - 1
	}
	if i == length {
		return make([]byte, 0), err
	}
	return buffer[i:], err
}

// 打包数据
func (hd *tcpMessageHandler) pack(msg *Message) (res []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("tcpMessageHandler pack panic: %v", e)
		}
	}()
	buf := new(bytes.Buffer)
	body := msg.Body
	buf.Write(uint32ToBytes(uint32(len(body) + 12))) // 写入总长度
	buf.Write(uint32ToBytes(msg.CId))                // 写入命令ID
	buf.Write(uint32ToBytes(msg.SId))                // 写入流水号
	buf.Write(body)                                  // 写入数据
	return buf.Bytes(), nil
}
