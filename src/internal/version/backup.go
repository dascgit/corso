package version

import "math"

const Backup = 3

// Various labels to refer to important version changes.
// Labels don't need 1:1 service:version representation.  Add a new
// label when it's important to mark a delta in behavior that's handled
// somewhere in the logic.
// Labels should state their application, the backup version number,
// and the colloquial purpose of the label.
const (
	// OneDrive1DataAndMetaFiles is the corso backup format version
	// in which we split from storing just the data to storing both
	// the data and metadata in two files.
	OneDrive1DataAndMetaFiles = 1

	// OneDrive3IsMetaMarker is a small improvement on
	// VersionWithDataAndMetaFiles, but has a marker IsMeta which
	// specifies if the file is a meta file or a data file.
	OneDrive3IsMetaMarker = 3

	// OneDrive4IncludesPermissions includes permissions in the backup.
	// Note that this is larger than the current backup version.  That's
	// because it isn't implemented yet.  But we have tests based on this,
	// so maybe we just keep bumping the verson ahead of the backup until
	// it gets implemented.
	OneDriveXIncludesPermissions = Backup + 1

	// OneDriveXNameInMeta points to the backup format version where we begin
	// storing files in kopia with their item ID instead of their OneDrive file
	// name.
	// TODO(ashmrtn): Update this to a real value when we merge the file name
	// change. Set to MAXINT for now to keep the if-check using it working.
	OneDriveXNameInMeta = math.MaxInt
)