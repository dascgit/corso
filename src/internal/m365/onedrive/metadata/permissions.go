package metadata

import (
	"context"
	"time"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"golang.org/x/exp/slices"

	"github.com/alcionai/corso/src/internal/common/ptr"
	"github.com/alcionai/corso/src/pkg/logger"
)

type SharingMode int

const (
	SharingModeCustom = SharingMode(iota)
	SharingModeInherited
)

type GV2Type string

const (
	GV2App       GV2Type = "application"
	GV2Device    GV2Type = "device"
	GV2Group     GV2Type = "group"
	GV2SiteUser  GV2Type = "site_user"
	GV2SiteGroup GV2Type = "site_group"
	GV2User      GV2Type = "user"
)

// FilePermission is used to store permissions of a specific resource owner
// to a drive item.
type Permission struct {
	ID         string     `json:"id,omitempty"`
	Roles      []string   `json:"role,omitempty"`
	Email      string     `json:"email,omitempty"`    // DEPRECATED: Replaced with EntityID in newer backups
	EntityID   string     `json:"entityId,omitempty"` // this is the resource owner's ID
	EntityType GV2Type    `json:"entityType,omitempty"`
	Expiration *time.Time `json:"expiration,omitempty"`
}

// isSamePermission checks equality of two UserPermission objects
func (p Permission) Equals(other Permission) bool {
	// EntityID can be empty for older backups and Email can be empty
	// for newer ones. It is not possible for both to be empty.  Also,
	// if EntityID/Email for one is not empty then the other will also
	// have EntityID/Email as we backup permissions for all the
	// parents and children when we have a change in permissions.
	if p.EntityID != "" && p.EntityID != other.EntityID {
		return false
	}

	if p.Email != "" && p.Email != other.Email {
		return false
	}

	p1r := p.Roles
	p2r := other.Roles

	slices.Sort(p1r)
	slices.Sort(p2r)

	return slices.Equal(p1r, p2r)
}

// DiffPermissions compares the before and after set, returning
// the permissions that were added and removed (in that order)
// in the after set.
func DiffPermissions(before, after []Permission) ([]Permission, []Permission) {
	var (
		added   = []Permission{}
		removed = []Permission{}
	)

	for _, cp := range after {
		found := false

		for _, pp := range before {
			if cp.Equals(pp) {
				found = true
				break
			}
		}

		if !found {
			added = append(added, cp)
		}
	}

	for _, pp := range before {
		found := false

		for _, cp := range after {
			if cp.Equals(pp) {
				found = true
				break
			}
		}

		if !found {
			removed = append(removed, pp)
		}
	}

	return added, removed
}

func FilterPermissions(ctx context.Context, perms []models.Permissionable) []Permission {
	up := []Permission{}

	for _, p := range perms {
		if p.GetGrantedToV2() == nil {
			// For link shares, we get permissions without a user
			// specified
			continue
		}

		var (
			// Below are the mapping from roles to "Advanced" permissions
			// screen entries:
			//
			// owner - Full Control
			// write - Design | Edit | Contribute (no difference in /permissions api)
			// read  - Read
			// empty - Restricted View
			//
			// helpful docs:
			// https://devblogs.microsoft.com/microsoft365dev/controlling-app-access-on-specific-sharepoint-site-collections/
			roles    = p.GetRoles()
			gv2      = p.GetGrantedToV2()
			entityID string
			gv2t     GV2Type
		)

		switch true {
		case gv2.GetUser() != nil:
			gv2t = GV2User
			entityID = ptr.Val(gv2.GetUser().GetId())
		case gv2.GetSiteUser() != nil:
			gv2t = GV2SiteUser
			entityID = ptr.Val(gv2.GetSiteUser().GetId())
		case gv2.GetGroup() != nil:
			gv2t = GV2Group
			entityID = ptr.Val(gv2.GetGroup().GetId())
		case gv2.GetSiteGroup() != nil:
			gv2t = GV2SiteGroup
			entityID = ptr.Val(gv2.GetSiteGroup().GetId())
		case gv2.GetApplication() != nil:
			gv2t = GV2App
			entityID = ptr.Val(gv2.GetApplication().GetId())
		case gv2.GetDevice() != nil:
			gv2t = GV2Device
			entityID = ptr.Val(gv2.GetDevice().GetId())
		default:
			logger.Ctx(ctx).Info("untracked permission")
		}

		// Technically GrantedToV2 can also contain devices, but the
		// documentation does not mention about devices in permissions
		if entityID == "" {
			// This should ideally not be hit
			continue
		}

		up = append(up, Permission{
			ID:         ptr.Val(p.GetId()),
			Roles:      roles,
			EntityID:   entityID,
			EntityType: gv2t,
			Expiration: p.GetExpirationDateTime(),
		})
	}

	return up
}