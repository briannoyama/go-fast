package fast

import (
	"testing"

	"go-fast/assert"
)

var buffer CBuffer[int]

func setupBuffer() {
	buffer = MakeCBuffer[int](2)
	cPageRefs = make([]int, 12)
	for i := 0; i < len(cPageRefs); i++ {
		buffer.Add(&cPageRefs[i], i)
	}
	buffer.MarkClean()
}

func TestBufferAdd(t *testing.T) {
	setupBuffer()
	var ref int
	buffer.Add(&ref, 12)
	assert.Equals(t, ref, 12)
	assert.Equals(t, buffer.Pages()[3].Data(), []int{12})
}

func TestBufferGet(t *testing.T) {
	setupBuffer()

	assert.Equals(t, buffer.Get(cPageRefs[2]), 2)
	assert.Equals(t, buffer.Get(cPageRefs[4]), 4)
	assert.Equals(t, buffer.Get(cPageRefs[9]), 9)
}

func TestBufferModify(t *testing.T) {
	setupBuffer()

	pages := buffer.Pages()
	modified := []int{}
	buffer.Modify(cPageRefs[2], func(i *int) { modified = append(modified, *i) })
	buffer.Modify(cPageRefs[4], func(i *int) { modified = append(modified, *i) })
	buffer.Modify(cPageRefs[7], func(i *int) { modified = append(modified, *i) })
	buffer.Modify(cPageRefs[9], func(i *int) { modified = append(modified, *i) })
	assert.Equals(t, modified, []int{2, 4, 7, 9})

	assert.Equals(t, pages[0].CleanLen(), 2)
	assert.Equals(t, pages[1].CleanLen(), 0)
	assert.Equals(t, pages[2].CleanLen(), 1)

	buffer.MarkClean()
	modified = []int{}
	buffer.Modify(cPageRefs[2], func(i *int) { modified = append(modified, *i) })
	buffer.Modify(cPageRefs[4], func(i *int) { modified = append(modified, *i) })
	buffer.Modify(cPageRefs[7], func(i *int) { modified = append(modified, *i) })
	buffer.Modify(cPageRefs[9], func(i *int) { modified = append(modified, *i) })
	assert.Equals(t, modified, []int{2, 4, 7, 9})

	// Show that modified references are moved to the end.
	assert.Equals(t, pages[0].CleanLen(), 3)
	assert.Equals(t, pages[1].CleanLen(), 2)
	assert.Equals(t, pages[2].CleanLen(), 3)
}

func TestBufferRemove(t *testing.T) {
	setupBuffer()

	assert.Equals(t, buffer.Remove(cPageRefs[2]), 2)
	assert.Equals(t, buffer.Remove(cPageRefs[4]), 4)
	assert.Equals(t, buffer.Remove(cPageRefs[7]), 7)
	assert.Equals(t, buffer.Remove(cPageRefs[9]), 9)

	pages := buffer.Pages()
	assert.Contains(t, pages[0].Data(), 0, 1, 3, 11)
	assert.Contains(t, pages[1].Data(), 6, 5, 8, 10)
	assert.Equals(t, len(pages), 2)
	assert.Equals(t, buffer.Len(), 8)

	assert.Equals(t, pages[0].CleanLen(), 2)
	assert.Equals(t, pages[1].CleanLen(), 0)

	buffer.MarkClean()
	assert.Equals(t, buffer.Remove(cPageRefs[3]), 3)
	assert.Equals(t, buffer.Remove(cPageRefs[5]), 5)
	assert.Equals(t, buffer.Remove(cPageRefs[6]), 6)
	assert.Equals(t, buffer.Remove(cPageRefs[8]), 8)
	assert.Equals(t, buffer.Remove(cPageRefs[10]), 10)

	pages = buffer.Pages()
	assert.Contains(t, pages[0].Data(), 0, 1, 11)
	assert.Equals(t, len(pages), 1)
	assert.Equals(t, buffer.Len(), 3)

	buffer.MarkClean()
	assert.Equals(t, buffer.Remove(cPageRefs[0]), 0)
	assert.Equals(t, buffer.Remove(cPageRefs[1]), 1)
	assert.Equals(t, buffer.Remove(cPageRefs[11]), 11)

	pages = buffer.Pages()
	assert.Equals(t, pages[0].Data(), []int{})
	assert.Equals(t, len(pages), 1)
	assert.Equals(t, buffer.Len(), 0)
}
