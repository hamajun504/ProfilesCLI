package profile

import (
	"reflect"
	"slices"
	"testing"
)

func assertProfileNames(t *testing.T, profiles []Profile, want []string, ignoreOrder bool) {
	t.Helper()

	gotNames := make([]string, 0, len(profiles))
	for _, p := range profiles {
		gotNames = append(gotNames, p.Name)
	}

	got := gotNames
	expected := want
	if ignoreOrder {
		got = slices.Clone(gotNames)
		expected = slices.Clone(want)
		slices.Sort(got)
		slices.Sort(expected)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected names (ignoreOrder=%t): got %v, want %v", ignoreOrder, got, expected)
	}
}
