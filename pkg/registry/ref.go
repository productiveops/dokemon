package registry

import (
	"fmt"
	"strings"

	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/types"
	"github.com/distribution/reference"
)

func ImageReference(name string) (types.ImageReference, error) {
	ref, err := getNamedReference(name)
	if err != nil {
		return nil, err
	}

	refstring := ref.String()
	if !strings.HasPrefix(refstring, "//") {
		refstring = fmt.Sprintf("//%s", refstring)
	}

	return docker.ParseReference(refstring)
}

func getNamedReference(name string) (reference.Named, error) {
	name = strings.TrimPrefix(name, "//")

	ref, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return nil, err
	}

	if _, ok := ref.(reference.Named); !ok {
		return nil, err
	}

	if _, hasTag := ref.(reference.NamedTagged); hasTag {
		ref, err = normalizeTaggedDigestedNamed(ref)
		if err != nil {
			return nil, err
		}
	} else if _, hasDigest := ref.(reference.Digested); hasDigest {
		ref = reference.TrimNamed(ref)
	}

	return reference.TagNameOnly(ref), nil
}

func normalizeTaggedDigestedNamed(named reference.Named) (reference.Named, error) {
	_, isDigested := named.(reference.Digested)
	if !isDigested {
		return named, nil
	}

	tag, isTagged := named.(reference.NamedTagged)
	if !isTagged {
		return named, nil
	}

	newNamed := reference.TrimNamed(named)
	newNamed, err := reference.WithTag(newNamed, tag.Tag())
	if err != nil {
		return named, err
	}
	
	return newNamed, nil
}
