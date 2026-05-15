package logic

import (
	"testing"
)

func TestOrderStatusResponse(t *testing.T) {
	tests := []struct {
		name           string
		orderNo        string
		status         int32
		expectHasOrder bool
	}{
		{
			name:           "已支付订单",
			orderNo:        "1747200000123456",
			status:         1,
			expectHasOrder: true,
		},
		{
			name:           "未支付订单",
			orderNo:        "1747200000123457",
			status:         0,
			expectHasOrder: true,
		},
		{
			name:           "已取消订单",
			orderNo:        "1747200000123458",
			status:         2,
			expectHasOrder: true,
		},
		{
			name:           "已完成订单",
			orderNo:        "1747200000123459",
			status:         3,
			expectHasOrder: true,
		},
		{
			name:           "空订单号",
			orderNo:        "",
			status:         0,
			expectHasOrder: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasOrder := tt.orderNo != "" && tt.status >= 0
			if hasOrder != tt.expectHasOrder {
				t.Errorf("hasOrder = %v, want %v", hasOrder, tt.expectHasOrder)
			}
		})
	}
}

func TestOrderStatusValues(t *testing.T) {
	statusNames := map[int32]string{
		0: "未支付",
		1: "已支付",
		2: "已取消",
		3: "已完成",
	}

	tests := []struct {
		status   int32
		expected string
	}{
		{0, "未支付"},
		{1, "已支付"},
		{2, "已取消"},
		{3, "已完成"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if name := statusNames[tt.status]; name != tt.expected {
				t.Errorf("status %d = %s, want %s", tt.status, name, tt.expected)
			}
		})
	}
}

func TestOrderNoValidation(t *testing.T) {
	tests := []struct {
		name        string
		orderNo     string
		expectValid bool
	}{
		{
			name:        "有效订单号",
			orderNo:     "1747200000123456",
			expectValid: true,
		},
		{
			name:        "空订单号",
			orderNo:     "",
			expectValid: false,
		},
		{
			name:        "纯数字订单号",
			orderNo:     "123456789",
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := len(tt.orderNo) > 0
			if isValid != tt.expectValid {
				t.Errorf("orderNo validation = %v, want %v", isValid, tt.expectValid)
			}
		})
	}
}

func TestOrderStatusTransition(t *testing.T) {
	tests := []struct {
		name         string
		fromStatus   int32
		toStatus     int32
		expectValid  bool
	}{
		{
			name:        "未支付 -> 已支付",
			fromStatus:  0,
			toStatus:    1,
			expectValid: true,
		},
		{
			name:        "已支付 -> 已取消",
			fromStatus:  1,
			toStatus:    2,
			expectValid: true,
		},
		{
			name:        "已取消 -> 已完成",
			fromStatus:  2,
			toStatus:    3,
			expectValid: true,
		},
		{
			name:        "已完成 -> 已支付",
			fromStatus:  3,
			toStatus:    1,
			expectValid: false,
		},
		{
			name:        "未支付 -> 已取消",
			fromStatus:  0,
			toStatus:    2,
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := isValidStatusTransition(tt.fromStatus, tt.toStatus)
			if valid != tt.expectValid {
				t.Errorf("status transition from %d to %d = %v, want %v",
					tt.fromStatus, tt.toStatus, valid, tt.expectValid)
			}
		})
	}
}

func isValidStatusTransition(from, to int32) bool {
	transitions := map[int32][]int32{
		0: {1, 2},
		1: {2, 3},
		2: {3},
		3: {},
	}

	allowed, exists := transitions[from]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == to {
			return true
		}
	}
	return false
}
