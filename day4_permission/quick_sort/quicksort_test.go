package dont_ask

import (
	"reflect"
	"testing"
)

func TestQuickSort(t *testing.T) {
	tests := []struct {
		name  string
		input []int
		want  []int
	}{
		{"empty", []int{}, []int{}},
		{"single element", []int{1}, []int{1}},
		{"already sorted", []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 4, 5}},
		{"reverse sorted", []int{5, 4, 3, 2, 1}, []int{1, 2, 3, 4, 5}},
		{"duplicates", []int{3, 1, 4, 1, 5, 9, 2, 6, 5}, []int{1, 1, 2, 3, 4, 5, 5, 6, 9}},
		{"negative numbers", []int{-3, 0, 2, -1, 5}, []int{-3, -1, 0, 2, 5}},
		{"all same", []int{4, 4, 4}, []int{4, 4, 4}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QuickSort(tt.input)
			if tt.input == nil || len(tt.input) == 0 {
				if len(got) != 0 {
					t.Errorf("QuickSort(%v) = %v, want empty", tt.input, got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuickSort(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
