package communication

// Communicator 消息通讯接口
type Communicator interface {
	Connect() error
	Send(message *Message) error
	Recv() <-chan *Message
	Close() error
}

// Message 消息结构
type Message struct {
	TLen uint32 // 总数据长度
	CId  uint32 // 命令ID
	SId  uint32 // 消息流水号
	Body []byte // 消息体
}

func NewMessage(cid uint32, sid uint32, body []byte) *Message {
	return &Message{
		TLen: uint32(len(body) + 12),
		CId:  cid,
		SId:  sid,
		Body: body,
	}
}
