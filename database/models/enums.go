package models

type Gender int8

const (
	MALE   Gender = 1
	FEMALE Gender = 2
	OTHER  Gender = 3
)

type Occupation int8

const (
	Student    Occupation = 1
	PartTime   Occupation = 2
	FullTime   Occupation = 3
	Retired    Occupation = 4
	FreeLancer Occupation = 5
)

type EventType int8

const (
	Food     EventType = 1
	Ticket   EventType = 2
	Shopping EventType = 3
	Sports   EventType = 4
	Popup    EventType = 5
	Others   EventType = 6
)

type TagIndicator int8

const (
	Keyword     TagIndicator = 1
	Description TagIndicator = 2
	Celebrity   TagIndicator = 3
)
