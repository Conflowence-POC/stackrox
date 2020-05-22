// Code generated by genny. DO NOT EDIT.
// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/mauricelam/genny

package set

import (
	"fmt"
	"sort"
	"strings"
)

// If you want to add a set for your custom type, simply add another go generate line along with the
// existing ones. If you're creating a set for a primitive type, you can follow the example of "string"
// and create the generated file in this package.
// For non-primitive sets, please make the generated code files go outside this package.
// Sometimes, you might need to create it in the same package where it is defined to avoid import cycles.
// The permission set is an example of how to do that.
// You can also specify the -imp command to specify additional imports in your generated file, if required.

// int represents a generic type that we want to have a set of.

// IntSet will get translated to generic sets.
type IntSet map[int]struct{}

// Add adds an element of type int.
func (k *IntSet) Add(i int) bool {
	if *k == nil {
		*k = make(map[int]struct{})
	}

	oldLen := len(*k)
	(*k)[i] = struct{}{}
	return len(*k) > oldLen
}

// AddMatching is a utility function that adds all the elements that match the given function to the set.
func (k *IntSet) AddMatching(matchFunc func(int) bool, elems ...int) bool {
	oldLen := len(*k)
	for _, elem := range elems {
		if !matchFunc(elem) {
			continue
		}
		if *k == nil {
			*k = make(map[int]struct{})
		}
		(*k)[elem] = struct{}{}
	}
	return len(*k) > oldLen
}

// AddAll adds all elements of type int. The return value is true if any new element
// was added.
func (k *IntSet) AddAll(is ...int) bool {
	if len(is) == 0 {
		return false
	}
	if *k == nil {
		*k = make(map[int]struct{})
	}

	oldLen := len(*k)
	for _, i := range is {
		(*k)[i] = struct{}{}
	}
	return len(*k) > oldLen
}

// Remove removes an element of type int.
func (k *IntSet) Remove(i int) bool {
	if len(*k) == 0 {
		return false
	}

	oldLen := len(*k)
	delete(*k, i)
	return len(*k) < oldLen
}

// RemoveAll removes the given elements.
func (k *IntSet) RemoveAll(is ...int) bool {
	if len(*k) == 0 {
		return false
	}

	oldLen := len(*k)
	for _, i := range is {
		delete(*k, i)
	}
	return len(*k) < oldLen
}

// RemoveMatching removes all elements that match a given predicate.
func (k *IntSet) RemoveMatching(pred func(int) bool) bool {
	if len(*k) == 0 {
		return false
	}

	oldLen := len(*k)
	for elem := range *k {
		if pred(elem) {
			delete(*k, elem)
		}
	}
	return len(*k) < oldLen
}

// Contains returns whether the set contains an element of type int.
func (k IntSet) Contains(i int) bool {
	_, ok := k[i]
	return ok
}

// Cardinality returns the number of elements in the set.
func (k IntSet) Cardinality() int {
	return len(k)
}

// IsEmpty returns whether the underlying set is empty (includes uninitialized).
func (k IntSet) IsEmpty() bool {
	return len(k) == 0
}

// Clone returns a copy of this set.
func (k IntSet) Clone() IntSet {
	if k == nil {
		return nil
	}
	cloned := make(map[int]struct{}, len(k))
	for elem := range k {
		cloned[elem] = struct{}{}
	}
	return cloned
}

// Difference returns a new set with all elements of k not in other.
func (k IntSet) Difference(other IntSet) IntSet {
	if len(k) == 0 || len(other) == 0 {
		return k.Clone()
	}

	retained := make(map[int]struct{}, len(k))
	for elem := range k {
		if !other.Contains(elem) {
			retained[elem] = struct{}{}
		}
	}
	return retained
}

// Helper function for intersections.
func (k IntSet) getSmallerLargerAndMaxIntLen(other IntSet) (smaller IntSet, larger IntSet, maxIntLen int) {
	maxIntLen = len(k)
	smaller, larger = k, other
	if l := len(other); l < maxIntLen {
		maxIntLen = l
		smaller, larger = larger, smaller
	}
	return smaller, larger, maxIntLen
}

// Intersects returns whether the set has a non-empty intersection with the other set.
func (k IntSet) Intersects(other IntSet) bool {
	smaller, larger, maxIntLen := k.getSmallerLargerAndMaxIntLen(other)
	if maxIntLen == 0 {
		return false
	}
	for elem := range smaller {
		if _, ok := larger[elem]; ok {
			return true
		}
	}
	return false
}

// Intersect returns a new set with the intersection of the members of both sets.
func (k IntSet) Intersect(other IntSet) IntSet {
	smaller, larger, maxIntLen := k.getSmallerLargerAndMaxIntLen(other)
	if maxIntLen == 0 {
		return nil
	}

	retained := make(map[int]struct{}, maxIntLen)
	for elem := range smaller {
		if _, ok := larger[elem]; ok {
			retained[elem] = struct{}{}
		}
	}
	return retained
}

// Union returns a new set with the union of the members of both sets.
func (k IntSet) Union(other IntSet) IntSet {
	if len(k) == 0 {
		return other.Clone()
	} else if len(other) == 0 {
		return k.Clone()
	}

	underlying := make(map[int]struct{}, len(k)+len(other))
	for elem := range k {
		underlying[elem] = struct{}{}
	}
	for elem := range other {
		underlying[elem] = struct{}{}
	}
	return underlying
}

// Equal returns a bool if the sets are equal
func (k IntSet) Equal(other IntSet) bool {
	thisL, otherL := len(k), len(other)
	if thisL == 0 && otherL == 0 {
		return true
	}
	if thisL != otherL {
		return false
	}
	for elem := range k {
		if _, ok := other[elem]; !ok {
			return false
		}
	}
	return true
}

// AsSlice returns a slice of the elements in the set. The order is unspecified.
func (k IntSet) AsSlice() []int {
	if len(k) == 0 {
		return nil
	}
	elems := make([]int, 0, len(k))
	for elem := range k {
		elems = append(elems, elem)
	}
	return elems
}

// GetArbitraryElem returns an arbitrary element from the set.
// This can be useful if, for example, you know the set has exactly one
// element, and you want to pull it out.
// If the set is empty, the zero value is returned.
func (k IntSet) GetArbitraryElem() (arbitraryElem int) {
	for elem := range k {
		arbitraryElem = elem
		break
	}
	return arbitraryElem
}

// AsSortedSlice returns a slice of the elements in the set, sorted using the passed less function.
func (k IntSet) AsSortedSlice(less func(i, j int) bool) []int {
	slice := k.AsSlice()
	if len(slice) < 2 {
		return slice
	}
	// Since we're generating the code, we might as well use sort.Sort
	// and avoid paying the reflection penalty of sort.Slice.
	sortable := &sortableIntSlice{slice: slice, less: less}
	sort.Sort(sortable)
	return sortable.slice
}

// Clear empties the set
func (k *IntSet) Clear() {
	*k = nil
}

// Freeze returns a new, frozen version of the set.
func (k IntSet) Freeze() FrozenIntSet {
	return NewFrozenIntSetFromMap(k)
}

// ElementsString returns a string representation of all elements, with individual element strings separated by `sep`.
// The string representation of an individual element is obtained via `fmt.Fprint`.
func (k IntSet) ElementsString(sep string) string {
	if len(k) == 0 {
		return ""
	}
	var sb strings.Builder
	first := true
	for elem := range k {
		if !first {
			sb.WriteString(sep)
		}
		fmt.Fprint(&sb, elem)
		first = false
	}
	return sb.String()
}

// NewIntSet returns a new thread unsafe set with the given key type.
func NewIntSet(initial ...int) IntSet {
	underlying := make(map[int]struct{}, len(initial))
	for _, elem := range initial {
		underlying[elem] = struct{}{}
	}
	return underlying
}

type sortableIntSlice struct {
	slice []int
	less  func(i, j int) bool
}

func (s *sortableIntSlice) Len() int {
	return len(s.slice)
}

func (s *sortableIntSlice) Less(i, j int) bool {
	return s.less(s.slice[i], s.slice[j])
}

func (s *sortableIntSlice) Swap(i, j int) {
	s.slice[j], s.slice[i] = s.slice[i], s.slice[j]
}

// A FrozenIntSet is a frozen set of int elements, which
// cannot be modified after creation. This allows users to use it as if it were
// a "const" data structure, and also makes it slightly more optimal since
// we don't have to lock accesses to it.
type FrozenIntSet struct {
	underlying map[int]struct{}
}

// NewFrozenIntSetFromMap returns a new frozen set from the set-style map.
func NewFrozenIntSetFromMap(m map[int]struct{}) FrozenIntSet {
	if len(m) == 0 {
		return FrozenIntSet{}
	}
	underlying := make(map[int]struct{}, len(m))
	for elem := range m {
		underlying[elem] = struct{}{}
	}
	return FrozenIntSet{
		underlying: underlying,
	}
}

// NewFrozenIntSet returns a new frozen set with the provided elements.
func NewFrozenIntSet(elements ...int) FrozenIntSet {
	underlying := make(map[int]struct{}, len(elements))
	for _, elem := range elements {
		underlying[elem] = struct{}{}
	}
	return FrozenIntSet{
		underlying: underlying,
	}
}

// Contains returns whether the set contains the element.
func (k FrozenIntSet) Contains(elem int) bool {
	_, ok := k.underlying[elem]
	return ok
}

// Cardinality returns the cardinality of the set.
func (k FrozenIntSet) Cardinality() int {
	return len(k.underlying)
}

// IsEmpty returns whether the underlying set is empty (includes uninitialized).
func (k FrozenIntSet) IsEmpty() bool {
	return len(k.underlying) == 0
}

// AsSlice returns the elements of the set. The order is unspecified.
func (k FrozenIntSet) AsSlice() []int {
	if len(k.underlying) == 0 {
		return nil
	}
	slice := make([]int, 0, len(k.underlying))
	for elem := range k.underlying {
		slice = append(slice, elem)
	}
	return slice
}

// AsSortedSlice returns the elements of the set as a sorted slice.
func (k FrozenIntSet) AsSortedSlice(less func(i, j int) bool) []int {
	slice := k.AsSlice()
	if len(slice) < 2 {
		return slice
	}
	// Since we're generating the code, we might as well use sort.Sort
	// and avoid paying the reflection penalty of sort.Slice.
	sortable := &sortableIntSlice{slice: slice, less: less}
	sort.Sort(sortable)
	return sortable.slice
}

// ElementsString returns a string representation of all elements, with individual element strings separated by `sep`.
// The string representation of an individual element is obtained via `fmt.Fprint`.
func (k FrozenIntSet) ElementsString(sep string) string {
	if len(k.underlying) == 0 {
		return ""
	}
	var sb strings.Builder
	first := true
	for elem := range k.underlying {
		if !first {
			sb.WriteString(sep)
		}
		fmt.Fprint(&sb, elem)
		first = false
	}
	return sb.String()
}

// The following functions make use of casting `k.underlying` into a mutable Set. This is safe, since we never leak
// references to these objects, and only invoke mutable set methods that are guaranteed to return a new copy.

// Union returns a frozen set that represents the union between this and other.
func (k FrozenIntSet) Union(other FrozenIntSet) FrozenIntSet {
	if len(k.underlying) == 0 {
		return other
	}
	if len(other.underlying) == 0 {
		return k
	}
	return FrozenIntSet{
		underlying: IntSet(k.underlying).Union(other.underlying),
	}
}

// Intersect returns a frozen set that represents the intersection between this and other.
func (k FrozenIntSet) Intersect(other FrozenIntSet) FrozenIntSet {
	return FrozenIntSet{
		underlying: IntSet(k.underlying).Intersect(other.underlying),
	}
}

// Difference returns a frozen set that represents the set difference between this and other.
func (k FrozenIntSet) Difference(other FrozenIntSet) FrozenIntSet {
	return FrozenIntSet{
		underlying: IntSet(k.underlying).Difference(other.underlying),
	}
}

// Unfreeze returns a mutable set with the same contents as this frozen set. This set will not be affected by any
// subsequent modifications to the returned set.
func (k FrozenIntSet) Unfreeze() IntSet {
	return IntSet(k.underlying).Clone()
}
