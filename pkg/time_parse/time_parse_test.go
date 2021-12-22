package time_parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRFC3339ToCSTLayout(t *testing.T) {
	cstTimeString, err := RFC3339ToCSTLayout("2020-11-08T08:18:46+08:00")
	assert.NoError(t, err)
	assert.Equal(t, "2020-11-08 08:18:46", cstTimeString)
}

func TestCSTLayoutString(t *testing.T) {
	t.Log(CSTLayoutString())
}

func TestCSTLayoutStringToUnix(t *testing.T) {
	testTime, err := ParseCSTInLocation("2020-01-24 21:11:11")
	assert.NoError(t, err)

	unixTime, err := CSTLayoutStringToUnix("2020-01-24 21:11:11")
	assert.NoError(t, err)

	assert.Equal(t, unixTime, testTime.Unix())
}

func TestUnixToCSTLayoutString(t *testing.T) {
	cstLayout, err := UnixToCSTLayoutString(1579871471)
	assert.NoError(t, err)
	assert.Equal(t, cstLayout, "2020-01-24 21:11:11")
}

func TestGMTLayoutString(t *testing.T) {
	t.Log(GMTLayoutString())
}

func TestParseGMTInLocation(t *testing.T) {
	testTime, err := ParseGMTInLocation("Wed, 22 Dec 2021 10:07:29 GMT")
	assert.NoError(t, err)
	t.Log(testTime.String())
}
