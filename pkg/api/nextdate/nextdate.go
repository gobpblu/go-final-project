package nextdate

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go-final-project/pkg/utils"
)

var (
	RepeatNotSpecifiedErr = errors.New("не указано правило повторения")
	IncorrectFormatErr    = errors.New("формат не поддерживается или указан неправильно")
	InvalidDaysCountErr   = errors.New("количество дней должно быть от 0 до 400")
	InvalidDateParamErr   = errors.New("некорректный формат даты")
)

const (
	DateLayout = "20060102"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	nowQuery := query.Get("now")
	now := time.Now()
	var err error

	if len(nowQuery) != 0 {
		now, err = time.Parse(DateLayout, nowQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, query.Get("date"), query.Get("repeat"))
	if err != nil {
		utils.WriteBadRequestError(w)
		return
	}
	w.Write([]byte(nextDate))
	w.WriteHeader(http.StatusOK)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if len(repeat) == 0 {
		return "", RepeatNotSpecifiedErr
	}

	startDate, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", InvalidDateParamErr
	}

	repeatArray := strings.Split(repeat, " ")
	arrayLen := len(repeatArray)

	if arrayLen == 2 && repeatArray[0] == "d" {
		return handleDaysShift(now, startDate, repeatArray[1])
	}

	if arrayLen == 1 && repeatArray[0] == "y" {
		return handleYearShift(now, startDate)
	}

	if arrayLen == 2 && repeatArray[0] == "w" {
		return handleWeekShift(now, startDate, repeatArray[1])
	}

	if arrayLen == 2 || arrayLen == 3 && repeatArray[0] == "m" {
		return handleMonthShift(now, startDate, repeatArray)
	}

	return "", IncorrectFormatErr
}

func handleDaysShift(now time.Time, startDate time.Time, daysStr string) (string, error) {
	daysCount, err := strconv.Atoi(daysStr)
	if err != nil {
		return "", IncorrectFormatErr
	}

	if daysCount > 400 || daysCount <= 0 {
		return "", InvalidDaysCountErr
	}

	for {
		startDate = startDate.AddDate(0, 0, daysCount)
		if AfterNow(startDate, now) {
			break
		}
	}

	return formatDate(startDate), nil
}

func handleYearShift(now time.Time, startDate time.Time) (string, error) {

	for {
		startDate = startDate.AddDate(1, 0, 0)
		if AfterNow(startDate, now) {
			break
		}
	}

	return formatDate(startDate), nil
}

func handleWeekShift(now time.Time, startDate time.Time, weekStr string) (string, error) {
	weekdays, err := extractWeekdaysIntsFromStr(strings.Split(weekStr, ","))
	if err != nil {
		return "", IncorrectFormatErr
	}

	if startDate.Before(now) {
		startDate = now
	}

	for {
		startDate = startDate.AddDate(0, 0, 1)
		if isSameWeekday(startDate, weekdays) {
			break
		}
	}

	return formatDate(startDate), nil
}

func handleMonthShift(now time.Time, startDate time.Time, repeatArray []string) (string, error) {
	var monthDays []int
	var err error
	monthDays, err = extractMonthDaysIntsFromStr(strings.Split(repeatArray[1], ","))
	if err != nil {
		return "", err
	}

	if len(monthDays) == 0 {
		return "", IncorrectFormatErr
	}

	var exactMonths []int

	if len(repeatArray) == 3 {
		exactMonths, err = extractMonthIntsFromStr(strings.Split(repeatArray[2], ","))
		if err != nil {
			return "", err
		}
	}

	if startDate.Before(now) {
		startDate = now
	}

	for {
		startDate = startDate.AddDate(0, 0, 1)
		if isValidDay(startDate, monthDays, exactMonths) {
			break
		}
	}

	return formatDate(startDate), nil
}

func AfterNow(date, now time.Time) bool {
	return date.Format(DateLayout) > now.Format(DateLayout)
}

func isSameWeekday(date time.Time, weekdays []int) bool {
	for _, weekday := range weekdays {
		if weekday == 7 && date.Weekday() == 0 {
			return true
		}

		if date.Weekday() == time.Weekday(weekday) {
			return true
		}
	}

	return false
}

func isValidDay(date time.Time, monthDays []int, exactMonths []int) bool {
	checkMonth := date.Month()
	checkDay := date.Day()
	isValidMonth := len(exactMonths) == 0
	isValidDay := false

	for _, month := range exactMonths {
		if month == int(checkMonth) {
			isValidMonth = true
			break
		}
	}

	for _, day := range monthDays {
		if day == checkDay {
			isValidDay = true
			break
		}

		if day == -1 && isLastDay(date) {
			isValidDay = true
			break
		}

		if day == -2 && isPreLastDay(date) {
			isValidDay = true
			break
		}
	}

	return isValidMonth && isValidDay
}

func extractMonthIntsFromStr(strings []string) ([]int, error) {

	numbers := make([]int, len(strings))
	for index, str := range strings {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, IncorrectFormatErr
		}

		if num <= 0 || num > 12 {
			return nil, IncorrectFormatErr
		}

		numbers[index] = num
	}

	return numbers, nil
}

func extractMonthDaysIntsFromStr(strings []string) ([]int, error) {

	numbers := make([]int, len(strings))
	for index, str := range strings {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, IncorrectFormatErr
		}

		if num < -2 || num == 0 || num > 31 {
			return nil, IncorrectFormatErr
		}

		numbers[index] = num
	}

	return numbers, nil
}

func extractWeekdaysIntsFromStr(strings []string) ([]int, error) {

	numbers := make([]int, len(strings))
	for index, str := range strings {
		num, err := strconv.Atoi(str)
		if err != nil {
			return nil, IncorrectFormatErr
		}

		if num < 1 || num > 7 {
			return nil, IncorrectFormatErr
		}

		numbers[index] = num
	}

	return numbers, nil
}

func isLastDay(t time.Time) bool {
	return t.AddDate(0, 0, 1).Day() == 1
}

func isPreLastDay(t time.Time) bool {
	return t.AddDate(0, 0, 2).Day() == 1
}

func formatDate(t time.Time) string {
	return t.Format(DateLayout)
}
