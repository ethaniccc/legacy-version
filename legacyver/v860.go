package legacyver

import (
	"github.com/ethaniccc/legacy-version/legacyver/proto"
)

func New860(dragonflyMapping bool) *Protocol {
	p := New859(dragonflyMapping)
	p.ver = "1.21.124"
	p.id = proto.ID860
	return p
}
