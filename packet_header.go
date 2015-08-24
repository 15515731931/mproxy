// Copyright 2015 The mp Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"io"
	"log"
)

type PacketHeader struct {
	payloadLength [3]byte
	PayloadLength uint32
	sequenceID    byte
	SequenceID    uint8
}

func (p *PacketHeader) Read(c *Client) error {
	header := make([]byte, 4)
	_, err := io.ReadFull(c.conn, header)
	if err != nil {
		return err
	}

	log.Printf("%#v\n", header)

	p.payloadLength[0] = header[0]
	p.payloadLength[1] = header[1]
	p.payloadLength[2] = header[2]
	p.PayloadLength = uint32(header[0]) | uint32(header[1])<<8 | uint32(header[2])<<16
	if p.PayloadLength < 1 {
		return errors.New("Invalid payload length")
	}

	//FIXME: Check p.sequenceID
	p.sequenceID = header[3]
	p.SequenceID = uint8(header[3])

	log.Printf("%#v\n", p)

	return nil
}
