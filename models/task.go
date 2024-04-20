package models

import "time"

type Task struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

func (t *Task) Equals(other *Task) bool {
	aTime := t.CreatedAt.Truncate(time.Second)
	bTime := other.CreatedAt.Truncate(time.Second)

	return t.Id == other.Id &&
		t.Name == other.Name &&
		aTime == bTime

}
