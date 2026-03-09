package utils

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
)

const dateLayout = "02.01.2006"

// dd.mm.yyyy
func InputDate(r io.Reader) (time.Time, error) {
	fmt.Println("Введите дату в формате ДД.ММ.ГГГГ")

	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, fmt.Errorf("ошибка чтения ввода: %w", err)
	}

	input = strings.TrimSpace(input)

	if input == "" {
		return time.Now(), nil
	}

	date, err := time.Parse(dateLayout, input)
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный формат даты %q, ожидается %s: %w",
			input, dateLayout, err)
	}

	return date, nil
}
