package handlers

import (
	"taskManager/taskManagerLayout"

	"github.com/anthdm/superkit/kit"
)

//const activatedContent = "activatedContent"

var ActiveContent string

var percentTodayTasks = map[string]any{
	"progress": 66,
}
var processData = map[string]any{
	"progress": 90,
}

func HandelerBaseOverView(kit *kit.Kit) error {

	return kit.Render(taskManagerLayout.BaseOverView(percentTodayTasks, processData, "overview"))

}

func HandelOverview(kit *kit.Kit) error {
	return kit.Render(taskManagerLayout.Overview(percentTodayTasks, processData))
}

func HandelBoard(kit *kit.Kit) error {
	return kit.Render(taskManagerLayout.Board(processData))
}

func HandelMembers(kit *kit.Kit) error {

	return kit.Render(taskManagerLayout.Members(processData))
}

func HandelTimeline(kit *kit.Kit) error {

	return kit.Render(taskManagerLayout.Timeline(processData))
}

func HandelReport(kit *kit.Kit) error {

	return kit.Render(taskManagerLayout.Report(processData))
}

func HandelFiles(kit *kit.Kit) error {

	return kit.Render(taskManagerLayout.Files(processData))
}
