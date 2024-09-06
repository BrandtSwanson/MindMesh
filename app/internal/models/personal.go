// personal.go
package models

import "time"

type Personal struct {
	TimeStamp    time.Time `json:"timestamp"`
	ID           int       `json:"id"`
	Water        int       `json:"water"`
	Joyousness   int       `json:"joyousness"`
	Relaxation   int       `json:"relaxation"`
	Alertness    int       `json:"alertness"`
	ScreenTime   int       `json:"screenTime"`
	Satisfaction int       `json:"satisfaction"`
	OutsideTime  int       `json:"outsideTime"`
}
