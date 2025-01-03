package timeutil

import "fmt"
import "time"

const (
	MilliSecondPerSecond = 1000
	SecondPerMinute      = 60
	MinutePerHour        = 60
	SecondPerHour        = SecondPerMinute * MinutePerHour
	HourPerDay           = 24

	DurationDay = 24 * time.Hour
)

// 带当前时区的零值时间
var ZeroTime, _ = time.Parse("2006-01-02 15:04:05", "1970-01-01 00:00:00")

var weekdays = []time.Weekday{
	time.Sunday,
	time.Monday,
	time.Tuesday,
	time.Wednesday,
	time.Thursday,
	time.Friday,
	time.Saturday,
}

func SecToMS(sec int64) int64 {
	return sec * MilliSecondPerSecond
}

func MinToMS(min int64) int64 {
	return min * SecondPerMinute * MilliSecondPerSecond
}

func DayToMS(day int64) int64 {
	return day * HourPerDay * SecondPerHour * MilliSecondPerSecond
}

func MinToSec(min int64) int64 {
	return min * SecondPerMinute
}

func SecToMin(sec int64) int64 {
	return sec / SecondPerMinute
}

func MSToSec(milliSecond int64) int64 {
	return milliSecond / MilliSecondPerSecond
}

func ToWeekday(weekday int) (time.Weekday, error) {
	for _, wDay := range weekdays {
		if weekday == int(wDay) {
			return wDay, nil
		}
	}
	return 0, fmt.Errorf("unknown weekday(%d)", weekday)
}

// ChangeDate
//
//	@Description: 修改时间
func ChangeDate(t time.Time, year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	zeroTime := time.Date(year, month, day, hour, min, sec, nsec, time.Local)
	return zeroTime
}

// ToDayZeroTime
//
//	@Description: 转换到本地时区的零点
//	@param t
//	@return time.Time
func ToDayZeroTime(t time.Time) time.Time {
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	return zeroTime
}

// ToDayZeroTimeWithBiasHour
// @Description: 转换到零点且偏移指定小时
// @param t
// @param biasHour 必须在[0,24)小时之间
// @return time.Time
// @return error
func ToDayZeroTimeWithBiasHour(t time.Time, biasHour int) (time.Time, error) {
	if biasHour < 0 || biasHour >= 24 {
		return t, fmt.Errorf("invalid bias hour:%d", biasHour)
	}
	zeroTime := time.Date(t.Year(), t.Month(), t.Day(), biasHour, 0, 0, 0, time.Local)
	return zeroTime, nil
}

// ToNextDayZeroTimeWithBiasHour
//
//	@Description: 转换到下一天零点且偏移指定小时
//	@param t
//	@param biasHour
//	@return time.Time
//	@return error
func ToNextDayZeroTimeWithBiasHour(t time.Time, biasHour int) (time.Time, error) {
	if biasHour < 0 || biasHour >= 24 {
		return t, fmt.Errorf("invalid bias hour:%d", biasHour)
	}
	zeroTime := time.Date(t.Year(), t.Month(), t.Day()+1, biasHour, 0, 0, 0, time.Local)
	return zeroTime, nil
}

// getFirstWeekDay
//
//	@Description: 获取指定时间的当前周的第一天
//	@param t
//	@param firstWeekday
//	@return int
func getFirstWeekDay(t *time.Time, firstWeekday time.Weekday) int {
	// 获取当前星期几
	weekday := t.Weekday()
	if weekday < firstWeekday {
		return t.Day() - int(weekday) + int(firstWeekday) - 7
	} else {
		return t.Day() - int(weekday) + int(firstWeekday)
	}
}

// ToWeekZeroTime
//
//	@Description: 转换到本周零点
//	@param t
//	@param firstWeekday
//	@return time.Time
func ToWeekZeroTime(t time.Time, firstWeekday time.Weekday) time.Time {
	// 计算一周的第一天
	zeroTime := time.Date(t.Year(), t.Month(), getFirstWeekDay(&t, firstWeekday),
		0, 0, 0, 0, time.Local)
	return zeroTime
}

// ToWeekZeroTimeWithBiasHour
//
//	@Description: 转换到本周零点且偏移指定小时
//	@param t
//	@param firstWeekday
//	@param biasHour
//	@return time.Time
//	@return error
func ToWeekZeroTimeWithBiasHour(t time.Time, firstWeekday time.Weekday, biasHour int) (time.Time, error) {
	if biasHour < 0 || biasHour >= 24 {
		return t, fmt.Errorf("invalid bias hour:%d", biasHour)
	}
	// 计算一周的第一天
	zeroTime := time.Date(t.Year(), t.Month(), getFirstWeekDay(&t, firstWeekday),
		biasHour, 0, 0, 0, time.Local)
	return zeroTime, nil
}

// ToNextWeekZeroTimeWithBiasHour
//
//	@Description: 转换到零点且偏移指定小时
//	@param t
//	@param firstWeekday
//	@param biasHour
//	@return time.Time
//	@return error
func ToNextWeekZeroTimeWithBiasHour(t time.Time, firstWeekday time.Weekday, biasHour int) (time.Time, error) {
	zeroTime, err := ToWeekZeroTimeWithBiasHour(t, firstWeekday, biasHour)
	if err != nil {
		return zeroTime, err
	}

	return zeroTime, nil
}

// ToNextMonthZeroTimeWithBiasHour
//
//	@Description: 获取下个月1号的带偏移零点
//	@param t
//	@param biasHour
//	@return time.Time
//	@return error
func ToNextMonthZeroTimeWithBiasHour(t time.Time, biasHour int) (time.Time, error) {
	if biasHour < 0 || biasHour >= 24 {
		return t, fmt.Errorf("invalid bias hour:%d", biasHour)
	}
	// 获取当前年份和月份
	year, month, _ := t.Date()

	// 计算下一个月份
	nextMonth := month + 1
	if nextMonth > 12 {
		nextMonth = 1
		year += 1
	}
	// 返回下个月的第一天
	zeroTime := time.Date(year, time.Month(nextMonth), 1, biasHour, 0, 0, 0, time.Local)
	return zeroTime, nil
}

// IsSameDay
//
//	@Description: 是否是在同一天
//	@param t1
//	@param t2
//	@return bool
func IsSameDay(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
