package notifydata

import (
	"time"

	"github.com/kopia/kopia/snapshot"
)

const durationPrecision = 100 * time.Millisecond

// ManifestWithError represents information about the snapshot manifest with optional error.
type ManifestWithError struct {
	Manifest snapshot.Manifest  `json:"manifest"` // may not be filled out if there was an error, Manifst.Source is always set.
	Previous *snapshot.Manifest `json:"previous"` // may not be filled out

	Error string `json:"error"` // will be present if there was an error
}

// StartTimestamp returns the start time of the snapshot.
func (m *ManifestWithError) StartTimestamp() time.Time {
	return m.Manifest.StartTime.ToTime().UTC().Truncate(time.Second)
}

// EndTimestamp returns the end time of the snapshot.
func (m *ManifestWithError) EndTimestamp() time.Time {
	return m.Manifest.EndTime.ToTime().UTC().Truncate(time.Second)
}

// TotalSize returns the total size of the snapshot in bytes.
func (m *ManifestWithError) TotalSize() int64 {
	if m.Manifest.RootEntry == nil {
		return 0
	}

	if m.Manifest.RootEntry.DirSummary != nil {
		return m.Manifest.RootEntry.DirSummary.TotalFileSize
	}

	return m.Manifest.RootEntry.FileSize
}

// TotalSizeDelta returns the total size of the snapshot in bytes.
func (m *ManifestWithError) TotalSizeDelta() int64 {
	if m.Previous == nil {
		return 0
	}

	if m.Manifest.RootEntry == nil {
		return 0
	}

	if m.Manifest.RootEntry.DirSummary != nil && m.Previous.RootEntry.DirSummary != nil {
		return m.Manifest.RootEntry.DirSummary.TotalFileSize - m.Previous.RootEntry.DirSummary.TotalFileSize
	}

	return m.Manifest.RootEntry.FileSize
}

// New returns the total size of the snapshot in bytes.
func (m *ManifestWithError) New() int64 {
	if m.Manifest.RootEntry == nil {
		return 0
	}

	if m.Manifest.RootEntry.DirSummary != nil {
		return m.Manifest.RootEntry.DirSummary.TotalFileSize
	}

	return m.Manifest.RootEntry.FileSize
}

// TotalFiles returns the total number of files in the snapshot.
func (m *ManifestWithError) TotalFiles() int64 {
	if m.Manifest.RootEntry == nil {
		return 0
	}

	if m.Manifest.RootEntry.DirSummary != nil {
		return m.Manifest.RootEntry.DirSummary.TotalFileCount
	}

	return 1
}

// TotalFilesDelta returns the total number of new files in the snapshot.
func (m *ManifestWithError) TotalFilesDelta() int64 {
	if m.Previous == nil {
		return 0
	}

	if m.Manifest.RootEntry == nil || m.Previous.RootEntry == nil {
		return 0
	}

	if m.Manifest.RootEntry.DirSummary != nil && m.Previous.RootEntry.DirSummary != nil {
		return m.Manifest.RootEntry.DirSummary.TotalFileCount - m.Previous.RootEntry.DirSummary.TotalFileCount
	}

	return 1
}

// TotalDirs returns the total number of directories in the snapshot.
func (m *ManifestWithError) TotalDirs() int64 {
	if m.Manifest.RootEntry == nil {
		return 0
	}

	if m.Manifest.RootEntry.DirSummary != nil {
		return m.Manifest.RootEntry.DirSummary.TotalDirCount
	}

	return 0
}

// TotalDirsDelta returns the total number of new directories in the snapshot.
func (m *ManifestWithError) TotalDirsDelta() int64 {
	if m.Previous == nil {
		return 0
	}

	if m.Manifest.RootEntry == nil || m.Previous.RootEntry == nil {
		return 0
	}

	if m.Manifest.RootEntry.DirSummary != nil && m.Previous.RootEntry.DirSummary != nil {
		return m.Manifest.RootEntry.DirSummary.TotalDirCount - m.Previous.RootEntry.DirSummary.TotalDirCount
	}

	return 0
}

// Duration returns the duration of the snapshot.
func (m *ManifestWithError) Duration() time.Duration {
	return time.Duration(m.Manifest.EndTime - m.Manifest.StartTime).Round(durationPrecision)
}

// StatusCode returns the status code of the manifest.
func (m *ManifestWithError) StatusCode() string {
	if m.Error != "" {
		return "fatal"
	}

	if m.Manifest.IncompleteReason != "" {
		return "incomplete"
	}

	if m.Manifest.RootEntry != nil && m.Manifest.RootEntry.DirSummary != nil {
		if m.Manifest.RootEntry.DirSummary.FatalErrorCount > 0 {
			return "fatal"
		}

		if m.Manifest.RootEntry.DirSummary.IgnoredErrorCount > 0 {
			return "error"
		}
	}

	return "ok"
}

// MultiSnapshotStatus represents the status of multiple snapshots.
type MultiSnapshotStatus struct {
	Snapshots []*ManifestWithError `json:"snapshots"`
}
