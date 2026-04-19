package window

import (
	"context"
	"fmt"
	"strings"
)

// FindByTitle returns the first window whose title contains substr
// (case-insensitive). Error messages are standardized for callers.
func FindByTitle(ctx context.Context, m Manager, substr string) (Info, error) {
	wins, err := m.List(ctx)
	if err != nil {
		return Info{}, err
	}
	lc := strings.ToLower(substr)
	for _, w := range wins {
		if strings.Contains(strings.ToLower(w.Title), lc) {
			return w, nil
		}
	}
	return Info{}, fmt.Errorf("window matching %q not found", substr)
}
