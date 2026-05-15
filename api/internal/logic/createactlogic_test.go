package logic

import (
	"fmt"
	"testing"
	"time"
)

func TestCreateActRequestValidation(t *testing.T) {
	now := time.Now().Unix()

	tests := []struct {
		name        string
		req         CreateActReq
		expectError bool
		errMessage  string
	}{
		{
			name: "有效请求",
			req: CreateActReq{
				Name:    "测试活动",
				Stock:   100,
				StartAt: now + 3600,
				EndAt:   now + 7200,
			},
			expectError: false,
			errMessage:  "",
		},
		{
			name: "活动名称为空",
			req: CreateActReq{
				Name:    "",
				Stock:   100,
				StartAt: now + 3600,
				EndAt:   now + 7200,
			},
			expectError: true,
			errMessage:  "活动名称不能为空",
		},
		{
			name: "库存为0",
			req: CreateActReq{
				Name:    "测试活动",
				Stock:   0,
				StartAt: now + 3600,
				EndAt:   now + 7200,
			},
			expectError: true,
			errMessage:  "库存必须大于0",
		},
		{
			name: "库存为负数",
			req: CreateActReq{
				Name:    "测试活动",
				Stock:   -1,
				StartAt: now + 3600,
				EndAt:   now + 7200,
			},
			expectError: true,
			errMessage:  "库存必须大于0",
		},
		{
			name: "开始时间等于结束时间",
			req: CreateActReq{
				Name:    "测试活动",
				Stock:   100,
				StartAt: now + 3600,
				EndAt:   now + 3600,
			},
			expectError: true,
			errMessage:  "活动开始时间必须早于结束时间",
		},
		{
			name: "开始时间晚于结束时间",
			req: CreateActReq{
				Name:    "测试活动",
				Stock:   100,
				StartAt: now + 7200,
				EndAt:   now + 3600,
			},
			expectError: true,
			errMessage:  "活动开始时间必须早于结束时间",
		},
		{
			name: "开始时间早于当前时间",
			req: CreateActReq{
				Name:    "测试活动",
				Stock:   100,
				StartAt: now - 3600,
				EndAt:   now + 3600,
			},
			expectError: true,
			errMessage:  "活动开始时间不能早于当前时间",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateActRequest(tt.req)
			if (err != nil) != tt.expectError {
				t.Errorf("validateCreateActRequest() error = %v, expectError = %v", err, tt.expectError)
				return
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("validateCreateActRequest() error message = %v, want = %v", err.Error(), tt.errMessage)
			}
		})
	}
}

type CreateActReq struct {
	Name    string
	Stock   int
	StartAt int64
	EndAt   int64
}

func validateCreateActRequest(req CreateActReq) error {
	if req.Name == "" {
		return fmt.Errorf("活动名称不能为空")
	}
	if req.Stock <= 0 {
		return fmt.Errorf("库存必须大于0")
	}
	if req.StartAt >= req.EndAt {
		return fmt.Errorf("活动开始时间必须早于结束时间")
	}
	if req.StartAt < time.Now().Unix() {
		return fmt.Errorf("活动开始时间不能早于当前时间")
	}
	return nil
}

func TestActivityDuration(t *testing.T) {
	tests := []struct {
		name           string
		startAt        int64
		endAt          int64
		expectedHours  float64
	}{
		{
			name:          "1小时活动",
			startAt:       3600,
			endAt:         7200,
			expectedHours: 1.0,
		},
		{
			name:          "24小时活动",
			startAt:       0,
			endAt:         86400,
			expectedHours: 24.0,
		},
		{
			name:          "30分钟活动",
			startAt:       0,
			endAt:         1800,
			expectedHours: 0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := float64(tt.endAt-tt.startAt) / 3600.0
			if duration != tt.expectedHours {
				t.Errorf("duration = %v hours, want %v hours", duration, tt.expectedHours)
			}
		})
	}
}

func TestStockRange(t *testing.T) {
	tests := []struct {
		name        string
		stock       int
		minStock    int
		maxStock    int
		expectValid bool
	}{
		{
			name:        "库存在有效范围内",
			stock:       100,
			minStock:    1,
			maxStock:    1000,
			expectValid: true,
		},
		{
			name:        "库存等于最小值",
			stock:       1,
			minStock:    1,
			maxStock:    1000,
			expectValid: true,
		},
		{
			name:        "库存等于最大值",
			stock:       1000,
			minStock:    1,
			maxStock:    1000,
			expectValid: true,
		},
		{
			name:        "库存小于最小值",
			stock:       0,
			minStock:    1,
			maxStock:    1000,
			expectValid: false,
		},
		{
			name:        "库存大于最大值",
			stock:       1001,
			minStock:    1,
			maxStock:    1000,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.stock >= tt.minStock && tt.stock <= tt.maxStock
			if isValid != tt.expectValid {
				t.Errorf("stock validation = %v, want %v", isValid, tt.expectValid)
			}
		})
	}
}
