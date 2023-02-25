package data

import (
	"net/url"

	"gorm.io/gorm"
)

type BookmarkPrivacy string

const (
	BookmarkPrivacyPublic   BookmarkPrivacy = "public"
	BookmarkPrivacyPrivate  BookmarkPrivacy = "private"
	BookmarkPrivacyQueryAll BookmarkPrivacy = "all"
)

type Bookmark struct {
	gorm.Model
	URL         string `gorm:"not null;default:null"`
	Title       string
	Description string
	Privacy     BookmarkPrivacy `gorm:"default:'private'"`

	Tags []Tag `gorm:"many2many:bookmark_tags;"`
}

type BookmarkForm struct {
	URL         string
	Title       string
	Description string
	Public      bool
	Tags        string
}

type BookmarkListRequest struct {
	Privacy BookmarkPrivacy
	TagID   uint
	Offset  int
	Order   string

	PaginationPathPrefix string
}

type BookmarkListResult struct {
	Items   []Bookmark
	Count   int64
	HasPrev bool
	PrevURL string
	HasNext bool
	NextURL string
}

type BookmarkSearchRequest struct {
	Query  string
	Offset int

	PaginationPathPrefix string
}

type BookmarkSearchResponse struct {
	Items   []Bookmark
	Count   int64
	HasPrev bool
	PrevURL string
	HasNext bool
	NextURL string
}

func (b *Bookmark) IsPublic() bool {
	return b.Privacy == BookmarkPrivacyPublic
}

func (req *BookmarkForm) IsValid() *ValidationError {
	isErr := false
	fields := make(map[string]string)

	if req.URL == "" {
		isErr = true
		fields["URL"] = "URL is required"
	} else {
		parsedURL, err := url.Parse(req.URL)
		urlParseError := "URL format is invalid"
		if err != nil {
			isErr = true
			fields["URL"] = urlParseError
		}
		if parsedURL.Scheme == "" || parsedURL.Host == "" {
			isErr = true
			fields["URL"] = urlParseError
		}
	}

	if isErr {
		return NewValidationError("Bookmark is invalid", fields)
	}

	return nil
}
