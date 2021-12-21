package cmd

import (
	"fmt"
	"testing"
)

type parseTest struct {
	filename       string
	transform      bool
	expectedTotal  int
	expectedUnique int
}

func TestParseBase(t *testing.T) {

	tests := []parseTest{
		{filename: "../fixtures/test_01.png", transform: false, expectedTotal: 4, expectedUnique: 3},
		{filename: "../fixtures/test_02.png", transform: false, expectedTotal: 9, expectedUnique: 9},
		{filename: "../fixtures/test_03.png", transform: false, expectedTotal: 4, expectedUnique: 4},
		{filename: "../fixtures/test_01.png", transform: true, expectedTotal: 4, expectedUnique: 3},
		{filename: "../fixtures/test_02.png", transform: true, expectedTotal: 9, expectedUnique: 9},
		{filename: "../fixtures/test_03.png", transform: true, expectedTotal: 4, expectedUnique: 1},
	}

	for _, tc := range tests {

		transformations := []string{}
		verbose := false

		name := fmt.Sprintf("%s-%t", tc.filename, tc.transform)
		t.Run(name, func(t *testing.T) {
			if tc.transform {
				transformations = CreateTransformations()
			}
			tiles, frequencyTiles := Parse(tc.filename, transformations, verbose)
			t.Logf("Parsed %d total tiles, %d unique\n", len(tiles), len(frequencyTiles))
			for i, frequencyTile := range frequencyTiles {
				t.Logf("Index: %d \tHash: %s \tCount: %d \tFirst Location: %v \tTransformation %t\n", i, frequencyTile.Hash, frequencyTile.Count, frequencyTile.FirstLocation, frequencyTile.Transformations)
			}
			actualTotal := len(tiles)
			actualUnique := len(frequencyTiles)
			if actualTotal != tc.expectedTotal {
				t.Errorf("expected %d, got %d", tc.expectedTotal, actualTotal)
			}
			if actualUnique != tc.expectedUnique {
				t.Errorf("expected %d, got %d", tc.expectedUnique, actualUnique)
			}
		})
	}
}
