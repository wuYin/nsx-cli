package client

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"nsx/codec"
	"strconv"
	"tron"
)

var (
	CRLF       = []byte{'\r', '\n'}
	DELM       = []byte{'*'}
	GET_HEADER = []byte("*2\r\n$3\r\nGET\r\n")
)

type ClientCodec struct{}

func NewClientCodec() *ClientCodec {
	return &ClientCodec{}
}

func (c *ClientCodec) ReadPacket(r *bufio.Reader) ([]byte, error) {

	readUnit := func(r *bufio.Reader) ([]byte, error) {
		head, err := ReadFullLine(r)
		if err != nil {
			return nil, fmt.Errorf("packet: read fail: %v", err)
		}
		if len(head) < 2 {
			return nil, fmt.Errorf("packet: invalid $ head: %q", head)
		}
		dataLen, err := strconv.Atoi(string(head[1:]))
		if err != nil {
			return nil, fmt.Errorf("packet: atoi failed: %v", err)
		}
		data, err := ReadFullLine(r)
		if err != nil {
			return nil, fmt.Errorf("packet: read fail: %v", err)
		}
		// 校验一下
		if dataLen != len(data) {
			return nil, fmt.Errorf("packet: unmatch data length, head %d, data %d", dataLen, len(data))
		}
		return data, nil
	}

	packData, err := readUnit(r) // 读取数据部分
	if err != nil {
		return nil, err
	}

	return packData, nil
}

// 直接将调用返回的数据封装到 packet 中
func (c *ClientCodec) UnmarshalPacket(buf []byte) (*tron.Packet, error) {
	var callResp codec.CallResp
	if err := json.Unmarshal(buf, &callResp); err != nil {
		panic(err)
	}
	p := tron.NewRespPacket(callResp.Seq, buf)
	return p, nil
}

// 序列化
// 封装通信协议底层的数据部分
func (c *ClientCodec) MarshalPacket(p tron.Packet) []byte {
	var cmdReq codec.CmdReq
	if err := json.Unmarshal(p.Data, &cmdReq); err != nil {
		panic(err)
	}
	cmdReq.Seq = p.Seq()
	p.Data, _ = json.Marshal(cmdReq)

	n := len(p.Data)
	buf := bytes.NewBuffer(make([]byte, 0, packBufLen(n)))
	buf.Write(GET_HEADER)
	buf.WriteByte('$')
	buf.WriteString(fmt.Sprintf("%d", n))
	buf.Write(CRLF)
	buf.Write(p.Data)
	buf.Write(CRLF)
	return buf.Bytes()
}

// $N CR LF [DATA] CR LF
func packBufLen(dataLen int) int {
	return 1 + len(fmt.Sprintf("%d", dataLen)) + len(CRLF) + dataLen + len(CRLF)
}

// 读取完整的一行
func ReadFullLine(r *bufio.Reader) ([]byte, error) {
	var full []byte
	for {
		l, isPrefix, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		full = append(full, l...)
		if !isPrefix {
			break
		}
	}
	return full, nil
}
