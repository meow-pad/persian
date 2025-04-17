package timeutil

import (
	"fmt"
	"testing"
	"time"
)

func Test_ToWeekZeroTime(t *testing.T) {
	// 示例时间
	t0 := time.Date(2023, 9, 30, 12, 0, 0, 0, time.UTC) // 2023-10-01 是周六
	t1 := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC) // 2023-10-01 是周日
	t2 := time.Date(2023, 10, 2, 12, 0, 0, 0, time.UTC) // 2023-10-02 是周一

	// 指定一周的第一天为周一
	firstDayOfWeek := time.Monday

	// 判断是否在同一周内
	if IsSameWeek(t0, t1, firstDayOfWeek) {
		fmt.Println("t0 and t1 in same week")
	} else {
		fmt.Println("t0 and t1 not in same week")
	}
	if IsSameWeek(t1, t2, firstDayOfWeek) {
		fmt.Println("t1 and t2 in same week")
	} else {
		fmt.Println("t1 and t2 not in same week")
	}
}
