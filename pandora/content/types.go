package content

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type (
	Identifiable interface {
		// GetId returns the identifier of the content card
		GetId() uuid.UUID

		// GetSlug should return the unique url identifier for the content
		GetSlug() string

		// GetContentType returns the content types of the card
		GetContentType() Type
	}

	Generic interface {
		Identifiable
	}
)

func HasType(types []Type, contentType Type) bool {
	for _, t := range types {
		if t == contentType {
			return true
		}
	}

	return false
}

type Type string

const (
	TypeEbook        Type = "ebook"
	TypePodcast      Type = "podcast"
	TypeEvent        Type = "event"
	TypeExpert       Type = "expert"
	TypeStream       Type = "stream"
	TypeLearningPath Type = "learningPath"
)

func (e Type) IsValid() bool {
	switch e {
	case TypeEbook, TypePodcast, TypeEvent,
		TypeExpert, TypeStream, TypeLearningPath:
		return true
	}

	return false
}

func (e Type) String() string {
	return string(e)
}

func (e Type) GetAsCopy() string {
	switch e {
	case TypeEbook:
		return "Power read"

	case TypePodcast:
		return "Podcast"

	case TypeEvent:
		return "Experience"

	case TypeStream:
		return "Livestream"

	case TypeExpert:
		return "Thinkfluencer"

	case TypeLearningPath:
		return "Trail"

	default:
		return string(e)
	}
}

func (e *Type) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = Type(strings.ToLower(str))

	if *e == "learning_path" {
		*e = TypeLearningPath
	}

	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid Type", str)
	}

	return nil
}

func (e Type) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(strings.ToUpper(e.String())))
}

func (e Type) GetAsUrlSlug() string {
	switch e {
	case TypeEbook:
		return "power-read"

	case TypePodcast:
		return "podcast"

	case TypeEvent:
		return "experience"

	case TypeStream:
		return "stream"

	case TypeExpert:
		return "thinkfluencer"

	case TypeLearningPath:
		return "trail"

	default:
		return string(e)
	}
}
