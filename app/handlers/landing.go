package handlers

import (
	"taskManager/app/views/landing"

	"github.com/anthdm/superkit/kit"
)

func HandleLandingIndex(kit *kit.Kit) error {
	return kit.Render(landing.Index())
}
