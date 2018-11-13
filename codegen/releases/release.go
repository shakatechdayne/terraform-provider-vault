package releases

import (
	"fmt"
	"strconv"
	"strings"
)

func Parse(release string) (*Release, error) {
	canonical := strings.Replace(release, "v", "", -1)
	parts := strings.Split(canonical, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("received unexpected release name: %s", release)
	}
	r := &Release{}
	for i, part := range parts {
		if i == 2 {
			// Pull the product information off the string before
			// parsing it below.
			patchParts := strings.Split(part, "-")
			if len(patchParts) > 1 {
				// The first number is the patch version.
				part = patchParts[0]
				if strings.Contains(patchParts[1], "beta") {
					r.ProductName = "beta"

				} else if strings.Contains(patchParts[1], "rc") {
					r.ProductName = "rc"
				} else {
					return nil, fmt.Errorf("unable to parse %s", release)
				}
				productNum, err := strconv.Atoi(patchParts[1][len(r.ProductName):])
				if err != nil {
					return nil, fmt.Errorf("unable to parse %s due to %s", release, err)
				}
				r.ProductNum = productNum
			}
		}
		num, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("unable to parse %s due to %s", release, err)
		}
		switch i {
		case 0:
			r.Major = num
		case 1:
			r.Minor = num
		case 2:
			r.Patch = num
		}
	}
	return r, nil
}

type Release struct {
	// In "v0.8.2-beta1", Major would be 0.
	Major int

	// In "v0.8.2-beta1", Minor would be 8.
	Minor int

	// In "v0.8.2-beta1", Patch would be 2.
	Patch int

	// In "v0.8.2-beta1", ProductName would be beta.
	ProductName string

	// In "v0.8.2-beta1", ProductNum would be 1.
	ProductNum int
}

func (r *Release) IsAfter(other *Release) bool {
	if r.Major > other.Major {
		return true
	}
	if r.Major < other.Major {
		return false
	}

	if r.Minor > other.Minor {
		return true
	}
	if r.Minor < other.Minor {
		return false
	}

	if r.Patch > other.Patch {
		return true
	}
	if r.Patch < other.Patch {
		return false
	}

	if r.ProductName == "" {
		if other.ProductName == "rc" || other.ProductName == "beta" {
			return true
		}
	}
	if r.ProductName == "rc" {
		if other.ProductName == "" {
			return false
		}
		if other.ProductName == "beta" {
			return true
		}
	}
	if r.ProductName == "beta" {
		if other.ProductName == "" || other.ProductName == "rc" {
			return false
		}
	}

	if r.ProductNum > other.ProductNum {
		return true
	}
	return false
}
