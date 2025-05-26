package gormqonvert

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ing-bank/gormtestutil"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestGormQonvert_Initialize_TriggersConversionCorrectly(t *testing.T) {
	t.Parallel()

	type ObjectA struct {
		ID   uuid.UUID `gorm:"type:uuid"`
		Name string
		Age  int
		Date time.Time
	}

	defaultQuery := func(db *gorm.DB) *gorm.DB { return db }

	tests := map[string]struct {
		filter   []map[string]any
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
			filter: []map[string]any{{
				"name": "jessica",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("0abeed86-f60f-4fcc-988b-d240e5f2fb5f"), Name: "jessica", Age: 46},
				{Name: "amy", Age: 35},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("0abeed86-f60f-4fcc-988b-d240e5f2fb5f"), Name: "jessica", Age: 46},
			},
		},
		"more complex where query": {
			filter: []map[string]any{{
				"name": "jessica",
				"age":  53,
			},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("b0bf3c1f-ca33-45fb-a2d7-05281efafbae"), Name: "jessica", Age: 53},
				{ID: uuid.New(), Name: "jessica", Age: 20},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("b0bf3c1f-ca33-45fb-a2d7-05281efafbae"), Name: "jessica", Age: 53},
			},
		},
		"multi-value where query": {
			filter: []map[string]any{{
				"name": []string{"jessica", "amy"},
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("b0bf3c1f-ca33-45fb-a2d7-05281efafbae"), Name: "jessica", Age: 53},
				{ID: uuid.MustParse("436f9485-6209-43cf-b1c9-a49382f7df20"), Name: "amy", Age: 20},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("b0bf3c1f-ca33-45fb-a2d7-05281efafbae"), Name: "jessica", Age: 53},
				{ID: uuid.MustParse("436f9485-6209-43cf-b1c9-a49382f7df20"), Name: "amy", Age: 20},
			},
		},
		"more complex multi-value where query": {
			filter: []map[string]any{{
				"name": []string{"jessica", "amy"},
				"age":  []int{53, 20}},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("b0bf3c1f-ca33-45fb-a2d7-05281efafbae"), Name: "jessica", Age: 53},
				{ID: uuid.MustParse("436f9485-6209-43cf-b1c9-a49382f7df20"), Name: "amy", Age: 20},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("b0bf3c1f-ca33-45fb-a2d7-05281efafbae"), Name: "jessica", Age: 53},
				{ID: uuid.MustParse("436f9485-6209-43cf-b1c9-a49382f7df20"), Name: "amy", Age: 20},
			},
		},

		// On to the 'real' tests
		"greater or equal to value": {
			filter: []map[string]any{{
				"age": ">=30",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
		},
		"greater than value": {
			filter: []map[string]any{{
				"age": ">30",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
		},
		"less or equal to value": {
			filter: []map[string]any{{
				"age": "<=30",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
			},
		},
		"less than value": {
			filter: []map[string]any{{
				"age": "<30",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
			},
		},
		"like value": {
			filter: []map[string]any{{
				"name": "~%i%",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
		},
		"not like value": {
			filter: []map[string]any{{
				"name": "!~%a%",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
		},
		"not equal to value": {
			filter: []map[string]any{{
				"age": "!=30",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
		},
		"between certain values": {
			filter: []map[string]any{{
				"age": []string{">30"},
			}, {
				"age": []string{"<35"},
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
				{ID: uuid.MustParse("b89241fa-cec4-411a-a927-99d784fe3375"), Name: "ahmed", Age: 33},
				{ID: uuid.MustParse("699204f0-26f0-4e02-9e25-b73ac0b2300b"), Name: "jochem", Age: 36},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
				{ID: uuid.MustParse("b89241fa-cec4-411a-a927-99d784fe3375"), Name: "ahmed", Age: 33}},
		},
		"not between certain values": {
			filter: []map[string]any{{
				"age": []string{"<30", ">35"},
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
				{ID: uuid.MustParse("b89241fa-cec4-411a-a927-99d784fe3375"), Name: "ahmed", Age: 33},
				{ID: uuid.MustParse("699204f0-26f0-4e02-9e25-b73ac0b2300b"), Name: "jochem", Age: 36},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("699204f0-26f0-4e02-9e25-b73ac0b2300b"), Name: "jochem", Age: 36},
			},
		},
		"not between certain values with another filter": {
			filter: []map[string]any{{
				"name": []string{"boris"},
				"age":  []string{">=30"},
			}, {
				"age": "<=35",
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{Name: "boris", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
				{ID: uuid.MustParse("b89241fa-cec4-411a-a927-99d784fe3375"), Name: "ahmed", Age: 33},
				{ID: uuid.MustParse("4a2e5a60-044f-40f6-a501-74d1b18e52b3"), Name: "josh", Age: 34},
				{ID: uuid.MustParse("699204f0-26f0-4e02-9e25-b73ac0b2300b"), Name: "jochem", Age: 36},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
		},
		"not certain values": {
			filter: []map[string]any{
				{"age": "!=29"},
				{"age": "!=30"},
				{"age": "!=33"},
				{"age": "!=34"},
				{"age": "!=36"},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
				{ID: uuid.MustParse("b89241fa-cec4-411a-a927-99d784fe3375"), Name: "ahmed", Age: 33},
				{ID: uuid.MustParse("4a2e5a60-044f-40f6-a501-74d1b18e52b3"), Name: "josh", Age: 34},
				{ID: uuid.MustParse("699204f0-26f0-4e02-9e25-b73ac0b2300b"), Name: "jochem", Age: 36}},
			expected: []ObjectA{
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
		},
		"after date": {
			filter: []map[string]any{
				{"date": fmt.Sprintf(">=%v", time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("98472426-bcc3-4939-a9c3-03d875ad3014"), Name: "joris", Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
				{ID: uuid.MustParse("ccec8651-438f-461c-9ea2-4915f99be39c"), Name: "dane", Date: time.Date(2022, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("98472426-bcc3-4939-a9c3-03d875ad3014"), Name: "joris", Date: time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC)},
			},
		},
		"like multiple values": {
			filter: []map[string]any{{
				"name": []string{"~%ssica", "~a_y"},
			}},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
			},
		},
		"not like multiple values": {
			filter: []map[string]any{
				{
					"name": []string{"!~%ss%"},
				}, {
					"name": []string{"!~%ris%"},
				},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
			},
		},
		"almost everything really": {
			filter: []map[string]any{
				{
					"name": []string{"!~%bor%"},
					"age":  []string{"<30"},
				}, {
					"name": []string{"~%ss%"},
					"age":  []string{"!=30"},
				},
				{
					"age": "!=28",
				},
			},
			query: defaultQuery,
			existing: []ObjectA{
				{ID: uuid.MustParse("b4c582da-05ff-455a-8c7c-4e6b6fb3b858"), Name: "jessica", Age: 28},
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("b79f4afb-0e0b-41f7-aa2e-65a7b6d1cd0f"), Name: "jessica", Age: 30},
				{ID: uuid.MustParse("0d48e381-e9f5-4a4b-b24e-1fe4d70db755"), Name: "boris", Age: 25},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
			},
		},
		// With existing query
		"greater or equal to value with existing query": {
			filter: []map[string]any{{
				"age": ">=30",
			}},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "amy")
			},
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30}},
		},
		"greater than value with existing query": {
			filter: []map[string]any{{
				"age": ">30",
			}},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "boris")
			},
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
				{ID: uuid.MustParse("8c4de24b-03db-40d9-bc96-a589e88a4462"), Name: "john", Age: 32},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
		},
		"less or equal to value with existing query": {
			filter: []map[string]any{{
				"age": "<=30",
			}},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "jessica")
			},
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
			},
		},
		"less than value with existing query": {
			filter: []map[string]any{{
				"age": "<30",
			}},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "jessica")
			},
			existing: []ObjectA{
				{Name: "john", Age: 28},
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
			},
		},
		"not equal to value with existing query": {
			filter: []map[string]any{{
				"age": "!=30",
			}},
			query: func(db *gorm.DB) *gorm.DB {
				return db.Where("name = ?", "boris")
			},
			existing: []ObjectA{
				{ID: uuid.MustParse("49d3c60b-48e0-4bc8-a144-b0d823cd1373"), Name: "jessica", Age: 29},
				{ID: uuid.MustParse("d8e2b086-21d4-4671-b1e8-cedc97a804a6"), Name: "amy", Age: 30},
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
			expected: []ObjectA{
				{ID: uuid.MustParse("2709499e-8666-4775-959b-24289a6eabff"), Name: "boris", Age: 31},
			},
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
				LikePrefix:             "~",
				NotLikePrefix:          "!~",
			}

			plugin := New(config, testData.options...)

			if err := db.CreateInBatches(testData.existing, 10).Error; err != nil {
				t.Error(err)
				t.FailNow()
			}

			// Act
			err := db.Use(plugin)

			// Assert
			assert.NoError(t, err)

			var actual []ObjectA

			query := testData.query(db)
			for _, filter := range testData.filter {
				query = query.Where(filter)
			}

			assert.NoError(t, query.Find(&actual).Error)

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
				LikePrefix:             "~",
				NotLikePrefix:          "!~",
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
			assert.NoError(t, err)

			var actual []ObjectB
			err = db.Where(testData.filter).Find(&actual).Error
			assert.NoError(t, err)

			assert.Equal(t, testData.expected, actual)
		})
	}
}
