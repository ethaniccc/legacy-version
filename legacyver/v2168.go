package legacyver

import "github.com/akmalfairuz/legacy-version/legacyver/proto"

// New2168 returns the protocol used by Minecraft 1.26.40.
// Protocol 2168 is the same ID as 1.26.44 (the current gophertunnel
// version), but the wire format has small differences that the
// legacyver conversion functions handle via IsProtoGTE(r, ID2168).
func New2168() *Protocol {
	return &Protocol{
		ver: "1.26.40",
		id:  proto.ID2168,
	}
}
