package sharkd

import (
	"encoding/json"
	"net"
)

type SharkdClient struct {
	sockpath  string
	conn      net.Conn
	connected bool
}

func NewSharkdClient() *SharkdClient {
	return &SharkdClient{
		sockpath:  sockpath,
		connected: false,
	}
}
func (c *SharkdClient) Connect() error {
	conn, err := net.Dial("unix", c.sockpath)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *SharkdClient) Send(data interface{}) error {
	if !c.connected {
		if err := c.Connect(); err != nil {
			return err
		}
	}
	bytedata, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(bytedata)
	if err != nil {
		return err
	}
	return nil
}
