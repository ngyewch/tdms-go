package tdms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	{
		_, err := ObjectPathFromString("")
		var invalidPathError *InvalidPathError
		assert.ErrorAs(t, err, &invalidPathError)
	}
	{
		objectPath, err := ObjectPathFromString("/")
		if assert.NoError(t, err) {
			assert.True(t, objectPath.IsRoot())
			assert.False(t, objectPath.IsGroup())
			assert.False(t, objectPath.IsChannel())
			assert.Empty(t, objectPath.Group)
			assert.Empty(t, objectPath.Channel)
		}
	}
	{
		_, err := ObjectPathFromString("/group 1")
		var invalidPathError *InvalidPathError
		assert.ErrorAs(t, err, &invalidPathError)
	}
	{
		_, err := ObjectPathFromString("/'group 1")
		var invalidPathError *InvalidPathError
		assert.ErrorAs(t, err, &invalidPathError)
	}
	{
		objectPath, err := ObjectPathFromString("/'group 1'")
		if assert.NoError(t, err) {
			assert.False(t, objectPath.IsRoot())
			assert.True(t, objectPath.IsGroup())
			assert.False(t, objectPath.IsChannel())
			assert.Equal(t, objectPath.Group, "group 1")
			assert.Empty(t, objectPath.Channel)
		}
	}
	{
		_, err := ObjectPathFromString("/'group 1'/channel 1")
		var invalidPathError *InvalidPathError
		assert.ErrorAs(t, err, &invalidPathError)
	}
	{
		_, err := ObjectPathFromString("/'group 1'/'channel 1")
		var invalidPathError *InvalidPathError
		assert.ErrorAs(t, err, &invalidPathError)
	}
	{
		objectPath, err := ObjectPathFromString("/'group 1'/'channel 1'")
		if assert.NoError(t, err) {
			assert.False(t, objectPath.IsRoot())
			assert.False(t, objectPath.IsGroup())
			assert.True(t, objectPath.IsChannel())
			assert.Equal(t, objectPath.Group, "group 1")
			assert.Equal(t, objectPath.Channel, "channel 1")
		}
	}
}
