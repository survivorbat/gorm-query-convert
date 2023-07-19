package gormqonvert

import (
	"github.com/ing-bank/gormtestutil"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeepGorm_Name_ReturnsExpectedName(t *testing.T) {
	t.Parallel()
	// Arrange
	plugin := New(CharacterConfig{})

	// Act
	result := plugin.Name()

	// Assert
	assert.Equal(t, "gormQonvert", result)
}

func TestDeepGorm_Initialize_RegistersCallback(t *testing.T) {
	t.Parallel()
	// Arrange
	db := gormtestutil.NewMemoryDatabase(t)
	plugin := New(CharacterConfig{})

	// Act
	err := plugin.Initialize(db)

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, db.Callback().Query().Get("gormQonvert:query"))
}
