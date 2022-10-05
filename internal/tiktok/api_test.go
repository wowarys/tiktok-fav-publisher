package tiktok

import (
	"fmt"
	"testing"

	"go.uber.org/zap"
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
			logger, err := zap.NewProduction()
			assert(t, err == nil)
			tt := NewServiceApi(test.name, test.count, logger)
			videos, err := tt.GetLikedVideos()
			assert(t, err == nil, err)
			assert(t, test.count == len(videos), fmt.Sprintf("expected: %d, got: %d", test.count, len(videos)))
		})
	}
}
