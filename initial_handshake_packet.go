// Copyright 2015 The mp Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"log"
)

type InitialHandshakePacket struct {
	PacketHeader
	protocolVersion     uint8
	serverVersion       []byte // human-readable server version
	ServerVersion       string
	connectionID        [4]byte
	ConnectionID        uint32
	authPluginDataPart1 [8]byte // [len=8] first 8 bytes of the auth-plugin data
	AuthPluginDataPart1 string
	filler1             uint8
	capabilityFlag1     [2]byte // lower 2 bytes of the Protocol::CapabilityFlags (optional)
	characterSet        uint8   // default server character-set, only the lower 8-bits Protocol::CharacterSet (optional)
	statusFlags         [2]byte // Protocol::StatusFlags (optional)
	capabilityFlag2     [2]byte // upper 2 bytes of the Protocol::CapabilityFlags
	authPluginDataLen   uint8   // length of the combined auth_plugin_data, if auth_plugin_data_len is > 0
	authPluginName      []byte  // name of the auth_method that the auth_plugin_data belongs to
}

func (p *InitialHandshakePacket) Read(c *Client) error {
	err := p.PacketHeader.Read(c)
	if err != nil {
		return err
	}

	data := make([]byte, p.PayloadLength)
	_, err = io.ReadFull(c.conn, data)
	if err != nil {
		return err
	}

	log.Printf("%#v\n", data)

	return nil
}
