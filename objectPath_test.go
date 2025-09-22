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

			assert.Equal(t, "/", objectPath.String())
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
			assert.Equal(t, "group 1", objectPath.Group)
			assert.Empty(t, objectPath.Channel)

			assert.Equal(t, "/'group 1'", objectPath.String())
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
			assert.Equal(t, "group 1", objectPath.Group)
			assert.Equal(t, "channel 1", objectPath.Channel)

			assert.Equal(t, "/'group 1'/'channel 1'", objectPath.String())
		}
	}
	{
		_, err := ObjectPathFromString("/'group 1'/'channel 1'/'item 1'")
		var invalidPathError *InvalidPathError
		assert.ErrorAs(t, err, &invalidPathError)
	}
	{
		_, err := ObjectPathFromString("/'group 1'/'bob's channel'")
		var invalidPathError *InvalidPathError
		assert.ErrorAs(t, err, &invalidPathError)
	}
	{
		objectPath, err := ObjectPathFromString("/'group 1'/'bob''s channel'")
		if assert.NoError(t, err) {
			assert.False(t, objectPath.IsRoot())
			assert.False(t, objectPath.IsGroup())
			assert.True(t, objectPath.IsChannel())
			assert.Equal(t, "group 1", objectPath.Group)
			assert.Equal(t, "bob's channel", objectPath.Channel)

			assert.Equal(t, "/'group 1'/'bob''s channel'", objectPath.String())
		}
	}
}
