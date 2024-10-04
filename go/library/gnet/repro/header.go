package gnet

import (
	"encoding/binary"
)

const ProtocolHeaderSize = 4 + 4

type (
	ProtocolHeader struct {
		PacketLength uint32
		Sequence     uint32
	}
)

func NewProtocolHeader(packetLength uint32, sequence uint32) ProtocolHeader {
	return ProtocolHeader{
		PacketLength: packetLength,
		Sequence:     sequence,
	}
}

func (h ProtocolHeader) ToBytes() []byte {
	bz := make([]byte, ProtocolHeaderSize)
	binary.BigEndian.PutUint32(bz, h.PacketLength)
	binary.BigEndian.PutUint32(bz[4:], h.Sequence)
	return bz
}

func (h ProtocolHeader) FromBytes(src []byte) ProtocolHeader {
	h.PacketLength = binary.BigEndian.Uint32(src)
	h.Sequence = binary.BigEndian.Uint32(src[4:])
	return h
}

func (h ProtocolHeader) FromPeek(src []byte, err error) (ProtocolHeader, error) {
	if err != nil {
		return h, err
	}
	return h.FromBytes(src), nil
}
