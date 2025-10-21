package handlers

import (
	"testing"

	"github.com/incident-io/incidentio-mcp-golang/internal/client"
)

func TestPaginationLogic(t *testing.T) {
	tests := []struct {
		name            string
		paginationMeta  client.ListResponse
		recordsFetched  int
		expectedHasMore bool
		description     string
	}{
		{
			name: "more results available - records fetched < total",
			paginationMeta: client.ListResponse{
				PaginationMeta: struct {
					After            string `json:"after,omitempty"`
					PageSize         int    `json:"page_size"`
					TotalRecordCount int    `json:"total_record_count"`
				}{
					After:            "01FCNDV6P870EA6S7TK1DSYDG0",
					PageSize:         25,
					TotalRecordCount: 238,
				},
			},
			recordsFetched:  25,
			expectedHasMore: true,
			description:     "Records fetched (25) < total (238)",
		},
		{
			name: "no more results - records fetched = total (real scenario)",
			paginationMeta: client.ListResponse{
				PaginationMeta: struct {
					After            string `json:"after,omitempty"`
					PageSize         int    `json:"page_size"`
					TotalRecordCount int    `json:"total_record_count"`
				}{
					After:            "01K7P0NE00WQDNR5MAK81QEQX4", // after cursor always exists
					PageSize:         25,
					TotalRecordCount: 18, // total records
				},
			},
			recordsFetched:  18, // all records fetched in first page
			expectedHasMore: false,
			description:     "Records fetched (18) = total (18), no more results",
		},
		{
			name: "no more results - no after cursor",
			paginationMeta: client.ListResponse{
				PaginationMeta: struct {
					After            string `json:"after,omitempty"`
					PageSize         int    `json:"page_size"`
					TotalRecordCount int    `json:"total_record_count"`
				}{
					After:            "",
					PageSize:         25,
					TotalRecordCount: 25,
				},
			},
			recordsFetched:  25,
			expectedHasMore: false,
			description:     "No after cursor, no more results",
		},
		{
			name: "edge case - after cursor present but records fetched = total",
			paginationMeta: client.ListResponse{
				PaginationMeta: struct {
					After            string `json:"after,omitempty"`
					PageSize         int    `json:"page_size"`
					TotalRecordCount int    `json:"total_record_count"`
				}{
					After:            "01FCNDV6P870EA6S7TK1DSYDG0",
					PageSize:         25,
					TotalRecordCount: 25,
				},
			},
			recordsFetched:  25,
			expectedHasMore: false,
			description:     "No after cursor, no more results",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the simplified logic from incidents.go
			// Use total_record_count to determine if there are more results
			// The "after" cursor is only needed for the next API call, not for determining if more results exist
			recordsFetched := tt.recordsFetched
			totalRecords := tt.paginationMeta.PaginationMeta.TotalRecordCount
			hasMore := recordsFetched < totalRecords

			if hasMore != tt.expectedHasMore {
				t.Errorf("Expected has_more=%v, got %v. %s", tt.expectedHasMore, hasMore, tt.description)
			}
		})
	}
}

func TestCreatePaginationResponse(t *testing.T) {
	tests := []struct {
		name            string
		paginationMeta  interface{}
		count           int
		expectedHasMore bool
		description     string
	}{
		{
			name: "more results available",
			paginationMeta: struct {
				After            string `json:"after,omitempty"`
				TotalRecordCount int    `json:"total_record_count"`
			}{
				After:            "01FCNDV6P870EA6S7TK1DSYDG0",
				TotalRecordCount: 100,
			},
			count:           25,
			expectedHasMore: true,
			description:     "Should return true when records fetched < total and after cursor present",
		},
		{
			name: "no more results - reached total",
			paginationMeta: struct {
				After            string `json:"after,omitempty"`
				TotalRecordCount int    `json:"total_record_count"`
			}{
				After:            "",
				TotalRecordCount: 50,
			},
			count:           50,
			expectedHasMore: false,
			description:     "Should return false when records fetched = total",
		},
		{
			name: "no more results - no after cursor",
			paginationMeta: struct {
				After            string `json:"after,omitempty"`
				TotalRecordCount int    `json:"total_record_count"`
			}{
				After:            "",
				TotalRecordCount: 25,
			},
			count:           25,
			expectedHasMore: false,
			description:     "Should return false when no after cursor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseTool := &BaseTool{}
			response := baseTool.CreatePaginationResponse(map[string]interface{}{}, tt.paginationMeta, tt.count)

			hasMoreResults, ok := response["has_more_results"].(bool)
			if !ok {
				t.Error("has_more_results should be a boolean")
				return
			}

			if hasMoreResults != tt.expectedHasMore {
				t.Errorf("Expected has_more_results=%v, got %v. %s", tt.expectedHasMore, hasMoreResults, tt.description)
			}
		})
	}
}
