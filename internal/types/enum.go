package types

type Role string

var (
	Admin      Role = "admin"
	Registered Role = "registered"
)

type Status string

var (
	Active  Status = "active"
	Passive Status = "passive"
)
