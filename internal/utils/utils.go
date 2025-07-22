// Package utils provides common utility functions for MoBot 2025
package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// File utilities

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// DirExists checks if a directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// EnsureDir creates directory if it doesn't exist
func EnsureDir(path string) error {
	if !DirExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer source.Close()
	
	// Create destination directory if needed
	if err := EnsureDir(filepath.Dir(dst)); err != nil {
		return fmt.Errorf("failed to create destination dir: %w", err)
	}
	
	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer destination.Close()
	
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("failed to copy: %w", err)
	}
	
	return destination.Sync()
}

// GetFileSize returns file size in bytes
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// String utilities

// TruncateString truncates string to max length with suffix
func TruncateString(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	
	if maxLen <= len(suffix) {
		return suffix[:maxLen]
	}
	
	return s[:maxLen-len(suffix)] + suffix
}

// NormalizeString normalizes string for comparison
func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// SanitizeFilename removes invalid characters from filename
func SanitizeFilename(filename string) string {
	// Replace invalid characters
	invalid := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	result := filename
	
	for _, char := range invalid {
		result = strings.ReplaceAll(result, char, "_")
	}
	
	// Remove leading/trailing spaces and dots
	result = strings.Trim(result, " .")
	
	// Ensure not empty
	if result == "" {
		result = "unnamed"
	}
	
	return result
}

// ID generation

// GenerateID generates a random hex ID
func GenerateID(prefix string, length int) string {
	bytes := make([]byte, length/2)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp
		return fmt.Sprintf("%s%d", prefix, time.Now().UnixNano())
	}
	
	if prefix != "" {
		return fmt.Sprintf("%s%s", prefix, hex.EncodeToString(bytes))
	}
	return hex.EncodeToString(bytes)
}

// GenerateShortID generates a short random ID
func GenerateShortID() string {
	return GenerateID("", 8)
}

// Time utilities

// FormatDuration formats duration in human-readable form
func FormatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	
	if d < time.Hour {
		mins := int(d.Minutes())
		secs := int(d.Seconds()) % 60
		return fmt.Sprintf("%dm%ds", mins, secs)
	}
	
	hours := int(d.Hours())
	mins := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh%dm", hours, mins)
}

// TimeAgo returns human-readable time difference
func TimeAgo(t time.Time) string {
	diff := time.Since(t)
	
	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(diff.Hours()/24))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%d weeks ago", int(diff.Hours()/(24*7)))
	default:
		return t.Format("2006-01-02")
	}
}

// Slice utilities

// Contains checks if slice contains value
func Contains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// Unique returns unique values from slice
func Unique[T comparable](slice []T) []T {
	seen := make(map[T]bool)
	result := make([]T, 0, len(slice))
	
	for _, v := range slice {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	
	return result
}

// Filter returns slice with only values that match predicate
func Filter[T any](slice []T, predicate func(T) bool) []T {
	result := make([]T, 0, len(slice))
	
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	
	return result
}

// Map transforms slice values
func Map[T any, R any](slice []T, transform func(T) R) []R {
	result := make([]R, len(slice))
	
	for i, v := range slice {
		result[i] = transform(v)
	}
	
	return result
}

// Math utilities

// Min returns minimum of two values
func Min[T int | int32 | int64 | float32 | float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns maximum of two values
func Max[T int | int32 | int64 | float32 | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Clamp constrains value between min and max
func Clamp[T int | int32 | int64 | float32 | float64](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Percentage calculates percentage
func Percentage(value, total float64) float64 {
	if total == 0 {
		return 0
	}
	return (value / total) * 100
}

// Error utilities

// Must panics if error is not nil (use only in initialization)
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}

// Retry retries function with exponential backoff
func Retry(attempts int, delay time.Duration, fn func() error) error {
	var err error
	
	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}
		
		if i < attempts-1 {
			time.Sleep(delay)
			delay *= 2 // Exponential backoff
		}
	}
	
	return fmt.Errorf("failed after %d attempts: %w", attempts, err)
}