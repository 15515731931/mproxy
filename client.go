// Copyright 2015 The mp Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package main

import (
	"net"
)

type Client struct {
	conn     net.Conn
	username string
	password string
	database string
}

func (c *Client) Connet(network, address string) error {
	var err error
	c.conn, err = net.Dial(network, address)
	if err != nil {
		return err
	}

	ihp := new(InitialHandshakePacket)
	err = ihp.Read(c)
	if err != nil {
		return err
	}

	return nil
}
