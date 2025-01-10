package istanbul

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiredMessageCountAndF(t *testing.T) {
	// Expected results for each committeeSize. key is committeeSize.
	expectedResults := map[uint64]struct {
		required int // ceil(N * 2 / 3)
		f        int // ceil(N / 3) - 1
	}{
		4:  {3, 1},
		5:  {4, 1},
		6:  {4, 1}, // notice that f is not exactly same as N - required
		7:  {5, 2},
		8:  {6, 2},
		9:  {6, 2},
		10: {7, 3},
		11: {8, 3},
		12: {8, 3},
		13: {9, 4},
		14: {10, 4},
		15: {10, 4},
		16: {11, 5},
		17: {12, 5},
		18: {12, 5},
		19: {13, 6},
		20: {14, 6},
		21: {14, 6},
		22: {15, 7},
		23: {16, 7},
		24: {16, 7},
		25: {17, 8},
		26: {18, 8},
		27: {18, 8},
		28: {19, 9},
		29: {20, 9},
		30: {20, 9},
		31: {21, 10},
		32: {22, 10},
		33: {22, 10},
		34: {23, 11},
		35: {24, 11},
		36: {24, 11},
		37: {25, 12},
		38: {26, 12},
		39: {26, 12},
		40: {27, 13},
		50: {34, 16},
		60: {40, 19},
	}

	for committeeSize := uint64(4); committeeSize <= 60; committeeSize++ {
		// Get expected results
		expected, exists := expectedResults[committeeSize]
		if !exists {
			// Skip if no expected result is defined
			t.Logf("Skipping committeeSize %d as no expected result is defined.", committeeSize)
			continue
		}
		t.Run("committeeSize "+strconv.FormatUint(committeeSize, 10), func(t *testing.T) {
			// Calculate actual results and check with expected result
			assert.Equal(t, expected.required, requiredMessageCount(int(committeeSize), committeeSize))
			assert.Equal(t, expected.f, f(int(committeeSize), committeeSize))
		},
		)
	}
}
