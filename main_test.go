package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidVideoFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{
			name:     "mp4 file",
			filename: "test.mp4",
			expected: true,
		},
		{
			name:     "avi file",
			filename: "test.avi",
			expected: true,
		},
		{
			name:     "mov file",
			filename: "test.mov",
			expected: true,
		},
		{
			name:     "mkv file",
			filename: "test.mkv",
			expected: true,
		},
		{
			name:     "wmv file",
			filename: "test.wmv",
			expected: true,
		},
		{
			name:     "flv file",
			filename: "test.flv",
			expected: true,
		},
		{
			name:     "webm file",
			filename: "test.webm",
			expected: true,
		},
		{
			name:     "uppercase mp4",
			filename: "test.MP4",
			expected: true,
		},
		{
			name:     "mixed case avi",
			filename: "test.AVI",
			expected: true,
		},
		{
			name:     "txt file",
			filename: "test.txt",
			expected: false,
		},
		{
			name:     "jpg file",
			filename: "test.jpg",
			expected: false,
		},
		{
			name:     "png file",
			filename: "test.png",
			expected: false,
		},
		{
			name:     "pdf file",
			filename: "test.pdf",
			expected: false,
		},
		{
			name:     "no extension",
			filename: "test",
			expected: false,
		},
		{
			name:     "empty filename",
			filename: "",
			expected: false,
		},
		{
			name:     "only extension",
			filename: ".mp4",
			expected: true,
		},
		{
			name:     "path with video file",
			filename: "path/to/test.mp4",
			expected: true,
		},
		{
			name:     "windows path with video file",
			filename: "C:\\path\\to\\test.mp4",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidVideoFile(tt.filename)
			assert.Equal(t, tt.expected, result)
		})
	}
}
