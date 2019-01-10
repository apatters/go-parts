package parts_test

import (
	"fmt"
	"testing"

	"github.com/apatters/go-parts"
	"github.com/stretchr/testify/assert"
)

const (
	regularTypeStr = "regular"
	dirTypeStr     = "dir"
	execTypeStr    = "executable"
	allTypeStr     = "all"
)

type modeDatum struct {
	Mode       parts.FileMode
	StringRepr string
	Type       string
}

var (
	modeData = []modeDatum{
		{parts.ModeRegular | 0666, "frw-rw-rw-", regularTypeStr},
		{parts.ModeDir | 0775, "drwxrwxr-x", dirTypeStr},
		{parts.ModeRegular | 0775, "frwxrwxr-x", execTypeStr},
		{parts.ModeType | parts.ModePerm, "dLDpSfrwxrwxrwx", allTypeStr},
	}
)

func TestModeString(t *testing.T) {
	for _, datum := range modeData {
		t.Logf("Mode: %s", datum.Mode)
		t.Logf("StringRep: %s", datum.StringRepr)
		assert.EqualValues(t, datum.StringRepr, datum.Mode.String())
	}
}

func TestModeIsFunctions(t *testing.T) {
	for _, datum := range modeData {
		t.Logf("Mode: %s", datum.Mode)
		switch {
		case datum.Type == regularTypeStr || datum.Type == allTypeStr:
			assert.True(t, datum.Mode.IsRegular())
		case datum.Type == dirTypeStr || datum.Type == allTypeStr:
			assert.True(t, datum.Mode.IsDir())
		case datum.Type == execTypeStr || datum.Type == allTypeStr:
			assert.True(t, datum.Mode.IsExecutable())
		default:
			panic(fmt.Sprintf("Unknown datum type '%s", datum.Type))
		}
	}
}

func TestModePerm(t *testing.T) {
	for _, datum := range modeData {
		t.Logf("Mode: %s", datum.Mode)
		perm := datum.Mode.Perm()
		t.Logf("perm: = %s", perm)
		assert.EqualValues(t, "-"+datum.StringRepr[len(datum.StringRepr)-9:], perm.String())
	}
}
