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

	protocolVersion byte
	ProtocolVersion uint8

	serverVersion []byte // human-readable server version
	ServerVersion string

	connectionID []byte // 4 byte Protocol::FixedLengthInteger
	ConnectionID uint32

	authPluginDataPart1 []byte // [len=8] first 8 bytes of the auth-plugin data

	filler1 byte

	capabilityFlag1 []byte // lower 2 bytes of the Protocol::CapabilityFlags (optional)

	characterSet byte // default server character-set, only the lower 8-bits Protocol::CharacterSet (optional)
	CharacterSet uint8

	statusFlags []byte // 2 byte Protocol::StatusFlags (optional)

	capabilityFlag2 []byte // upper 2 bytes of the Protocol::CapabilityFlags

	authPluginDataLen byte // length of the combined auth_plugin_data, if auth_plugin_data_len is > 0

	reserved []byte // 10 byte reserved (all [00])

	authPluginDataPart2 []byte

	authPluginName []byte // name of the auth_method that the auth_plugin_data belongs to

	authPluginData []byte

	capabilityFlag []byte
	CapabilityFlag uint32
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

	// 4a 00 00 00 0a 35 2e 35    2e 33 34 00 15 00 00 00    J....5.5.34.....
	// 34 28 69 59 65 6d 72 77    00 ff f7 08 02 00 0f 80    4(iYemrw........
	// 15 00 00 00 00 00 00 00    00 00 00 21 3b 43 7c 6c    ...........!;C|l
	// 26 77 45 34 74 6b 3c 00    6d 79 73 71 6c 5f 6e 61    &wE4tk<.mysql_na
	// 74 69 76 65 5f 70 61 73    73 77 6f 72 64 00          tive_password.

	r := NewPacketReader(data)

	p.protocolVersion = r.Byte()
	p.ProtocolVersion = uint8(p.protocolVersion)

	p.serverVersion = r.BytesNul()
	p.ServerVersion = string(p.serverVersion)

	p.connectionID = r.Bytes(4)
	p.ConnectionID = uint32(p.connectionID[0]) | uint32(p.connectionID[1])<<8 | uint32(p.connectionID[2])<<16 | uint32(p.connectionID[3])<<24

	p.authPluginDataPart1 = r.Bytes(8)
	p.authPluginData = make([]byte, 0)
	p.authPluginData = append(p.authPluginData, p.authPluginDataPart1...)

	p.filler1 = r.Byte()

	p.capabilityFlag1 = r.Bytes(2)
	p.capabilityFlag = make([]byte, 0)
	p.capabilityFlag = append(p.capabilityFlag, p.capabilityFlag1...)

	if len(data) <= r.Index() {
		p.CapabilityFlag = uint32(p.capabilityFlag[0]) | uint32(p.capabilityFlag[1])<<8

		log.Printf("%#v\n", p)
		return nil
	}

	p.characterSet = r.Byte()

	p.statusFlags = r.Bytes(2)

	p.capabilityFlag2 = r.Bytes(2)
	p.capabilityFlag = append(p.capabilityFlag, p.capabilityFlag2...)

	p.authPluginDataLen = r.Byte()
	p.reserved = r.Bytes(10)
	p.authPluginDataPart2 = r.Bytes(12)
	_ = r.Byte()

	p.authPluginName = r.BytesNul()

	p.authPluginData = append(p.authPluginData, p.authPluginDataPart2...)

	p.CapabilityFlag = uint32(p.capabilityFlag[0]) | uint32(p.capabilityFlag[1])<<8 | uint32(p.capabilityFlag[2])<<16 | uint32(p.capabilityFlag[3])<<24

	log.Printf("%#v\n", p)
	return nil
}
