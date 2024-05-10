package sandbox

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Problem:
// JSON encoding and decoding for type alias like this:
type Date time.Time

// looks like this:

func (d Date) MarshalJSON() ([]byte, error) {
	return time.Time(d).MarshalJSON()
}
func (d *Date) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := t.UnmarshalJSON(data); err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

// But manipulating Date as time.Time is not convenient.

func Test_DateManipulation(t *testing.T) {
	d := Date(time.Now())

	// d.AddDate(1, 0, 0) // Error: Date does not have AddDate method.
	// d.Format("2006-01-02") // Error: Date does not have Format method.
	// d.Year() // Error: Date does not have Year method.
	// etc.

	// We have to convert Date to time.Time first.
	tt := time.Time(d)
	tomorrow := tt.AddDate(1, 0, 0)

	assert.Equal(t,
		tomorrow.Format("2006-01-02"),
		time.Now().AddDate(1, 0, 0).Format("2006-01-02"))

}

// We want to use Date as time.Time directly.

// Solution:
// Use a composite type to wrap time.Time and implement the MarshalJSON and UnmarshalJSON methods.

type Date2 struct {
	time.Time
}

func (d Date2) MarshalJSON() ([]byte, error) {
	return d.Time.MarshalJSON()
}
func (d *Date2) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := t.UnmarshalJSON(data); err != nil {
		return err
	}
	d.Time = t
	return nil
}

// Now we can use Date2 as time.Time directly.
func Test_Date2Manipulation(t *testing.T) {
	d := Date2{Time: time.Now()}

	tomorrow := d.AddDate(1, 0, 0)

	assert.Equal(t,
		tomorrow.Format("2006-01-02"),
		time.Now().AddDate(1, 0, 0).Format("2006-01-02"))

}
