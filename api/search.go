package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// SearchOptions contains options for JQL search
type SearchOptions struct {
	JQL        string
	StartAt    int
	MaxResults int
	Fields     []string
}

// Search searches for issues using JQL
func (c *Client) Search(opts SearchOptions) (*SearchResult, error) {
	params := map[string]string{
		"jql": opts.JQL,
	}

	if opts.StartAt > 0 {
		params["startAt"] = strconv.Itoa(opts.StartAt)
	}

	if opts.MaxResults > 0 {
		params["maxResults"] = strconv.Itoa(opts.MaxResults)
	} else {
		params["maxResults"] = "50"
	}

	if len(opts.Fields) > 0 {
		fields := ""
		for i, f := range opts.Fields {
			if i > 0 {
				fields += ","
			}
			fields += f
		}
		params["fields"] = fields
	}

	urlStr := buildURL(fmt.Sprintf("%s/search", c.BaseURL), params)
	body, err := c.get(urlStr)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	return &result, nil
}

// SearchAll searches for all issues matching JQL (handles pagination)
func (c *Client) SearchAll(jql string, maxResults int) ([]Issue, error) {
	if maxResults <= 0 {
		maxResults = 1000
	}

	var allIssues []Issue
	startAt := 0
	pageSize := 100

	for {
		result, err := c.Search(SearchOptions{
			JQL:        jql,
			StartAt:    startAt,
			MaxResults: pageSize,
		})
		if err != nil {
			return nil, err
		}

		allIssues = append(allIssues, result.Issues...)

		if len(allIssues) >= result.Total || len(allIssues) >= maxResults {
			break
		}

		startAt += len(result.Issues)
		if len(result.Issues) == 0 {
			break
		}
	}

	if len(allIssues) > maxResults {
		allIssues = allIssues[:maxResults]
	}

	return allIssues, nil
}
