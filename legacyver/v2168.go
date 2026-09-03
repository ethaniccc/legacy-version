package legacyver

import (
	"github.com/akmalfairuz/legacy-version/legacyver/proto"
)

// New2168 returns the protocol for 1.26.40 through 1.26.43. These share the
// protocol number 2168 with 1.26.44, which gophertunnel implements natively,
// so the listener picks this entry by the client's game version instead of
// by number. The only wire difference to 1.26.44 is the objective of a
// SetScore remove entry, which is a single optional here and a double
// optional in 1.26.44 (see proto.MarshalScoreboardEntry).
func New2168() *Protocol {
	return &Protocol{
		ver: "1.26.40",
		id:  proto.ID2168,
	}
}
