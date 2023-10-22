package model

import "time"

type Segment struct {
	Id   int
	Name string
}

type SegmentWithExpires struct {
	Id      int
	Name    string
	Expires time.Time
}
