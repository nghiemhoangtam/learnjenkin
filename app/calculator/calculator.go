// Package calculator cung cấp các phép toán số học cơ bản.
// Đây là phần "business logic" để chúng ta viết unit test và cho Jenkins chạy.
package calculator

import "errors"

// ErrDivideByZero được trả về khi thực hiện phép chia cho 0.
var ErrDivideByZero = errors.New("không thể chia cho 0")

// Add trả về tổng của a và b.
func Add(a, b int) int {
	return a + b
}

// Subtract trả về hiệu của a và b.
func Subtract(a, b int) int {
	return a - b
}

// Multiply trả về tích của a và b.
func Multiply(a, b int) int {
	return a * b
}

// Divide trả về thương của a và b.
// Nếu b == 0 thì trả về lỗi ErrDivideByZero.
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}
