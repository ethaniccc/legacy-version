package legacyver

import (
	"github.com/sandertv/gophertunnel/minecraft"
)

// All returns a slice of all legacy protocol versions that are supported. dragonflyMapping
// must be set to true if you're using Dragonfly.
func All() []minecraft.Protocol {
	return []minecraft.Protocol{
		New2168(),
		New1001(),
		New975(),
		New924(),
		New944(),
		New898(),
		New860(),
		New859(),
		New844(),
		New827(),
		New819(),
		New818(),
		New800(),
		New786(),
		New776(),
		New766(),
		New748(),
		New729(),
		New712(),
		New686(),
		New685(),
		New671(),
		New662(),
		New649(),
	}
}
