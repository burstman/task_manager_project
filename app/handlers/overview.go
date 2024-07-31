package handlers

import (
	"taskManager/overview"

	"github.com/anthdm/superkit/kit"
)

func HandelOverView(kit *kit.Kit) error {
	data := map[string]any{
		"progress": 66,
	}
	return kit.Render(overview.Overview(data))
}
