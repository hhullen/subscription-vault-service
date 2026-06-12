package supports

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

var (
	validatorInstance validator.Validate = *validator.New()

	dateFormats = []string{
		"2006-01-02",
		"2006/01/02",
		"2006.01.02",
		"02-01-2006",
		"02.01.2006",
		"02/01/2006",
	}
)

func StructValidator() *validator.Validate {
	return &validatorInstance
}

func ParseDate(s string) (time.Time, error) {
	if strings.ToLower(s) == "now" {
		return time.Now(), nil
	}

	var monthsAgo int
	_, err := fmt.Sscanf(s, "%d", &monthsAgo)
	if err == nil && monthsAgo < 0 {
		return time.Now().AddDate(0, monthsAgo, 0), nil
	}

	for _, f := range dateFormats {
		if t, err := time.Parse(f, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("incorrect date format: %s", s)
}

func Concat(ss ...string) string {
	length := 0
	for i := range ss {
		length += len(ss[i])
	}

	var b strings.Builder
	b.Grow(length)

	for i := range ss {
		b.WriteString(ss[i])
	}

	return b.String()
}

func ReadSecretFile(path string) (string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	secret := ""
	_, err = fmt.Fscan(f, &secret)
	if err != nil {
		return "", err
	}

	return string(secret), nil
}

func IsInContainer() bool {
	return os.Getenv("RUNNING_IN_CONTAINER") == "true"
}

func MakeKVMessagesJSON(kvs ...any) (bytes []byte, err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("failed MakeKVMessagesJSON: %v", p)
		}
	}()

	msgs := map[string]any{}
	for i := 0; i < len(kvs)-1; i += 2 {
		key := fmt.Sprint(kvs[i])
		value := kvs[i+1]
		msgs[key] = value
	}

	bytes, err = json.Marshal(msgs)
	return
}
