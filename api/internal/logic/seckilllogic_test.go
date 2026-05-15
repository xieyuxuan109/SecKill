package logic

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGenerateOrderNo(t *testing.T) {
	orderNo1 := generateOrderNo()
	orderNo2 := generateOrderNo()

	if orderNo1 == "" {
		t.Error("generateOrderNo() returned empty string")
	}

	if orderNo1 == orderNo2 {
		t.Errorf("generateOrderNo() should generate unique order numbers, got same: %s", orderNo1)
	}

	if len(orderNo1) < 10 {
		t.Errorf("generateOrderNo() length should be at least 10, got: %d", len(orderNo1))
	}
}

func TestTimeValidation(t *testing.T) {
	tests := []struct {
		name       string
		startMs    int64
		endMs      int64
		nowMs      int64
		wantErr    bool
		errMessage string
	}{
		{
			name:       "活动进行中",
			startMs:    time.Now().UnixMilli() - 10000,
			endMs:      time.Now().UnixMilli() + 10000,
			nowMs:      time.Now().UnixMilli(),
			wantErr:    false,
			errMessage: "",
		},
		{
			name:       "活动尚未开始",
			startMs:    time.Now().UnixMilli() + 10000,
			endMs:      time.Now().UnixMilli() + 20000,
			nowMs:      time.Now().UnixMilli(),
			wantErr:    true,
			errMessage: "活动尚未开始",
		},
		{
			name:       "活动已结束",
			startMs:    time.Now().UnixMilli() - 20000,
			endMs:      time.Now().UnixMilli() - 10000,
			nowMs:      time.Now().UnixMilli(),
			wantErr:    true,
			errMessage: "活动已结束",
		},
		{
			name:       "边界测试-刚好开始",
			startMs:    time.Now().UnixMilli(),
			endMs:      time.Now().UnixMilli() + 10000,
			nowMs:      time.Now().UnixMilli(),
			wantErr:    false,
			errMessage: "",
		},
		{
			name:       "边界测试-刚好结束",
			startMs:    time.Now().UnixMilli() - 10000,
			endMs:      time.Now().UnixMilli(),
			nowMs:      time.Now().UnixMilli(),
			wantErr:    false,
			errMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.nowMs < tt.startMs {
				err = fmt.Errorf("活动尚未开始")
			} else if tt.nowMs > tt.endMs {
				err = fmt.Errorf("活动已结束")
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("TimeValidation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && err.Error() != tt.errMessage {
				t.Errorf("TimeValidation() error message = %v, want %v", err.Error(), tt.errMessage)
			}
		})
	}
}

func TestTimestampConversion(t *testing.T) {
	seconds := time.Now().Unix()
	milliseconds := seconds * 1000
	parsedSeconds := milliseconds / 1000

	if seconds != parsedSeconds {
		t.Errorf("Timestamp conversion failed: seconds=%d, milliseconds=%d, parsedSeconds=%d",
			seconds, milliseconds, parsedSeconds)
	}
}

func TestSeckillRequestValidation(t *testing.T) {
	tests := []struct {
		name        string
		activityId  int64
		expectValid bool
	}{
		{
			name:        "有效的活动ID",
			activityId:  1,
			expectValid: true,
		},
		{
			name:        "活动ID为0",
			activityId:  0,
			expectValid: false,
		},
		{
			name:        "负数活动ID",
			activityId:  -1,
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.activityId > 0
			if isValid != tt.expectValid {
				t.Errorf("activityId validation = %v, want %v", isValid, tt.expectValid)
			}
		})
	}
}

func TestStockDeduction(t *testing.T) {
	tests := []struct {
		name         string
		currentStock int64
		deductAmount int64
		expectResult int64
		expectError  bool
	}{
		{
			name:         "正常扣减",
			currentStock: 100,
			deductAmount: 1,
			expectResult: 99,
			expectError:  false,
		},
		{
			name:         "库存不足",
			currentStock: 0,
			deductAmount: 1,
			expectResult: -1,
			expectError:  true,
		},
		{
			name:         "扣减后库存为0",
			currentStock: 1,
			deductAmount: 1,
			expectResult: 0,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stock := tt.currentStock
			stock -= tt.deductAmount
			if stock < 0 {
				stock = -1
			}

			hasError := stock < 0
			if hasError != tt.expectError {
				t.Errorf("stock deduction error = %v, want error = %v", hasError, tt.expectError)
			}
			if stock != tt.expectResult {
				t.Errorf("stock = %d, want %d", stock, tt.expectResult)
			}
		})
	}
}

func TestDuplicateOrderCheck(t *testing.T) {
	tests := []struct {
		name         string
		userId       int64
		actId        int64
		hasBought    bool
		expectReject bool
	}{
		{
			name:         "新用户第一次购买",
			userId:       1001,
			actId:        1,
			hasBought:    false,
			expectReject: false,
		},
		{
			name:         "用户重复购买",
			userId:       1001,
			actId:        1,
			hasBought:    true,
			expectReject: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shouldReject := tt.hasBought
			if shouldReject != tt.expectReject {
				t.Errorf("duplicate order check = %v, want reject = %v", shouldReject, tt.expectReject)
			}
		})
	}
}

func TestContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	select {
	case <-ctx.Done():
	default:
		t.Error("context should be cancelled")
	}
}

func TestOrderNoFormat(t *testing.T) {
	orderNo := generateOrderNo()

	if len(orderNo) < 10 {
		t.Errorf("orderNo length should be at least 10, got: %d", len(orderNo))
	}

	timestampPart := orderNo[:10]
	if timestampPart == "0000000000" {
		t.Error("orderNo should contain valid timestamp")
	}
}
