package supports

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	// "github.com/sergi/go-diff/diffmatchpatch"
	// "golang.org/x/crypto/argon2"
)

const (
// defaultHashLen = 97
)

// type HashSettings struct {
// 	memory     uint32
// 	iterations uint32
// 	threads    uint8
// 	saltLength uint32
// 	keyLength  uint32
// }

// var hashSettings = HashSettings{
// 	memory:     64 * 1024,
// 	iterations: 3,
// 	threads:    4,
// 	saltLength: 16,
// 	keyLength:  32,
// }

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

// func ArgonHash(s string) string {
// 	salt := make([]byte, hashSettings.saltLength)
// 	_, _ = rand.Read(salt)

// 	hash := argon2.IDKey([]byte(s), salt, hashSettings.iterations,
// 		hashSettings.memory, hashSettings.threads, hashSettings.keyLength)

// 	var b strings.Builder
// 	b.Grow(defaultHashLen)

// 	b.WriteString("$argon2id$v=")
// 	b.WriteString(strconv.Itoa(argon2.Version))
// 	b.WriteString("$m=")
// 	b.WriteString(strconv.FormatUint(uint64(hashSettings.memory), 10))
// 	b.WriteString(",t=")
// 	b.WriteString(strconv.FormatUint(uint64(hashSettings.iterations), 10))
// 	b.WriteString(",p=")
// 	b.WriteString(strconv.FormatUint(uint64(hashSettings.threads), 10))
// 	b.WriteByte('$')
// 	b.WriteString(base64.RawStdEncoding.EncodeToString(salt))
// 	b.WriteByte('$')
// 	b.WriteString(base64.RawStdEncoding.EncodeToString(hash))

// 	return b.String()
// }

// func IsStringArgonHash(s, hash string) (bool, error) {
// 	vals := strings.Split(hash, "$")

// 	if len(vals) != 6 {
// 		return false, fmt.Errorf("invalid hash format")
// 	}

// 	memory, iterations, parallelism, err := parseParams(vals[3])
// 	if err != nil {
// 		return false, err
// 	}

// 	unb64salt, err := base64.RawStdEncoding.DecodeString(vals[4])
// 	if err != nil {
// 		return false, err
// 	}

// 	unb64hash, err := base64.RawStdEncoding.DecodeString(vals[5])
// 	if err != nil {
// 		return false, err
// 	}

// 	comparisonHash := argon2.IDKey([]byte(s), unb64salt, iterations, memory, parallelism, uint32(len(unb64hash)))

// 	return subtle.ConstantTimeCompare(unb64hash, comparisonHash) == 1, nil
// }

// func parseParams(paramsStr string) (m, t uint32, p uint8, err error) {
// 	parts := strings.Split(paramsStr, ",")
// 	if len(parts) != 3 {
// 		return 0, 0, 0, fmt.Errorf("invalid params format")
// 	}

// 	for _, part := range parts {
// 		val, err := strconv.ParseUint(part[2:], 10, 32)
// 		if err != nil {
// 			return 0, 0, 0, err
// 		}

// 		switch part[0:1] {
// 		case "m":
// 			m = uint32(val)
// 		case "t":
// 			t = uint32(val)
// 		case "p":
// 			p = uint8(val)
// 		}
// 	}
// 	return m, t, p, nil
// }

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

// func FNV1Hash(data []byte) string {
// 	h := fnv.New64a()
// 	h.Write(data)
// 	return strconv.FormatUint(h.Sum64(), 10)
// }

// func MakePatchFromTexts(v1, v2 string) string {
// 	dmp := diffmatchpatch.New()
// 	diffs := dmp.DiffMain(v1, v2, false)
// 	patch := dmp.PatchMake(v1, diffs)
// 	return dmp.PatchToText(patch)
// }

// func ApplyPatchToText(text, patch string) (string, error) {
// 	dmp := diffmatchpatch.New()
// 	loadedPatches, _ := dmp.PatchFromText(patch)
// 	newText, applies := dmp.PatchApply(loadedPatches, text)
// 	for i := range applies {
// 		if !applies[i] {
// 			return "", fmt.Errorf("Patch is partially or not applied")
// 		}
// 	}

// 	return newText, nil
// }
