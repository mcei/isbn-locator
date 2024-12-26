package internal

import (
	"errors"
	"strconv"
	"strings"
)

var (
	errEmptyString      = errors.New("ISBN cannot be an empty string")
	errUnexpectedFormat = errors.New("unexpected format of ISBN")
	errInvalidNumber    = errors.New("invalid ISBN")
)

// CheckISBN проверяет корректность формата строки с международным книжным номером
// https://en.wikipedia.org/wiki/ISBN
// Возвращает errEmptyString если пустая строка
// Возвращает errUnexpectedFormat если строка не соответствует стандарту ISBN-10 или ISBN-13
func CheckISBN(isbn string) error {
	if isbn == "" {
		return errEmptyString
	}
	nums := strings.Replace(isbn, "-", "", -1)

	if len(nums) == 10 {
		return sum10(nums)
	} else if len(nums) == 13 {
		return sum13(nums)
	} else {
		return errUnexpectedFormat
	}
}

// sum10 проверяет сумму цифр в строке формата ISBN-10 без дефисов
// Возвращает errUnexpectedFormat если формат строки не соответствует стандарту
// Возвращает errInvalidNumber в случае неуспешной проверки контрольной суммы цифр
func sum10(s string) error {
	sum := 0
	length := len(s)
	for i, r := range s {
		n, err := strconv.Atoi(string(r))
		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				if i == 9 && r == 88 {
					// обрабатываем случай, когда
					// контрольная цифра может быть римской X
					sum += 10
					continue
				}
				return errUnexpectedFormat
			}
			return err
		}
		k := length - i
		sum += n * k
	}
	if sum%11 != 0 {
		return errInvalidNumber
	}
	return nil
}

// sum13 проверяет сумму цифр в строке формата ISBN-13 без дефисов
// Возвращает errUnexpectedFormat если формат строки не соответствует стандарту
// Возвращает errInvalidNumber в случае неуспешной проверки контрольной суммы цифр
func sum13(s string) error {
	sum, k, j := 0, 1, 3
	for _, r := range s {
		n, err := strconv.Atoi(string(r))
		if err != nil {
			if errors.Is(err, strconv.ErrSyntax) {
				return errUnexpectedFormat
			}
			return err
		}
		sum += n * k
		k, j = j, k
	}
	if sum%10 != 0 {
		return errInvalidNumber
	}
	return nil
}
