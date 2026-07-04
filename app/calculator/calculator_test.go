package calculator

import (
	"errors"
	"testing"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		name     string
		a, b     int
		expected int
	}{
		{"hai số dương", 2, 3, 5},
		{"có số âm", -2, 3, 1},
		{"cộng với 0", 10, 0, 10},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := Add(c.a, c.b); got != c.expected {
				t.Errorf("Add(%d, %d) = %d; muốn %d", c.a, c.b, got, c.expected)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	if got := Subtract(10, 4); got != 6 {
		t.Errorf("Subtract(10, 4) = %d; muốn 6", got)
	}
}

func TestMultiply(t *testing.T) {
	if got := Multiply(6, 7); got != 42 {
		t.Errorf("Multiply(6, 7) = %d; muốn 42", got)
	}
}

func TestDivide(t *testing.T) {
	t.Run("chia bình thường", func(t *testing.T) {
		got, err := Divide(10, 2)
		if err != nil {
			t.Fatalf("không mong đợi lỗi: %v", err)
		}
		if got != 5 {
			t.Errorf("Divide(10, 2) = %d; muốn 5", got)
		}
	})

	t.Run("chia cho 0 phải trả về lỗi", func(t *testing.T) {
		_, err := Divide(10, 0)
		if !errors.Is(err, ErrDivideByZero) {
			t.Errorf("Divide(10, 0) trả về lỗi %v; muốn %v", err, ErrDivideByZero)
		}
	})
}
