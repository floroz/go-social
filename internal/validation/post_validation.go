package validation

import (
	"strings"

	"github.com/floroz/go-social/internal/domain"
)

const maxContentLength = 500

func ValidateCreatePostDTO(createPost *domain.CreatePostDTO) error {
	if createPost == nil {
		return domain.NewValidationError("body", "create post data is required")
	}

	if trimmedContent := strings.TrimSpace(createPost.Content); trimmedContent == "" {
		return domain.NewValidationError("content", "content is required")
	} else if len(trimmedContent) > maxContentLength {
		return domain.NewValidationError("content", "content is too long")
	}

	return nil
}
