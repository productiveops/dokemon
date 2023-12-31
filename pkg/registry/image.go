package registry

import (
	"github.com/containers/image/v5/docker/reference"
)

type Image struct {
	Domain  string
	Path    string
	Tag     string

	named reference.Named
}

func ParseImage(name string) (Image, error) {
	named, err := reference.ParseNormalizedNamed(name)
	if err != nil {
		return Image{}, err
	}
	named = reference.TagNameOnly(named)

	i := Image{
		named:  named,
		Domain: reference.Domain(named),
		Path:   reference.Path(named),
	}

	if tagged, ok := named.(reference.Tagged); ok {
		i.Tag = tagged.Tag()
	}

	return i, nil
}

func (i Image) String() string {
	return i.named.String()
}
