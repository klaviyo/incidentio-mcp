package tools

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	// Maximum lengths for various fields
	MaxNameLength        = 200
	MaxSummaryLength     = 1000
	MaxDescriptionLength = 5000
	MaxMessageLength     = 5000
	MaxIDLength          = 100
)

// ValidateStringInput validates a string input with common security checks
func ValidateStringInput(value string, fieldName string, maxLength int, required bool) error {
	if required && value == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	if len(value) > maxLength {
		return fmt.Errorf("%s exceeds maximum length of %d characters", fieldName, maxLength)
	}

	// Check for potential injection attempts (very basic)
	if containsSuspiciousPatterns(value) {
		return fmt.Errorf("%s contains invalid characters", fieldName)
	}

	return nil
}

// ValidateID validates an ID format
func ValidateID(id string, fieldName string) error {
	if id == "" {
		return fmt.Errorf("%s is required", fieldName)
	}

	if len(id) > MaxIDLength {
		return fmt.Errorf("%s exceeds maximum length of %d characters", fieldName, MaxIDLength)
	}

	// IDs should typically be alphanumeric with underscores/hyphens
	if !isValidIDFormat(id) {
		return fmt.Errorf("%s has invalid format", fieldName)
	}

	return nil
}

// containsSuspiciousPatterns checks for basic injection patterns
func containsSuspiciousPatterns(input string) bool {
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"onclick=",
		"onerror=",
		"\x00", // null bytes
		"../",  // path traversal
		"..\\", // path traversal windows
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range suspiciousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return true
		}
	}

	return false
}

// isValidIDFormat checks if an ID has a valid format
func isValidIDFormat(id string) bool {
	// Allow alphanumeric, underscores, hyphens, and periods
	validID := regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
	return validID.MatchString(id)
}

// SanitizeErrorMessage removes sensitive information from error messages
func SanitizeErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	errMsg := err.Error()

	// Remove potential API keys or tokens
	apiKeyPattern := regexp.MustCompile(`(?i)(api[_-]?key|token|bearer|authorization)[:\s=]+[^\s]+`)
	errMsg = apiKeyPattern.ReplaceAllString(errMsg, "$1: [REDACTED]")

	// Remove potential URLs with credentials
	urlPattern := regexp.MustCompile(`https?://[^:]+:[^@]+@[^\s]+`)
	errMsg = urlPattern.ReplaceAllString(errMsg, "https://[REDACTED]@...")

	return errMsg
}