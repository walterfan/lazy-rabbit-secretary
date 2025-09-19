package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/walterfan/lazy-rabbit-secretary/pkg/log"
)

type ChangeItem struct {
	Diff          string `json:"diff"`
	NewPath       string `json:"new_path"`
	OldPath       string `json:"old_path"`
	AMode         string `json:"a_mode"`
	BMode         string `json:"b_mode"`
	NewFile       bool   `json:"new_file"`
	RenamedFile   bool   `json:"renamed_file"`
	DeletedFile   bool   `json:"deleted_file"`
	GeneratedFile bool   `json:"generated_file"`
}

type MergeRequestInfo struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Changes     []ChangeItem `json:"changes"`
}

func formatBoolAsEmoji(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

// GetProjectIDByName gets the numeric project ID from GitLab using the project path (e.g., "namespace/project-name")
func GetProjectIDByName(gitlabURL, projectName, privateToken string) (string, error) {
	logger := log.GetLogger()

	// Construct the search API URL
	searchURL := fmt.Sprintf("%s/api/v4/projects?search=%s", gitlabURL, url.QueryEscape(projectName))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("PRIVATE-TOKEN", privateToken)

	httpClient, err := createHttpClient()
	if err != nil {
		return "", err
	}

	logger.Infof("Fetching project ID from: %s", searchURL)
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Errorf("Error fetching project ID: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitLab project search failed: %s - %s", resp.Status, string(body))
	}

	var projects []struct {
		ID                int    `json:"id"`
		Name              string `json:"name"`
		PathWithNamespace string `json:"path_with_namespace"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return "", err
	}

	if len(projects) == 0 {
		return "", fmt.Errorf("no project found matching name: %s", projectName)
	}

	for _, p := range projects {
		if p.PathWithNamespace == projectName {
			return fmt.Sprintf("%d", p.ID), nil
		}
	}

	// If not exact match, return first result as fallback
	logger.Warnf("No exact match for project %s, returning first match: %s", projectName, projects[0].PathWithNamespace)
	return fmt.Sprintf("%d", projects[0].ID), nil
}

func ChangeItemsToMarkdown(mrInfo MergeRequestInfo) string {
	if len(mrInfo.Changes) == 0 {
		return "No changes found."
	}

	var markdown strings.Builder

	markdown.WriteString("## Merge Request\n\n")
	markdown.WriteString(fmt.Sprintf("* Title: `%s`\n", mrInfo.Title))
	markdown.WriteString(fmt.Sprintf("* Description: `%s`\n", mrInfo.Description))
	markdown.WriteString("\n### changes\n\n")
	for _, change := range mrInfo.Changes {

		markdown.WriteString(fmt.Sprintf("* Old Path: `%s`\n", change.OldPath))
		markdown.WriteString(fmt.Sprintf("* New Path: `%s`\n", change.NewPath))
		markdown.WriteString(fmt.Sprintf("* Added: %s\n", formatBoolAsEmoji(change.NewFile)))
		markdown.WriteString(fmt.Sprintf("* Renamed: %s\n", formatBoolAsEmoji(change.RenamedFile)))
		markdown.WriteString(fmt.Sprintf("* Deleted: %s\n", formatBoolAsEmoji(change.DeletedFile)))

		markdown.WriteString("\n* Code Diff:\n")
		markdown.WriteString("```\n")
		markdown.WriteString(change.Diff)
		markdown.WriteString("\n```\n")
		markdown.WriteString("\n------\n\n")
	}

	return markdown.String()
}

func GetMergeRequestChange(gitlabURL, projectID, mergeRequestID, privateToken string) (string, error) {
	logger := log.GetLogger()
	apiURL := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests/%s/changes",
		gitlabURL,
		url.PathEscape(projectID), // Handles namespace or numeric ID
		url.PathEscape(mergeRequestID))
	// Create the request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("PRIVATE-TOKEN", privateToken)
	// Send request
	httpClient, err := createHttpClient()
	if err != nil {
		return "", err
	}
	logger.Infof("sending request to %s", apiURL)
	resp, err := httpClient.Do(req)
	if err != nil {
		logger.Error("Error sending request to %s:", apiURL, err)
		return "", err
	}
	defer resp.Body.Close()
	logger.Infof("get response %d:", resp.StatusCode)
	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitLab API error: %s - %s", resp.Status, string(body))
	}

	// Parse response JSON
	var result MergeRequestInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return ChangeItemsToMarkdown(result), nil
}

func GetGitLabFileContent(gitlabURL, projectID, filePath, branch, privateToken string) (string, error) {
	// Construct the API URL
	apiURL := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s?ref=%s",
		gitlabURL,
		url.PathEscape(projectID), // Handles namespace or numeric ID
		url.PathEscape(filePath),
		url.PathEscape(branch),
	)

	// Create the request
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("PRIVATE-TOKEN", privateToken)

	// Send request
	httpClient, err := createHttpClient()
	if err != nil {
		return "", err
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitLab API error: %s - %s", resp.Status, string(body))
	}

	// Parse response JSON
	var result struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	// Decode base64 content if needed
	if result.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(result.Content)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}

	return result.Content, nil
}
