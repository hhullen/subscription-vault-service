package supports

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	testKey = "support_test"
)

func TestParseDate(t *testing.T) {
	t.Parallel()

	t.Run("now", func(t *testing.T) {
		t.Parallel()
		tm, err := ParseDate("now")
		require.NotNil(t, tm)
		require.Nil(t, err)
	})

	t.Run("-months", func(t *testing.T) {
		t.Parallel()
		tm, err := ParseDate("-13")

		tb := time.Now().AddDate(0, 12, 0)

		require.True(t, tm.Before(tb))
		require.NotNil(t, tm)
		require.Nil(t, err)
	})

	t.Run("parse formats", func(t *testing.T) {
		t.Parallel()

		for _, tf := range dateFormats {
			tm, err := ParseDate(tf)
			require.NotNil(t, tm)
			require.Nil(t, err)

		}
	})

	t.Run("incorrect value", func(t *testing.T) {
		t.Parallel()

		tm, err := ParseDate("10--5")
		require.NotNil(t, err)
		require.Equal(t, tm, time.Time{})
	})

	t.Run("incorrect date", func(t *testing.T) {
		t.Parallel()

		tm, err := ParseDate("1998-20.02")
		require.NotNil(t, err)
		require.Equal(t, tm, time.Time{})
	})
}

func TestConcat(t *testing.T) {
	t.Parallel()

	require.NotPanics(t, func() {
		res := Concat("Power ", "of ", "Winx!")
		require.Equal(t, "Power of Winx!", res)
	})

}

func TestReadSecretFile(t *testing.T) {
	t.Parallel()

	tempDir := t.TempDir()

	filePath := filepath.Join(tempDir, testKey)

	err := os.WriteFile(filePath, []byte(testKey), 0644)
	require.NoError(t, err)

	t.Run("Ok", func(t *testing.T) {
		v, err := ReadSecretFile(filePath)
		require.Nil(t, err)
		require.Equal(t, testKey, v)
	})

	t.Run("error", func(t *testing.T) {
		v, err := ReadSecretFile("unexisted")
		require.NotNil(t, err)
		require.Equal(t, "", v)
	})
}

func TestIsInContainer(t *testing.T) {
	t.Parallel()

	require.Equal(t, IsInContainer(), false)
}

func TestMakeKVMessagesJSON(t *testing.T) {
	t.Parallel()

	b, _ := MakeKVMessagesJSON("key1", "val1", "key2", "val2", "KeyNoVal")
	data := map[string]string{}
	err := json.Unmarshal(b, &data)
	require.Nil(t, err)

	require.Equal(t, data["key1"], "val1")
	require.Equal(t, data["key2"], "val2")
	_, exists := data["KeyNoVal"]
	require.False(t, exists)
}
