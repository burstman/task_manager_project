package handlers

import (
	"taskManager/taskManagerLayout"

	"github.com/anthdm/superkit/kit"
)

func HandelOverView(kit *kit.Kit) error {
	data := map[string]any{
		"progress": 66,
	}
	return kit.Render(taskManagerLayout.Overview(data))
}
