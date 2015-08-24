// Copyright 2015 The mp Authors. All rights reserved.
// Use of this source code is governed by a GNU-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
)

type PacketReader struct {
	data  []byte
	pLow  int
	pHigh int
}

func (p *PacketReader) Index() int {
	return p.pLow
}

func (p *PacketReader) Byte() byte {
	d := p.data[p.pLow]

	p.pLow++
	return d
}

func (p *PacketReader) Bytes(l int) []byte {
	d := make([]byte, l)
	for i := 0; i < l; i++ {
		d[i] = p.data[p.pLow]
		p.pLow++
	}

	return d
}

func (p *PacketReader) BytesNul() []byte {
	p.pHigh = p.pLow + bytes.IndexByte(p.data[p.pLow:], 0x00)
	d := p.data[p.pLow:p.pHigh]
	p.pLow = p.pHigh

	p.pLow++
	return d
}

func (p *PacketReader) Bytes4() [4]byte {
	var d [4]byte
	for i := 0; i < 4; i++ {
		d[i] = p.data[p.pLow]
		p.pLow++
	}

	return d
}

func (p *PacketReader) Uint8() uint8 {
	return uint8(p.Byte())
}

func (p *PacketReader) StringNul() string {
	return string(p.BytesNul())
}

func NewPacketReader(data []byte) *PacketReader {
	return &PacketReader{data: data, pLow: 0, pHigh: 0}
}
