package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"taskManager/app/db"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
)

// handelGetProject is an HTTP handler that fetches projects from the database
// based on the provided ID parameter in the URL. It then returns the projects
// as a JSON response.
//
// The function expects the ID parameter to be provided in the URL path, and it
// will attempt to convert the parameter to an integer. If the conversion fails,
// an error is returned.
//
// The function calls the fetchProjectsFromDatabase function to retrieve the
// projects, and then it creates a JSON response with the projects data.
func handelGetProject(kit *kit.Kit) error {
	id := chi.URLParam(kit.Request, "user_id")
	project := chi.URLParam(kit.Request, "name")
	searchTerm := chi.URLParam(kit.Request, "searchterm")

	userID, err := strconv.Atoi(id)

	if err != nil {
		return kit.JSON(http.StatusUnprocessableEntity, map[string]string{"error": fmt.Sprintf("invalide id : %s", err)})
	}

	projects := fetchProjectsFromDatabase(userID, project, searchTerm)

	if len(projects) == 0 {
		return kit.JSON(http.StatusOK, map[string]string{"info": "No Project found"})
	}

	response := map[string]any{
		"projects": projects,
	}

	err = kit.JSON(http.StatusOK, response)
	if err != nil {
		return err
	}
	return nil
}

// handlePostProjects is an HTTP handler that creates a new project in the database.
//
// It expects the request body to contain a JSON-encoded Project struct. If the
// decoding of the request body fails, it returns a 400 Bad Request response.
//
// If the creation of the project in the database fails, it returns a 500 Internal
// Server Error response with the error message.
//
// If the project already exists in the database, it returns a 500 Internal Server
// Error response with an informational message.
//
// If the project is created successfully, it returns a 200 OK response with the
// project ID and the user ID who created the project.
func handlePostProjects(kit *kit.Kit) error {
	var project Project
	err := json.NewDecoder(kit.Request.Body).Decode(&project)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request payload : %s", err)})
	}

	row, err := createProjectsFromDatabase(&project)
	if err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create project : %s", err)})
	}

	if row == 0 {
		return kit.JSON(http.StatusOK, map[string]string{"info": fmt.Sprintf("Project %s all ready exist", project.Name)})
	}

	return kit.JSON(http.StatusOK, map[string]string{
		"info":       "Project created successfully ",
		"project_id": fmt.Sprintf("%d", project.ProjectID),
		"user_id":    fmt.Sprintf("%d", project.CreatedBy),
	})
}

// fetchProjectsFromDatabase fetches a list of project names from the database based on the provided project name or search term.
// The function takes an integer id, a string project, and a string searchTerm as input parameters.
// It constructs a database query to select project names where the name matches the provided project name or the name or description matches the search term.
// The function returns a slice of strings containing the matching project names.
func fetchProjectsFromDatabase(id int, project, searchTerm string) []string {
	var results []string

	// Use a database query to fetch projects based on the project name or search term
	// This is a placeholder query, adjust it according to your database schema and ORM
	query := db.Get().Table("projects")

	if project != "" {
		query = query.Where("name LIKE ?", "%"+project+"%")
	}

	if searchTerm != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+searchTerm+"%", "%"+searchTerm+"%")
	}

	query.Select("name").Find(&results)

	fmt.Printf("id = %d\nproject = %s\nsearch-term = %s\n", id, project, searchTerm)

	return results
}

// createProjectsFromDatabase creates a new project in the database if it doesn't already exist.
// It takes a pointer to a Project struct as input, and returns the number of rows affected
// and any error that occurred during the database operation.
func createProjectsFromDatabase(project *Project) (int, error) {

	result := db.Get().Where(Project{Name: project.Name}).FirstOrCreate(project)
	return int(result.RowsAffected), result.Error
}

func handelUpdateProjects(kit *kit.Kit) error {
	var project Project
	err := json.NewDecoder(kit.Request.Body).Decode(&project)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request payload : %s", err)})
	}
	err = updateProjectsFromDatabase(&project)

	if err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to update project : %s", err)})
	}

	return kit.JSON(http.StatusOK, map[string]string{
		"info":       "Project updated successfully ",
		"project_id": fmt.Sprintf("%d", project.ProjectID),
		"user_id":    fmt.Sprintf("%d", project.CreatedBy),
	})
}

// updateProjectsFromDatabase updates an existing project in the database.
// It takes a pointer to a Project struct as input, and updates the
// TypeProject, Description, Deadline, and CreatedBy fields of the project.
// If the update is successful, it returns nil, otherwise it returns the error.
func updateProjectsFromDatabase(project *Project) error {
	result := db.Get().Model(project).Updates(Project{
		TypeProject: project.TypeProject,
		Description: project.Description,
		Deadline:    project.Deadline,
		CreatedBy:   project.CreatedBy,
	})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func handlerGetAllProjects(kit *kit.Kit) error {

	projects, err := getAllProjectsFromDatabase()
	if err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to fetch projects: %s", err)})
	}

	return kit.JSON(http.StatusOK, map[string]any{
		"projects": projects,
	})

}

func getAllProjectsFromDatabase() ([]Project, error) {
	var projects []Project
	result := db.Get().Find(&projects)
	if result.Error != nil {
		return nil, result.Error
	}
	return projects, nil
}
