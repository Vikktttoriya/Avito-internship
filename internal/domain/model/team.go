package model

type Team struct {
	Name    string
	Members []*User
}

func NewTeam(name string, members []*User) *Team {
	return &Team{
		Name:    name,
		Members: members,
	}
}
