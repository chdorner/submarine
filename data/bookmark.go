package data

import (
	"net/url"

	"gorm.io/gorm"
)

type Bookmark struct {
	gorm.Model
	URL        string `gorm:"not null;default:null"`
	Title      string
	Descripton string
}

type BookmarkCreate struct {
	URL         string
	Title       string
	Description string
}

func (req *BookmarkCreate) IsValid() *ValidationError {
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
