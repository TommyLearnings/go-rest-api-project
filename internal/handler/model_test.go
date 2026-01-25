package handler_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/TommyLearning/go-rest-api-project/internal/news"

	"github.com/TommyLearning/go-rest-api-project/internal/handler"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewsPostReqBody_Validate(t *testing.T) {
	type expectaions struct {
		err  string
		news *news.Record
	}
	testCases := []struct {
		name        string
		req         handler.NewsPostReqBody
		expectaions expectaions
	}{
		{
			name: "author empty",
			req:  handler.NewsPostReqBody{},
			expectaions: expectaions{
				err: "author is empty",
			},
		},
		{
			name: "title empty",
			req: handler.NewsPostReqBody{
				Author: "test-author",
			},
			expectaions: expectaions{
				err: "title is empty",
			},
		},
		{
			name: "content empty",
			req: handler.NewsPostReqBody{
				Author: "test-author",
				Title:  "test-title",
			},
			expectaions: expectaions{
				err: "content is empty",
			},
		},
		{
			name: "summary empty",
			req: handler.NewsPostReqBody{
				Author: "test-author",
				Title:  "test-title",
			},
			expectaions: expectaions{
				err: "summary is empty",
			},
		},
		{
			name: "time invalid",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "invalid-time",
			},
			expectaions: expectaions{
				err: `parsing time "invalid-time"`,
			},
		},
		{
			name: "source invalid",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
			},
			expectaions: expectaions{
				err: "source is empty",
			},
		},
		{
			name: "tags empty",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
				Source:    "https://google.com",
			},
			expectaions: expectaions{
				err: "tags cannot be empty",
			},
		},
		{
			name: "validate",
			req: handler.NewsPostReqBody{
				Author:    "test-author",
				Title:     "test-title",
				Content:   "test-content",
				Summary:   "test-summary",
				CreatedAt: "2024-04-07T05:13:27+00:00",
				Source:    "https://google.com",
				Tags:      []string{"tag1", "tag2"},
			},
			expectaions: expectaions{
				news: &news.Record{
					Author:  "test-author",
					Title:   "test-title",
					Content: "test-content",
					Summary: "test-summary",
					Tags:    []string{"tag1", "tag2"},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			news, err := tc.req.Validate()

			if tc.expectaions.err != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectaions.err)
			} else {
				assert.NoError(t, err)

				parseTime, paresErr := time.Parse(time.RFC3339, tc.req.CreatedAt)
				require.NoError(t, paresErr)
				tc.expectaions.news.CreatedAt = parseTime

				parseSource, parseErr := url.Parse(tc.req.Source)
				require.NoError(t, parseErr)
				tc.expectaions.news.Source = parseSource.String()

				assert.Equal(t, tc.expectaions.news, news)
			}

		})
	}
}
