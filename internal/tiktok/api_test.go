package tiktok

import (
	"fmt"
	"testing"
)

func assert(t testing.TB, expected bool, args ...interface{}) {
	if !expected {
		t.Helper()
		t.Fatal(args...)
	}
}

func TestGetLikedVideos(t *testing.T) {
	var tests = []struct {
		name     string
		username string
		count    int
	}{
		{
			"TestCount10",
			"tiktoklikebot341",
			10,
		},
		{
			"TestCount35",
			"tiktoklikebot341",
			25,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := NewServiceApi(test.name, test.count)
			videos, err := tt.GetLikedVideos()
			assert(t, err == nil, err)
			assert(t, test.count == len(videos), fmt.Sprintf("expected: %d, got: %d", test.count, len(videos)))
		})
	}
}
