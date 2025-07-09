package hetzner

import (
	"time"
)

// StorageBox represents a storagebox record prepared for the exporter.
type StorageBox struct {
	ID       int64             `json:"id"`
	Name     string            `json:"name"`
	Server   string            `json:"server"`
	System   string            `json:"system"`
	Username string            `json:"username"`
	Status   string            `json:"status"`
	Created  time.Time         `json:"created"`
	Labels   map[string]string `json:"labels"`
	Type     struct {
		ID                     int64                `json:"id"`
		Name                   string               `json:"name"`
		Description            string               `json:"description"`
		Size                   float64              `json:"size"`
		SubaccountsLimit       int                  `json:"subaccounts_limit"`
		SnapshotLimit          int                  `json:"snapshot_limit"`
		AutomaticSnapshotLimit int                  `json:"automatic_snapshot_limit"`
		Deprecation            map[string]time.Time `json:"deprecation"`
		Prices                 []struct {
			Location string `json:"location"`
			Hourly   struct {
				Gross string `json:"gross"`
				Net   string `json:"net"`
			} `json:"price_hourly"`
			Monthly struct {
				Gross string `json:"gross"`
				Net   string `json:"net"`
			} `json:"price_monthly"`
			Setup struct {
				Gross string `json:"gross"`
				Net   string `json:"net"`
			} `json:"setup_fee"`
		} `json:"prices"`
	} `json:"storage_box_type"`
	Stats struct {
		Size      int64 `json:"size"`
		Data      int64 `json:"size_data"`
		Snapshots int64 `json:"size_snapshots"`
	} `json:"stats"`
	Protection struct {
		Delete bool `json:"delete"`
	} `json:"protection"`
	Location struct {
		ID          int64   `json:"id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		City        string  `json:"city"`
		Country     string  `json:"country"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		NetworkZone string  `json:"network_zone"`
	} `json:"location"`
	AccessSettings struct {
		ReachableExternally bool `json:"reachable_externally"`
		SambaEnabled        bool `json:"samba_enabled"`
		SSHEnabled          bool `json:"ssh_enabled"`
		WebDAVEnabled       bool `json:"webdav_enabled"`
		ZFSEnabled          bool `json:"zfs_enabled"`
	} `json:"access_settings"`
	SnapshotPlan struct {
		DayOfMonth   int `json:"day_of_month"`
		DayOfWeek    int `json:"day_of_week"`
		Hour         int `json:"hour"`
		Minute       int `json:"minute"`
		MaxSnapshots int `json:"max_snapshots"`
	} `json:"snapshot_plan"`
}

type storageBoxAllResponse struct {
	Meta struct {
		Pagination struct {
			Page         int `json:"page"`
			PerPage      int `json:"per_page"`
			PreviousPage int `json:"previous_page"`
			NextPage     int `json:"next_page"`
			LastPage     int `json:"last_page"`
			TotalEntries int `json:"total_entries"`
		} `json:"pagination"`
	} `json:"meta"`
	StorageBoxes []StorageBox `json:"storage_boxes"`
}
