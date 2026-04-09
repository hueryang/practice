package input

import (
	"fmt"
	"os"
)

// ReadMeetingFile reads the entire meeting minutes file and returns raw bytes.
// Content is printed by the caller (preserves exact bytes, including UTF-8).
func ReadMeetingFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取会议纪要文件 %q: %w", path, err)
	}
	return data, nil
}
