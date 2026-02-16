package proto

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	ID924 = 924 // v1.26.0
	ID898 = 898 // v1.21.130
	ID860 = 860 // v1.21.124
	ID859 = 859 // v1.21.120
	ID844 = 844 // v1.21.110
	ID827 = 827 // v1.21.100
	ID819 = 819 // v1.21.93
	ID818 = 818 // v1.21.90
	ID800 = 800 // v1.21.80
	ID786 = 786 // v1.21.70
	ID776 = 776 // v1.21.60
	ID766 = 766 // v1.21.50
	ID748 = 748 // v1.21.40
	ID729 = 729 // v1.21.30
	ID712 = 712 // v1.21.20
	ID686 = 686 // v1.21.2
	ID685 = 685 // v1.21.0
	ID671 = 671 // v1.20.80
	ID662 = 662 // v1.20.70
	ID649 = 649 // v1.20.60
)

func IsProtoGTE(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() >= proto
}

func IsProtoLTE(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() <= proto
}

func IsProtoLT(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() < proto
}

func IsProtoGT(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() > proto
}

func IsProto(io protocol.IO, proto int32) bool {
	return io.(IO).ProtocolID() == proto
}

func FetchProtoID(io protocol.IO) int32 {
	return io.(IO).ProtocolID()
}
