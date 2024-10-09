package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"taskManager/app/db"

	"github.com/anthdm/superkit/kit"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

// handelGetProjectWithDeadline is an HTTP handler that fetches projects from the database
// based on the provided ID parameter in the URL. It then returns the projects
// as a JSON response.
//
// The function expects the ID parameter to be provided in the URL path, and it
// will attempt to convert the parameter to an integer. If the conversion fails,
// an error is returned.
//
// The function calls the fetchProjectsFromDatabase function to retrieve the
// projects, and then it creates a JSON response with the projects data.
func handelGetProjectWithDeadline(kit *kit.Kit) error {

	ctx := kit.Request.Context()
	id := chi.URLParam(kit.Request, "user_id")
	project := chi.URLParam(kit.Request, "name")
	deadline := chi.URLParam(kit.Request, "deadline")
	userID, err := strconv.Atoi(id)
	fmt.Println(userID)
	if err != nil {
		return kit.JSON(http.StatusUnprocessableEntity, map[string]string{"error": fmt.Sprintf("invalide id : %s", err)})
	}

	projects, err := fetchProjectsFromDatabase(ctx, project, deadline)
	if err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err)})
	}

	if len(projects) == 0 {
		return kit.JSON(http.StatusOK, map[string]string{"info": "No Project found"})
	}

	response := map[string]any{
		"projects": projects,
	}

	return kit.JSON(http.StatusOK, response)
}

func handelGetProject(kit *kit.Kit) error {

	ctx := kit.Request.Context()
	id := chi.URLParam(kit.Request, "user_id")
	project := chi.URLParam(kit.Request, "name")
	userID, err := strconv.Atoi(id)
	fmt.Println(userID)
	if err != nil {
		return kit.JSON(http.StatusUnprocessableEntity, map[string]string{"error": fmt.Sprintf("invalide id : %s", err)})
	}
	projects, err := fetchProjectsFromDatabase(ctx, project, "")
	if err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprint(err)})
	}

	if len(projects) == 0 {
		return kit.JSON(http.StatusOK, map[string]string{"info": "No Project found"})
	}

	response := map[string]any{
		"projects": projects,
	}

	return kit.JSON(http.StatusOK, response)
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
	//fmt.Println("test_request_project")
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
	kit.Response.Header().Set("HX-Trigger", "menuTrigger")
	return kit.JSON(http.StatusCreated, map[string]string{
		"info":       "Project created successfully ",
		"project_id": fmt.Sprintf("%d", project.ProjectID),
		"user_id":    fmt.Sprintf("%d", project.CreatedBy),
	})
}

// fetchProjectsFromDatabase fetches a list of project names from the database based on the provided project name or search term.
// The function takes an integer id, a string project, and a string searchTerm as input parameters.
// It constructs a database query to select project names where the name matches the provided project name or the name or description matches the search term.
// The function returns a slice of strings containing the matching project names.
func fetchProjectsFromDatabase(ctx context.Context, name, deadline string) ([]Project, error) {
	var results []Project

	// Get the database connection and pass the context
	query := db.Get().Table("projects").WithContext(ctx)

	// Apply filtering based on project name or description
	if len(name) != 0 {
		lowerName := strings.ToLower(name)
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", "%"+lowerName+"%", "%"+lowerName+"%")
	}
	// Optionally filter by deadline if provided
	if len(deadline) != 0 {
		query = query.Where("deadline = ?", deadline)
	}
	// Limit the results and apply offset for pagination
	//query = query.Limit(limit).Offset(offset)

	// Execute the query and check for errors
	err := query.Select("project_id,name,project_type, description, created_at,deadline, created_by").Find(&results).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching projects: %v", err)
	}

	return results, nil
}

// handelGetSingleProject is an HTTP handler function that fetches a single project from the database based on the provided name parameter.
// It takes a *kit.Kit as input, which contains the HTTP request context and other request-related information.
// The function extracts the "user_id" and "name" parameters from the request URL, and then calls the fetchProjectFromDatabase function
// to retrieve the project details. If the project is found, it is returned in the HTTP response with a status code of http.StatusCreated.
// If an error occurs during the database fetch, a 500 Internal Server Error response is returned with the error message.
func handelGetSingleProject(kit *kit.Kit) error {
	ctx := kit.Request.Context()
	id := chi.URLParam(kit.Request, "user_id")
	fmt.Println(id)
	name := chi.URLParam(kit.Request, "name")
	project, err := fetchProjectFromDatabase(ctx, name)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return kit.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Failed to fetch project : %s", err)})
		}
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to fetch project : %s", err)})
	}
	return kit.JSON(http.StatusCreated, map[string]any{"projects": project})
}
func fetchProjectFromDatabase(ctx context.Context, name string) (*Project, error) {
	var result Project
	query := db.Get().Table("projects").WithContext(ctx)
	query = query.Where("LOWER(name) = ?", strings.ToLower(name))
	err := query.Select("project_id,name,project_type, description, created_at,deadline, created_by").First(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// createProjectsFromDatabase creates a new project in the database if it doesn't already exist.
// It takes a pointer to a Project struct as input, and returns the number of rows affected
// and any error that occurred during the database operation.
func createProjectsFromDatabase(project *Project) (int, error) {

	result := db.Get().Where(Project{Name: project.Name}).FirstOrCreate(project)
	return int(result.RowsAffected), result.Error
}

func handelUpdateProjects(kit *kit.Kit) error {
	ctx := kit.Request.Context()
	userid := chi.URLParam(kit.Request, "user_id")
	fmt.Println(userid)
	var project Project
	err := json.NewDecoder(kit.Request.Body).Decode(&project)
	if err != nil {
		return kit.JSON(http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Invalid request payload : %s", err)})
	}
	err = updateProjectsFromDatabase(&project, ctx)

	if err == gorm.ErrRecordNotFound {
		return kit.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Failed to fetch project : %s", err)})
	}
	if err != nil {
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to update project : %s", err)})
	}

	return kit.JSON(http.StatusOK, map[string]string{"info": "Project updated successfully"})
}

// updateProjectsFromDatabase updates an existing project in the database.
// It takes a pointer to a Project struct as input, and updates the
// TypeProject, Description, Deadline, and CreatedBy fields of the project.
// If the update is successful, it returns nil, otherwise it returns the error.
func updateProjectsFromDatabase(project *Project, ctx context.Context) error {
	result := db.Get().Model(project).Updates(Project{
		ProjectType: project.ProjectType,
		Description: project.Description,
		Deadline:    project.Deadline,
		CreatedBy:   project.CreatedBy,
	}).WithContext(ctx)
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

func handelDeleteProjects(kit *kit.Kit) error {
	ctx := kit.Request.Context()
	userID := chi.URLParam(kit.Request, "user_id")
	fmt.Println(userID)
	projectName := chi.URLParam(kit.Request, "name")
	fmt.Println(projectName)
	result := db.Get().Where("name = ?", projectName).Delete(&Project{}).WithContext(ctx)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return kit.JSON(http.StatusNotFound, map[string]string{"error": fmt.Sprintf("Failed to delete project : %s", result.Error)})
		}
		return kit.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to delete project : %s", result.Error)})
	}
	kit.Response.Header().Set("HX-Trigger", "menuTrigger")
	return kit.JSON(http.StatusOK, map[string]string{"info": "Project deleted successfully"})
}
