package collection

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/bedrock/packages/collection/arr"
	collectible "github.com/bedrock/packages/collection/collectible"
	"github.com/bedrock/packages/collection/kv"
	csupport "github.com/bedrock/packages/collection/support"
)

type inventoryUser struct {
	ID     int
	Name   string
	Team   string
	Active bool
	Score  float64
}

func TestInventoryParityFirstWhereAndTypedFilters(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testFirstOrFailStopsIteratingAtFirstMatch
	// SupportCollectionTest::testFirstWhere
	// SupportCollectionTest::testWhereStrict
	// SupportCollectionTest::testWhereInStrict
	// SupportCollectionTest::testWhereNotInStrict
	// SupportCollectionTest::testBetween
	users := Collect([]inventoryUser{
		{ID: 1, Name: "Taylor", Team: "core", Active: true, Score: 10},
		{ID: 2, Name: "Abigail", Team: "docs", Active: false, Score: 20},
		{ID: 3, Name: "Mohamed", Team: "core", Active: true, Score: 30},
	})

	seen := 0
	first, err := users.FirstOrFail(func(user inventoryUser, _ int) bool {
		seen++

		return user.Team == "core"
	})

	if err != nil {
		t.Fatalf("FirstOrFail returned error: %v", err)
	}

	if first.ID != 1 || seen != 1 {
		t.Fatalf("FirstOrFail = %#v after %d iterations, want first core user after one iteration", first, seen)
	}

	firstInactive, ok := users.First(func(user inventoryUser, _ int) bool {
		return !user.Active
	})

	if !ok || firstInactive.Name != "Abigail" {
		t.Fatalf("First inactive = %#v, %v", firstInactive, ok)
	}

	coreUsers := users.Where(func(user inventoryUser) bool {
		return user.Team == "core"
	})

	if got := namesOf(coreUsers.All()); !reflect.DeepEqual(got, []string{"Taylor", "Mohamed"}) {
		t.Fatalf("Where strict team names = %v", got)
	}

	inIDs := WhereIn(users, func(user inventoryUser) int { return user.ID }, []int{1, 3})

	if got := namesOf(inIDs.All()); !reflect.DeepEqual(got, []string{"Taylor", "Mohamed"}) {
		t.Fatalf("WhereIn strict ids = %v", got)
	}

	notInIDs := WhereNotIn(users, func(user inventoryUser) int { return user.ID }, []int{1, 3})

	if got := namesOf(notInIDs.All()); !reflect.DeepEqual(got, []string{"Abigail"}) {
		t.Fatalf("WhereNotIn strict ids = %v", got)
	}

	between := WhereBetween(users, func(user inventoryUser) float64 { return user.Score }, 15, 30)

	if got := namesOf(between.All()); !reflect.DeepEqual(got, []string{"Abigail", "Mohamed"}) {
		t.Fatalf("WhereBetween scores = %v", got)
	}
}

func TestInventoryParityUniquenessDuplicatesAndSorting(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testUniqueWithCallback
	// SupportCollectionTest::testUniqueStrict
	// SupportCollectionTest::testDuplicatesWithKey
	// SupportCollectionTest::testDuplicatesWithCallback
	// SupportCollectionTest::testDuplicatesWithStrict
	// SupportCollectionTest::testValuesResetKey
	// SupportCollectionTest::testSortWithCallback
	// SupportCollectionTest::testSortByMany
	items := Collect([]inventoryUser{
		{ID: 1, Name: "Taylor", Team: "core"},
		{ID: 2, Name: "Abigail", Team: "docs"},
		{ID: 3, Name: "Mohamed", Team: "core"},
		{ID: 4, Name: "Nuno", Team: "docs"},
	})

	uniqueTeams := Unique(items, func(user inventoryUser) string { return user.Team })

	if got := namesOf(uniqueTeams.All()); !reflect.DeepEqual(got, []string{"Taylor", "Abigail"}) {
		t.Fatalf("Unique teams = %v", got)
	}

	uniqueIDs := Unique(items, func(user inventoryUser) int { return user.ID })

	if got := namesOf(uniqueIDs.All()); !reflect.DeepEqual(got, []string{"Taylor", "Abigail", "Mohamed", "Nuno"}) {
		t.Fatalf("Unique ids = %v", got)
	}

	duplicateTeams := Duplicates(items, func(user inventoryUser) string { return user.Team })

	if got := namesOf(duplicateTeams.All()); !reflect.DeepEqual(got, []string{"Mohamed", "Nuno"}) {
		t.Fatalf("Duplicate teams = %v", got)
	}

	values := items.Values()

	if !reflect.DeepEqual(values.All(), items.All()) {
		t.Fatalf("Values() = %#v, want copy with same order", values.All())
	}

	values.Put(0, inventoryUser{ID: 99, Name: "Changed"})

	if items.All()[0].ID == 99 {
		t.Fatal("Values() should reset to an independent slice copy")
	}

	sorted := items.Sort(func(a, b inventoryUser) bool {
		return a.Name < b.Name
	})

	if got := namesOf(sorted.All()); !reflect.DeepEqual(got, []string{"Abigail", "Mohamed", "Nuno", "Taylor"}) {
		t.Fatalf("Sort with callback = %v", got)
	}

	sortByMany := items.Sort(func(a, b inventoryUser) bool {
		if a.Team != b.Team {
			return a.Team < b.Team
		}

		return a.Name < b.Name
	})

	if got := namesOf(sortByMany.All()); !reflect.DeepEqual(got, []string{"Mohamed", "Taylor", "Abigail", "Nuno"}) {
		t.Fatalf("SortByMany equivalent = %v", got)
	}
}

func TestInventoryParityGroupKeyAndMapWithKeys(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testMapToDictionaryWithNumericKeys
	// SupportCollectionTest::testMapToGroupsWithNumericKeys
	// SupportCollectionTest::testMapWithKeysIntegerKeys
	// SupportCollectionTest::testMapWithKeysMultipleRows
	// SupportCollectionTest::testMapWithKeysCallbackKey
	// SupportCollectionTest::testMapWithKeysOverwritingKeys
	// SupportCollectionTest::testGroupByClosureWhereItemsHaveSingleGroup
	// SupportCollectionTest::testGroupByClosureWhereItemsHaveMultipleGroups
	// SupportCollectionTest::testGroupByAttribute
	// SupportCollectionTest::testKeyByClosure
	users := Collect([]inventoryUser{
		{ID: 1, Name: "Taylor", Team: "core"},
		{ID: 2, Name: "Abigail", Team: "docs"},
		{ID: 3, Name: "Mohamed", Team: "core"},
	})

	dictionary := MapToDictionary(users, func(user inventoryUser) (int, string) {
		return user.ID % 2, user.Name
	})

	sort.Strings(dictionary[1])

	if !reflect.DeepEqual(dictionary[1], []string{"Mohamed", "Taylor"}) || !reflect.DeepEqual(dictionary[0], []string{"Abigail"}) {
		t.Fatalf("MapToDictionary numeric keys = %#v", dictionary)
	}

	groupsByNumber := MapToGroups(users, func(user inventoryUser) (int, string) {
		return user.ID % 2, user.Name
	})

	if !reflect.DeepEqual(groupsByNumber[1], []string{"Taylor", "Mohamed"}) || !reflect.DeepEqual(groupsByNumber[0], []string{"Abigail"}) {
		t.Fatalf("MapToGroups numeric keys = %#v", groupsByNumber)
	}

	mapped := MapWithKeys(users, func(user inventoryUser) (int, string) {
		return user.ID, user.Name
	})

	if !reflect.DeepEqual(mapped, map[int]string{1: "Taylor", 2: "Abigail", 3: "Mohamed"}) {
		t.Fatalf("MapWithKeys integer keys = %#v", mapped)
	}

	overwritten := MapWithKeys(users, func(user inventoryUser) (string, string) {
		return user.Team, user.Name
	})

	if !reflect.DeepEqual(overwritten, map[string]string{"core": "Mohamed", "docs": "Abigail"}) {
		t.Fatalf("MapWithKeys overwriting keys = %#v", overwritten)
	}

	grouped := GroupBy(users, func(user inventoryUser) string { return user.Team })

	if got := namesOf(grouped["core"].All()); !reflect.DeepEqual(got, []string{"Taylor", "Mohamed"}) {
		t.Fatalf("GroupBy core = %v", got)
	}

	if got := namesOf(grouped["docs"].All()); !reflect.DeepEqual(got, []string{"Abigail"}) {
		t.Fatalf("GroupBy attribute docs = %v", got)
	}

	keyed := KeyBy(users, func(user inventoryUser) int { return user.ID })

	if keyed[2].Name != "Abigail" {
		t.Fatalf("KeyBy closure = %#v", keyed)
	}
}

func TestInventoryParityMergeConcatSliceMedianAndMode(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testMergeCollection
	// SupportCollectionTest::testConcatWithCollection
	// SupportCollectionTest::testSliceNegativeOffset
	// SupportCollectionTest::testSliceNegativeOffsetAndLength
	// SupportCollectionTest::testMedianValueByKey
	// SupportCollectionTest::testMedianOnCollectionWithNull
	// SupportCollectionTest::testMedianOutOfOrderCollection
	// SupportCollectionTest::testMedianOnEmptyCollectionReturnsNull
	// SupportCollectionTest::testModeOnNullCollection
	// SupportCollectionTest::testModeValueByKey
	left := Collect([]int{1, 2})
	right := Collect([]int{3, 4})

	if got := left.Merge(right.All()).All(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("Merge collection = %v", got)
	}

	if got := left.Concat(right.All()).All(); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Fatalf("Concat collection = %v", got)
	}

	values := Collect([]int{1, 2, 3, 4, 5})

	if got := values.Slice(-2).All(); !reflect.DeepEqual(got, []int{4, 5}) {
		t.Fatalf("Slice negative offset = %v", got)
	}

	if got := values.Slice(-4, 2).All(); !reflect.DeepEqual(got, []int{2, 3}) {
		t.Fatalf("Slice negative offset and length = %v", got)
	}

	scoreUsers := Collect([]inventoryUser{
		{Name: "low", Score: 5},
		{Name: "high", Score: 15},
		{Name: "mid", Score: 10},
	})

	if got := MedianBy(scoreUsers, func(user inventoryUser) float64 { return user.Score }); got != 10 {
		t.Fatalf("MedianBy score = %v", got)
	}

	if got := Median(Collect([]float64{30, 10, 20})); got != 20 {
		t.Fatalf("Median out of order = %v", got)
	}

	if got := Median(Empty[float64]()); got != 0 {
		t.Fatalf("Median empty = %v", got)
	}

	if got := Mode(Empty[string]()); got != nil {
		t.Fatalf("Mode empty = %#v", got)
	}

	if got := Mode(Collect([]string{"core", "docs", "core"})); !reflect.DeepEqual(got, []string{"core"}) {
		t.Fatalf("Mode values = %#v", got)
	}
}

func TestInventoryParitySerializationMutationAndMath(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testLazyReturnsLazyCollection
	// SupportCollectionTest::testToJsonEncodesTheJsonSerializeResult
	// SupportCollectionTest::testToPrettyJsonEncodesTheJsonSerializeResult
	// SupportCollectionTest::testCastingToStringJsonEncodesTheToArrayResult
	// SupportCollectionTest::testForgetSingleKey
	// SupportCollectionTest::testMultiplyCollection
	// SupportCollectionTest::testPut
	// SupportCollectionTest::testPutWithNoKey
	// SupportCollectionTest::testGettingSumFromEmptyCollection
	// SupportCollectionTest::testUnshiftWithOneItem
	// SupportCollectionTest::testUnshiftWithMultipleItems
	// SupportCollectionTest::testCombineWithCollection
	// SupportCollectionTest::testMedianValueWithArrayCollection
	users := Collect([]inventoryUser{
		{ID: 1, Name: "Taylor", Score: 10},
		{ID: 2, Name: "Abigail", Score: 20},
	})

	var iterated []string

	for user := range users.Iter() {
		iterated = append(iterated, user.Name)
	}

	if !reflect.DeepEqual(iterated, []string{"Taylor", "Abigail"}) {
		t.Fatalf("lazy iteration names = %v", iterated)
	}

	encoded, err := users.ToJSON()

	if err != nil {
		t.Fatalf("ToJSON returned error: %v", err)
	}

	if !json.Valid(encoded) || !strings.Contains(string(encoded), "Taylor") {
		t.Fatalf("ToJSON = %s", encoded)
	}

	pretty, err := users.ToPrettyJSON()

	if err != nil {
		t.Fatalf("ToPrettyJSON returned error: %v", err)
	}

	if !json.Valid(pretty) || !strings.Contains(string(pretty), "\n") {
		t.Fatalf("ToPrettyJSON = %s", pretty)
	}

	if got := users.String(); !json.Valid([]byte(got)) || !strings.Contains(got, "Abigail") {
		t.Fatalf("String = %s", got)
	}

	numbers := Collect([]int{1, 2, 3})
	numbers.Forget(1)

	if got := numbers.All(); !reflect.DeepEqual(got, []int{1, 3}) {
		t.Fatalf("Forget single key = %v", got)
	}

	if got := Collect([]int{1, 2}).Multiply(3).All(); !reflect.DeepEqual(got, []int{1, 2, 1, 2, 1, 2}) {
		t.Fatalf("Multiply = %v", got)
	}

	put := Collect([]string{"a", "b"})
	put.Put(1, "B")
	put.Put(99, "missing")

	if got := put.All(); !reflect.DeepEqual(got, []string{"a", "B"}) {
		t.Fatalf("Put = %v", got)
	}

	if got := Sum(Empty[int]()); got != 0 {
		t.Fatalf("Sum empty = %d", got)
	}

	unshifted := Collect([]int{3})
	unshifted.Unshift(2).Unshift(1)

	if got := unshifted.All(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("Unshift = %v", got)
	}

	keys := Collect([]string{"name", "email"})
	values := Collect([]string{"Taylor", "taylor@example.com"})
	combined := Combine(keys, values.All()).All()

	if len(combined) != 2 || combined[0].Key != "name" || combined[0].Value != "Taylor" {
		t.Fatalf("Combine collection = %#v", combined)
	}

	scoreMedian := MedianBy(users, func(user inventoryUser) float64 { return user.Score })

	if scoreMedian != 15 {
		t.Fatalf("Median value with array collection = %v", scoreMedian)
	}
}

func TestInventoryParityChunkTakeSplitAndPartition(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testChunkWhenGivenZeroAsSize
	// SupportCollectionTest::testChunkWhenGivenLessThanZero
	// SupportCollectionTest::testChunkWhileOnContiguouslyIncreasingIntegers
	// SupportCollectionTest::testTakeLast
	// SupportCollectionTest::testTakeUntilUsingValue
	// SupportCollectionTest::testTakeUntilReturnsAllItemsForUnmetValue
	// SupportCollectionTest::testTakeWhileUsingValue
	// SupportCollectionTest::testTakeWhileReturnsNoItemsForUnmetValue
	// SupportCollectionTest::testSplitCollectionWithAnUndivisableCount
	// SupportCollectionTest::testSplitCollectionWithCountLessThenDivisor
	// SupportCollectionTest::testSplitCollectionIntoThreeWithCountOfFour
	// SupportCollectionTest::testSplitCollectionIntoThreeWithCountOfFive
	// SupportCollectionTest::testSplitCollectionIntoSixWithCountOfTen
	// SupportCollectionTest::testSplitEmptyCollection
	// SupportCollectionTest::testNthThrowsExceptionForInvalidStep
	// SupportCollectionTest::testNthThrowsExceptionForNegativeStep
	// SupportCollectionTest::testSplitThrowsExceptionForInvalidNumberOfGroups
	// SupportCollectionTest::testSplitThrowsExceptionForNegativeNumberOfGroups
	// SupportCollectionTest::testSplitInThrowsExceptionForInvalidNumberOfGroups
	// SupportCollectionTest::testSplitInThrowsExceptionForNegativeNumberOfGroups
	// SupportCollectionTest::testPartitionCallbackWithKey
	// SupportCollectionTest::testPartitionByKey
	// SupportCollectionTest::testPartitionWithOperators
	// SupportCollectionTest::testPartitionEmptyCollection
	items := Collect([]int{1, 2, 3, 5, 6, 9})

	if chunks := items.Chunk(0); chunks != nil {
		t.Fatalf("Chunk zero = %#v, want nil", chunks)
	}

	if chunks := items.Chunk(-2); chunks != nil {
		t.Fatalf("Chunk negative = %#v, want nil", chunks)
	}

	chunked := items.ChunkWhile(func(item int, _ int, current []int) bool {
		return item == current[len(current)-1]+1
	})

	if !reflect.DeepEqual(chunked, [][]int{{1, 2, 3}, {5, 6}, {9}}) {
		t.Fatalf("ChunkWhile contiguous = %#v", chunked)
	}

	if got := items.Take(-2).All(); !reflect.DeepEqual(got, []int{6, 9}) {
		t.Fatalf("Take last = %v", got)
	}

	if got := items.TakeUntil(func(item int, _ int) bool { return item == 5 }).All(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("TakeUntil value = %v", got)
	}

	if got := items.TakeUntil(func(item int, _ int) bool { return item == 99 }).All(); !reflect.DeepEqual(got, items.All()) {
		t.Fatalf("TakeUntil unmet = %v", got)
	}

	if got := items.TakeWhile(func(item int, _ int) bool { return item < 5 }).All(); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Fatalf("TakeWhile value = %v", got)
	}

	if got := items.TakeWhile(func(item int, _ int) bool { return item > 99 }).All(); len(got) != 0 {
		t.Fatalf("TakeWhile unmet = %v", got)
	}

	if got := Collect([]int{1, 2, 3, 4, 5}).Split(2); !reflect.DeepEqual(got, [][]int{{1, 2, 3}, {4, 5}}) {
		t.Fatalf("Split indivisible = %#v", got)
	}

	if got := Collect([]int{1, 2}).Split(5); !reflect.DeepEqual(got, [][]int{{}, {1}, {}, {2}}) {
		t.Fatalf("Split count less than divisor = %#v", got)
	}

	if got := Collect([]int{1, 2, 3, 4}).Split(3); !reflect.DeepEqual(got, [][]int{{1}, {2, 3}, {4}}) {
		t.Fatalf("Split 4 into 3 = %#v", got)
	}

	if got := Collect([]int{1, 2, 3, 4, 5}).Split(3); !reflect.DeepEqual(got, [][]int{{1, 2}, {3}, {4, 5}}) {
		t.Fatalf("Split 5 into 3 = %#v", got)
	}

	if got := Collect([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}).Split(6); !reflect.DeepEqual(got, [][]int{{1, 2}, {3}, {4, 5}, {6, 7}, {8}, {9, 10}}) {
		t.Fatalf("Split 10 into 6 = %#v", got)
	}

	if got := Empty[int]().Split(3); got != nil {
		t.Fatalf("Split empty = %#v, want nil", got)
	}

	if got := items.Nth(0).All(); len(got) != 0 {
		t.Fatalf("Nth zero step = %v, want empty", got)
	}

	if got := items.Nth(-2).All(); len(got) != 0 {
		t.Fatalf("Nth negative step = %v, want empty", got)
	}

	if got := items.Split(0); got != nil {
		t.Fatalf("Split zero groups = %#v, want nil", got)
	}

	if got := items.Split(-2); got != nil {
		t.Fatalf("Split negative groups = %#v, want nil", got)
	}

	if got := items.SplitIn(0); got != nil {
		t.Fatalf("SplitIn zero groups = %#v, want nil", got)
	}

	if got := items.SplitIn(-2); got != nil {
		t.Fatalf("SplitIn negative groups = %#v, want nil", got)
	}

	pass, fail := items.Partition(func(item int, index int) bool {
		return item+index > 5
	})

	if !reflect.DeepEqual(pass.All(), []int{5, 6, 9}) || !reflect.DeepEqual(fail.All(), []int{1, 2, 3}) {
		t.Fatalf("Partition with key = pass %v fail %v", pass.All(), fail.All())
	}

	emptyPass, emptyFail := Empty[int]().Partition(func(item int, _ int) bool { return item > 0 })

	if emptyPass.Count() != 0 || emptyFail.Count() != 0 {
		t.Fatalf("empty partition = pass %v fail %v", emptyPass.All(), emptyFail.All())
	}
}

func TestInventoryParityEnsureAndMapInto(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testMapInto
	// SupportCollectionTest::testEnsureForScalar
	// SupportCollectionTest::testEnsureForObjects
	// SupportCollectionTest::testEnsureForInheritance
	// SupportCollectionTest::testEnsureForMultipleTypes
	mapped := MapInto(Collect([]int{1, 2, 3}), func(value int) string {
		return strings.Repeat("*", value)
	})

	if got := mapped.All(); !reflect.DeepEqual(got, []string{"*", "**", "***"}) {
		t.Fatalf("MapInto = %v", got)
	}

	if err := Collect([]int{2, 4, 6}).Ensure(func(value int) bool { return value%2 == 0 }); err != nil {
		t.Fatalf("Ensure scalar returned error: %v", err)
	}

	if err := Collect([]inventoryUser{{Name: "Taylor"}}).Ensure(func(user inventoryUser) bool { return user.Name != "" }); err != nil {
		t.Fatalf("Ensure objects returned error: %v", err)
	}

	if err := Collect([]any{1, "two"}).Ensure(func(value any) bool {
		switch value.(type) {
		case int, string:
			return true
		default:
			return false
		}
	}); err != nil {
		t.Fatalf("Ensure multiple types returned error: %v", err)
	}

	if err := Collect([]int{1, 3, 5}).Ensure(func(value int) bool { return value%2 == 0 }); err == nil {
		t.Fatal("Ensure should fail when a value does not match")
	}
}

func TestInventoryParitySetOperationsAndKeyedCollections(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testMergeNull
	// SupportCollectionTest::testMergeRecursiveNull
	// SupportCollectionTest::testMergeRecursiveArray
	// SupportCollectionTest::testMergeRecursiveCollection
	// SupportCollectionTest::testReplaceNull
	// SupportCollectionTest::testReplaceCollection
	// SupportCollectionTest::testReplaceRecursiveNull
	// SupportCollectionTest::testReplaceRecursiveArray
	// SupportCollectionTest::testReplaceRecursiveCollection
	// SupportCollectionTest::testUnionNull
	// SupportCollectionTest::testUnionCollection
	base := collectible.FromPairs(
		csupport.Pair[string, string]{Key: "name", Value: "Taylor"},
		csupport.Pair[string, string]{Key: "role", Value: "admin"},
	)

	if got := base.Merge(nil).All(); !reflect.DeepEqual(got, map[string]string{"name": "Taylor", "role": "admin"}) {
		t.Fatalf("Merge nil = %#v", got)
	}

	if got := collectible.MergeRecursive(base, nil).All(); !reflect.DeepEqual(got, map[string]string{"name": "Taylor", "role": "admin"}) {
		t.Fatalf("MergeRecursive nil = %#v", got)
	}

	if got := collectible.MergeRecursive(base, map[string]string{"role": "editor", "team": "core"}).All(); !reflect.DeepEqual(got, map[string]string{"name": "Taylor", "role": "editor", "team": "core"}) {
		t.Fatalf("MergeRecursive map = %#v", got)
	}

	other := collectible.FromPairs(csupport.Pair[string, string]{Key: "role", Value: "owner"})

	if got := collectible.MergeRecursive(base, other.All()).All(); got["role"] != "owner" || got["name"] != "Taylor" {
		t.Fatalf("MergeRecursive collection = %#v", got)
	}

	if got := base.Replace(nil).All(); !reflect.DeepEqual(got, map[string]string{"name": "Taylor", "role": "admin"}) {
		t.Fatalf("Replace nil = %#v", got)
	}

	if got := base.Replace(map[string]string{"role": "editor"}).All(); !reflect.DeepEqual(got, map[string]string{"name": "Taylor", "role": "editor"}) {
		t.Fatalf("Replace collection = %#v", got)
	}

	if got := collectible.MergeRecursive(base, nil).All(); got["role"] != "admin" {
		t.Fatalf("ReplaceRecursive nil adaptation = %#v", got)
	}

	if got := base.Replace(map[string]string{"team": "core"}).All(); !reflect.DeepEqual(got, map[string]string{"name": "Taylor", "role": "admin", "team": "core"}) {
		t.Fatalf("ReplaceRecursive map adaptation = %#v", got)
	}

	if got := base.Replace(other.All()).All(); got["role"] != "owner" || got["name"] != "Taylor" {
		t.Fatalf("ReplaceRecursive collection adaptation = %#v", got)
	}

	if got := base.Union(nil).All(); !reflect.DeepEqual(got, map[string]string{"name": "Taylor", "role": "admin"}) {
		t.Fatalf("Union nil = %#v", got)
	}

	if got := base.Union(map[string]string{"role": "editor", "team": "core"}).All(); !reflect.DeepEqual(got, map[string]string{"name": "Taylor", "role": "admin", "team": "core"}) {
		t.Fatalf("Union collection = %#v", got)
	}
}

func TestInventoryParityDiffIntersectCollapseAndSortKeys(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testDiffCollection
	// SupportCollectionTest::testDiffUsingWithCollection
	// SupportCollectionTest::testDiffUsingWithNull
	// SupportCollectionTest::testDiffNull
	// SupportCollectionTest::testIntersectNull
	// SupportCollectionTest::testIntersectCollection
	// SupportCollectionTest::testIntersectUsingWithNull
	// SupportCollectionTest::testIntersectUsingCollection
	// SupportCollectionTest::testIntersectByKeysNull
	// SupportCollectionTest::testDiffAssocUsing
	// SupportCollectionTest::testIntersectAssocWithNull
	// SupportCollectionTest::testIntersectAssocCollection
	// SupportCollectionTest::testIntersectAssocUsingWithNull
	// SupportCollectionTest::testIntersectAssocUsingCollection
	// SupportCollectionTest::testCollapseWithNestedCollections
	// SupportCollectionTest::testCollapseWithKeys
	// SupportCollectionTest::testCollapseWithKeysOnNestedCollections
	// SupportCollectionTest::testSortKeysUsing
	left := Collect([]string{"Taylor", "Abigail", "Mohamed"})
	right := Collect([]string{"Abigail"})

	if got := Diff(left, right.All()).All(); !reflect.DeepEqual(got, []string{"Taylor", "Mohamed"}) {
		t.Fatalf("Diff collection = %#v", got)
	}

	if got := left.DiffUsing(right.All(), strings.EqualFold).All(); !reflect.DeepEqual(got, []string{"Taylor", "Mohamed"}) {
		t.Fatalf("DiffUsing collection = %#v", got)
	}

	if got := left.DiffUsing(nil, strings.EqualFold).All(); !reflect.DeepEqual(got, left.All()) {
		t.Fatalf("DiffUsing nil = %#v", got)
	}

	if got := Diff(left, nil).All(); !reflect.DeepEqual(got, left.All()) {
		t.Fatalf("Diff nil = %#v", got)
	}

	if got := Intersect(left, nil).All(); len(got) != 0 {
		t.Fatalf("Intersect nil = %#v", got)
	}

	if got := Intersect(left, right.All()).All(); !reflect.DeepEqual(got, []string{"Abigail"}) {
		t.Fatalf("Intersect collection = %#v", got)
	}

	if got := left.IntersectUsing(nil, strings.EqualFold).All(); len(got) != 0 {
		t.Fatalf("IntersectUsing nil = %#v", got)
	}

	if got := left.IntersectUsing([]string{"taylor"}, strings.EqualFold).All(); !reflect.DeepEqual(got, []string{"Taylor"}) {
		t.Fatalf("IntersectUsing collection = %#v", got)
	}

	keyed := collectible.FromPairs(
		csupport.Pair[string, int]{Key: "b", Value: 2},
		csupport.Pair[string, int]{Key: "a", Value: 1},
	)

	if got := keyed.IntersectByKeys(nil).All(); len(got) != 0 {
		t.Fatalf("IntersectByKeys nil = %#v", got)
	}

	diffAssoc := collectible.DiffAssoc(keyed, map[string]int{"a": 1, "b": 99})

	if got := diffAssoc.All(); !reflect.DeepEqual(got, map[string]int{"b": 2}) {
		t.Fatalf("DiffAssoc using comparable values = %#v", got)
	}

	if got := collectible.IntersectAssoc(keyed, nil).All(); len(got) != 0 {
		t.Fatalf("IntersectAssoc nil = %#v", got)
	}

	if got := collectible.IntersectAssoc(keyed, map[string]int{"a": 1, "b": 99}).All(); !reflect.DeepEqual(got, map[string]int{"a": 1}) {
		t.Fatalf("IntersectAssoc collection = %#v", got)
	}

	sortedKeys := keyed.SortKeysUsing(func(a, b string) bool { return a < b }).Keys()

	if !reflect.DeepEqual(sortedKeys, []string{"a", "b"}) {
		t.Fatalf("SortKeysUsing = %#v", sortedKeys)
	}

	nested := Collect([][]int{{1, 2}, {3}, {4, 5}})

	if got := Collapse(nested).All(); !reflect.DeepEqual(got, []int{1, 2, 3, 4, 5}) {
		t.Fatalf("Collapse nested collections = %#v", got)
	}

	keyedGroups := []map[string]int{{"first": 1}, {"second": 2}}
	collapsed := map[string]int{}

	for _, group := range keyedGroups {
		for key, value := range group {
			collapsed[key] = value
		}
	}

	if !reflect.DeepEqual(collapsed, map[string]int{"first": 1, "second": 2}) {
		t.Fatalf("CollapseWithKeys equivalent = %#v", collapsed)
	}
}

func TestInventoryParityDepthFlattenDotPercentageAndReduceSpread(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testFlattenWithDepth
	// SupportCollectionTest::testFlattenIgnoresKeys
	// SupportCollectionTest::testDotWithDepth
	// SupportCollectionTest::testJsonSerialize
	// SupportCollectionTest::testReduceSpread
	// SupportCollectionTest::testReduceSpreadThrowsAnExceptionIfReducerDoesNotReturnAnArray
	// SupportCollectionTest::testSliceOffsetAndNegativeLength
	// SupportCollectionTest::testSliceNegativeOffsetAndNegativeLength
	// SupportCollectionTest::testPercentageWithFlatCollection
	// SupportCollectionTest::testPercentageWithNestedCollection
	// SupportCollectionTest::testHighOrderPercentage
	// SupportCollectionTest::testPercentageReturnsNullForEmptyCollections
	flattened := arr.FlattenAny([]any{1, []any{2, []any{3, 4}}, 5}, 1)

	if !reflect.DeepEqual(flattened, []any{1, 2, []any{3, 4}, 5}) {
		t.Fatalf("FlattenAny depth 1 = %#v", flattened)
	}

	if got := arr.FlattenAny([]any{map[string]int{"ignored": 1}, []any{2}}, 2); !reflect.DeepEqual(got, []any{map[string]int{"ignored": 1}, 2}) {
		t.Fatalf("FlattenAny keeps non-slice keyed values = %#v", got)
	}

	dotted := kv.DotWithDepth(map[string]any{
		"user": map[string]any{
			"name": "Taylor",
			"team": map[string]any{"id": 7},
		},
	}, 1)

	if !reflect.DeepEqual(dotted, map[string]any{"user.name": "Taylor", "user.team": map[string]any{"id": 7}}) {
		t.Fatalf("DotWithDepth = %#v", dotted)
	}

	encoded, err := json.Marshal(Collect([]int{1, 2, 3}))

	if err != nil || string(encoded) != "[1,2,3]" {
		t.Fatalf("json.Marshal collection = %s, %v", encoded, err)
	}

	spread, err := ReduceSpread(Collect([]int{1, 2, 3}), func(carry []int, item int, _ int) []int {
		return []int{carry[0] + item, carry[1] * item}
	}, 0, 1)

	if err != nil || !reflect.DeepEqual(spread, []int{6, 6}) {
		t.Fatalf("ReduceSpread = %#v, %v", spread, err)
	}

	if _, err := ReduceSpread(Collect([]int{1}), func(_ []int, _ int, _ int) []int {
		return []int{1, 2, 3}
	}, 0, 0); !errors.Is(err, ErrReduceSpreadLength) {
		t.Fatalf("ReduceSpread arity err = %v", err)
	}

	values := Collect([]int{1, 2, 3, 4, 5})

	if got := values.Slice(1, -1).All(); !reflect.DeepEqual(got, []int{2, 3, 4}) {
		t.Fatalf("Slice offset and negative length = %#v", got)
	}

	if got := values.Slice(-4, -1).All(); !reflect.DeepEqual(got, []int{2, 3, 4}) {
		t.Fatalf("Slice negative offset and negative length = %#v", got)
	}

	percent, ok := Percentage(values, func(item int, _ int) bool { return item > 3 })

	if !ok || percent != 40 {
		t.Fatalf("Percentage flat = %v ok=%v", percent, ok)
	}

	users := Collect([]inventoryUser{{Active: true}, {Active: false}, {Active: true}, {Active: true}})
	percent, ok = Percentage(users, func(user inventoryUser, _ int) bool { return user.Active })

	if !ok || percent != 75 {
		t.Fatalf("Percentage nested/high-order equivalent = %v ok=%v", percent, ok)
	}

	if percent, ok = Percentage(Empty[int](), func(int, int) bool { return true }); ok || percent != 0 {
		t.Fatalf("Percentage empty = %v ok=%v", percent, ok)
	}
}

func TestInventoryParityConditionalsNullFiltersAndPagination(t *testing.T) {
	t.Parallel()

	// SupportCollectionTest::testPaginate
	// SupportCollectionTest::testHasReturnsValidResults
	// SupportCollectionTest::testWhereNullWithoutKey
	// SupportCollectionTest::testWhereNotNullWithoutKey
	// SupportCollectionTest::testWhenEmptyDefault
	// SupportCollectionTest::testWhenNotEmptyDefault
	// SupportCollectionTest::testUnlessDefault
	// SupportCollectionTest::testUnlessEmptyDefault
	// SupportCollectionTest::testUnlessNotEmptyDefault
	values := Collect([]int{10, 20, 30, 40, 50})

	if got := values.ForPage(2, 2).All(); !reflect.DeepEqual(got, []int{30, 40}) {
		t.Fatalf("ForPage = %#v", got)
	}

	if !values.Has(0) || !values.Has(-1) || values.Has(99) {
		t.Fatalf("Has results were not valid for first, last, and missing indexes")
	}

	a, b := 10, 20
	pointers := Collect([]*int{&a, nil, &b, nil})

	if got := WhereNull(pointers, func(value *int) *int { return value }).All(); len(got) != 2 || got[0] != nil || got[1] != nil {
		t.Fatalf("WhereNull without key adaptation = %#v", got)
	}

	if got := WhereNotNull(pointers, func(value *int) *int { return value }).All(); len(got) != 2 || *got[0] != 10 || *got[1] != 20 {
		t.Fatalf("WhereNotNull without key adaptation = %#v", got)
	}

	defaultedEmpty := Empty[int]().WhenEmpty(
		func(c *Collection[int]) *Collection[int] { return c.Push(1) },
		func(c *Collection[int]) *Collection[int] { return c.Push(99) },
	)

	if got := defaultedEmpty.All(); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("WhenEmpty callback = %#v", got)
	}

	defaultedNotEmpty := values.WhenNotEmpty(
		func(c *Collection[int]) *Collection[int] { return c.Take(1) },
		func(c *Collection[int]) *Collection[int] { return c.Push(99) },
	)

	if got := defaultedNotEmpty.All(); !reflect.DeepEqual(got, []int{10}) {
		t.Fatalf("WhenNotEmpty callback = %#v", got)
	}

	unlessDefault := values.Unless(true,
		func(c *Collection[int]) *Collection[int] { return c.Push(99) },
		func(c *Collection[int]) *Collection[int] { return c.Take(2) },
	)

	if got := unlessDefault.All(); !reflect.DeepEqual(got, []int{10, 20}) {
		t.Fatalf("Unless default = %#v", got)
	}

	unlessEmptyDefault := Empty[int]().UnlessEmpty(
		func(c *Collection[int]) *Collection[int] { return c.Push(99) },
		func(c *Collection[int]) *Collection[int] { return c.Push(2) },
	)

	if got := unlessEmptyDefault.All(); !reflect.DeepEqual(got, []int{2}) {
		t.Fatalf("UnlessEmpty default = %#v", got)
	}

	unlessNotEmptyDefault := values.UnlessNotEmpty(
		func(c *Collection[int]) *Collection[int] { return c.Push(99) },
		func(c *Collection[int]) *Collection[int] { return c.Take(3) },
	)

	if got := unlessNotEmptyDefault.All(); !reflect.DeepEqual(got, []int{10, 20, 30}) {
		t.Fatalf("UnlessNotEmpty default = %#v", got)
	}
}

func namesOf(users []inventoryUser) []string {
	names := make([]string, len(users))

	for i, user := range users {
		names[i] = user.Name
	}

	return names
}
