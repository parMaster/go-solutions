package sandbox

import (
	"sort"
	"testing"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/stretchr/testify/assert"
)

// Sane ways to sort a slice of structs:
// 1. Implement sort.Interface and use sort.Sort
// 2. Use sort.Slice
// 3. Use slices.SortFunc

// sort.Interface must be implemented by the type:
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

// Record is a struct with a StartTime field that will be used to sort Records
type Record struct {
	Id        string    `json:"id"`              // primary key for Record
	MeetingId string    `json:"meeting_id"`      // foreign key to Meeting.UUID
	StartTime time.Time `json:"recording_start"` // DateTime in RFC3339
	DateTime  string    `json:"date_time"`
}

// Records is a slice of Record, it implements sort.Interface
type Records []Record

func (r Records) Len() int {
	return len(r)
}

func (r Records) Less(i, j int) bool {
	return r[i].StartTime.Before(r[j].StartTime)
}

func (r Records) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// NewRecords returns a Records slice
func NewRecords() Records {
	return Records{}
}

func Test_SortInterface(t *testing.T) {

	r := NewRecords()

	r = append(r, Record{Id: "1", MeetingId: "1", StartTime: time.Now().AddDate(0, 0, -1)})
	r = append(r, Record{Id: "2", MeetingId: "1", StartTime: time.Now().AddDate(0, 0, -2)})
	r = append(r, Record{Id: "3", MeetingId: "1", StartTime: time.Now().AddDate(0, 0, -3)})

	log.Printf("Records: %v", r)

	sort.Sort(r)

	log.Printf("Records: %v", r) // sorted by StartTime ascending

	sort.Sort(sort.Reverse(r)) // functional options pattern, btw

	log.Printf("Records: %v", r) // sorted by StartTime descending

	r1 := make(Records, len(r))
	copy(r1, r)
	sort.Sort(sort.Reverse(r))
	assert.Equal(t, r1, r) // r hasn't changed because already sorted descending
}

func Test_SortSlice(t *testing.T) {

	r := NewRecords()

	r = append(r, Record{Id: "1", MeetingId: "1", StartTime: time.Now().AddDate(0, 0, -1)})
	r = append(r, Record{Id: "2", MeetingId: "1", StartTime: time.Now().AddDate(0, 0, -2)})
	r = append(r, Record{Id: "3", MeetingId: "1", StartTime: time.Now().AddDate(0, 0, -3)})

	log.Printf("Before Records: %v", r)

	r1 := make(Records, len(r))
	copy(r1, r)

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].StartTime.Before(r[j].StartTime)
	})

	log.Printf("After Records: %v", r) // sorted by StartTime ascending

	sort.SliceStable(r, func(i, j int) bool {
		return r[i].StartTime.After(r[j].StartTime) // reverse sort
	})

	assert.Equal(t, r1, r) // r was sorted then reversed, so it's like the original
}
