package helper

import (
	"reflect"
	"testing"
)

func TestSliceDiff(t *testing.T) {
	t.Run("test on strings", func(t *testing.T) {
		a := []string{"a", "b"}
		b := []string{"b", "c"}
		want := []string{"c"}
		res := SliceDiff[string](a, b)
		if !reflect.DeepEqual(res, want) {
			t.Errorf("got %v want %v", res, want)
		}
	})
}
