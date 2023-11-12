package types

type Role string

var (
	Admin      Role = "admin"
	Mod        Role = "mod"
	Registered Role = "registered"
)

type Status string

var (
	Active  Status = "active"
	Passive Status = "passive"
)
