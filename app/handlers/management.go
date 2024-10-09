package handlers

import (
	"fmt"
	"taskManager/plugins/auth"
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

func HandelerBase(kit *kit.Kit) error {

	session := kit.GetSession(auth.UserSessionName)

	firstName, ok := session.Values["firstName"].(string)
	if !ok {
		return fmt.Errorf("error first name values")
	}
	userID, ok := session.Values["userid"].(uint)
	if !ok {
		return fmt.Errorf("error first name values")
	}

	fmt.Printf("%v:%v\n", firstName, userID)

	return kit.Render(taskManagerLayout.BaseOverView(percentTodayTasks, processData, "overview", firstName, userID))
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

func HandelProjectList(kit *kit.Kit) error {
	projects := taskManagerLayout.GetAllProjectNames()

	return kit.Render(taskManagerLayout.ProjectList(projects))
}
