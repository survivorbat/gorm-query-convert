package gormqonvert

import (
	"github.com/google/uuid"
	"github.com/ing-bank/gormtestutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestGormQonvert_Initialize_TriggersConversionCorrectly(t *testing.T) {
	t.Parallel()

	type ObjectA struct {
		ID   uuid.UUID
		Name string
		Age  int
		Date time.Time
	}

	defaultQuery := func(db *gorm.DB) *gorm.DB { return db }

	tests := map[string]struct {
		filter   map[string]any
		existing []ObjectA
		options  []Option
		query    func(*gorm.DB) *gorm.DB

		expected []ObjectA
	}{
		"nothing": {
			expected: []ObjectA{},
			query:    defaultQuery,
		},

		// Check if everything still works
		"simple where query": {
			filter: map[string]any{
				"name": "jessica",
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 46}, {Name: "amy", Age: 35}},
			expected: []ObjectA{{Name: "jessica", Age: 46}},
		},
		"more complex where query": {
			filter: map[string]any{
				"name": "jessica",
				"age":  53,
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "jessica", Age: 20}},
			expected: []ObjectA{{Name: "jessica", Age: 53}},
		},
		"multi-value where query": {
			filter: map[string]any{
				"name": []string{"jessica", "amy"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
		},
		"more complex multi-value where query": {
			filter: map[string]any{
				"name": []string{"jessica", "amy"},
				"age":  []int{53, 20},
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
			expected: []ObjectA{{Name: "jessica", Age: 53}, {Name: "amy", Age: 20}},
		},

		// On to the 'real' tests
		"greater or equal to value": {
			filter: map[string]any{
				"age": ">=30",
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
		},
		"greater than value": {
			filter: map[string]any{
				"age": ">30",
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "boris", Age: 31}},
		},
		"less or equal to value": {
			filter: map[string]any{
				"age": "<=30",
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}},
		},
		"less than value": {
			filter: map[string]any{
				"age": "<30",
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "jessica", Age: 29}},
		},
		"not equal to value": {
			filter: map[string]any{
				"age": "!=30",
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "jessica", Age: 29}, {Name: "boris", Age: 31}},
		},

		"not between certain values": {
			filter: map[string]any{
				"age": []string{"<30", ">35"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}, {Name: "ahmed", Age: 33}, {Name: "jochem", Age: 36}},
			expected: []ObjectA{{Name: "jessica", Age: 29}, {Name: "jochem", Age: 36}},
		},
		"not between certain values with another filter": {
			filter: map[string]any{
				"name": []string{"boris"},
				"age":  []string{">=30", "<=35"},
			},
			query:    defaultQuery,
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}, {Name: "ahmed", Age: 33}, {Name: "josh", Age: 34}, {Name: "jochem", Age: 36}},
			expected: []ObjectA{{Name: "boris", Age: 31}},
		},

		// With existing query
		"greater or equal to value with existing query": {
			filter: map[string]any{
				"age": ">=30",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "amy")
			},
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "amy", Age: 30}},
		},
		"greater than value with existing query": {
			filter: map[string]any{
				"age": ">30",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "boris")
			},
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}, {Name: "john", Age: 32}},
			expected: []ObjectA{{Name: "boris", Age: 31}},
		},
		"less or equal to value with existing query": {
			filter: map[string]any{
				"age": "<=30",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "jessica")
			},
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "jessica", Age: 29}},
		},
		"less than value with existing query": {
			filter: map[string]any{
				"age": "<30",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "jessica")
			},
			existing: []ObjectA{{Name: "john", Age: 28}, {Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "jessica", Age: 29}},
		},
		"not equal to value with existing query": {
			filter: map[string]any{
				"age": "!=30",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "boris")
			},
			existing: []ObjectA{{Name: "jessica", Age: 29}, {Name: "amy", Age: 30}, {Name: "boris", Age: 31}},
			expected: []ObjectA{{Name: "boris", Age: 31}},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			db := gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))
			_ = db.AutoMigrate(&ObjectA{})

			config := CharacterConfig{
				GreaterThanPrefix:      ">",
				GreaterOrEqualToPrefix: ">=",
				LessThanPrefix:         "<",
				LessOrEqualToPrefix:    "<=",
				NotEqualToPrefix:       "!=",
			}

			plugin := New(config, testData.options...)

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			// Act
			err := db.Use(plugin)

			// Assert
			assert.Nil(t, err)

			var actual []ObjectA
			err = testData.query(db).Where(testData.filter).Find(&actual).Error
			assert.Nil(t, err)

			assert.Equal(t, testData.expected, actual)
		})
	}
}

func TestGormQonvert_Initialize_TriggersConversionCorrectlyWithSetting(t *testing.T) {
	t.Parallel()

	type ObjectB struct {
		Age int
	}

	tests := map[string]struct {
		filter   map[string]any
		query    func(*gorm.DB) *gorm.DB
		existing []ObjectB
		expected []ObjectB
	}{
		"conversion with query set to true": {
			filter: map[string]any{
				"age": ">=40",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, true)
			},
			existing: []ObjectB{{Age: 50}},
			expected: []ObjectB{{Age: 50}},
		},
		"conversion with query set to false": {
			filter: map[string]any{
				"age": ">=40",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, false)
			},
			existing: []ObjectB{{Age: 50}},
			expected: []ObjectB{},
		},
		"conversion with query set to random value": {
			filter: map[string]any{
				"age": ">=40",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Set(tagName, "yes")
			},
			existing: []ObjectB{{Age: 50}},
			expected: []ObjectB{},
		},
		"conversion with query unset": {
			filter: map[string]any{
				"age": ">=40",
			},
			query: func(db *gorm.DB) *gorm.DB {
				return db
			},
			existing: []ObjectB{{Age: 50}},
			expected: []ObjectB{},
		},
	}

	for name, testData := range tests {
		testData := testData
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			db := gormtestutil.NewMemoryDatabase(t, gormtestutil.WithName(t.Name()))
			_ = db.AutoMigrate(&ObjectB{})

			config := CharacterConfig{
				GreaterThanPrefix:      ">",
				GreaterOrEqualToPrefix: ">=",
				LessThanPrefix:         "<",
				LessOrEqualToPrefix:    "<=",
				NotEqualToPrefix:       "!=",
			}

			plugin := New(config, SettingOnly())

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			db = testData.query(db)

			// Act
			err := db.Use(plugin)

			// Assert
			assert.Nil(t, err)

			var actual []ObjectB
			err = db.Where(testData.filter).Find(&actual).Error
			assert.Nil(t, err)

			assert.Equal(t, testData.expected, actual)
		})
	}
}
