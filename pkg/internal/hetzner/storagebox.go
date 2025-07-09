package hetzner

import (
	"context"
	"fmt"
)

// StorageBoxClient is a client for the storageboxes API.
type StorageBoxClient struct {
	client *Client
}

// All fetches all available storageboxes from API.
func (c *StorageBoxClient) All(ctx context.Context) ([]*StorageBox, error) {
	records := make([]*StorageBox, 0)
	page := 1

	for {
		req, err := c.client.NewRequest(
			ctx,
			"GET",
			fmt.Sprintf("%s/storage_boxes?per_page=50&page=%d", Endpoint, page),
			nil,
		)

		if err != nil {
			return nil, err
		}

		result := &storageBoxAllResponse{}

		if _, err := c.client.Do(req, result); err != nil {
			return nil, err
		}

		for _, record := range result.StorageBoxes {
			records = append(
				records,
				&record,
			)
		}

		if result.Meta.Pagination.NextPage < 1 {
			break
		}

		page = result.Meta.Pagination.NextPage
	}

	return records, nil
}
