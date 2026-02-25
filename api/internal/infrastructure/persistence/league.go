package persistence

type League struct {
	Base

	Name     string         `gorm:"not null;uniqueIndex:ux_league_name" json:"name"`
	IsActive bool           `gorm:"default:true" json:"isActive"`
	Members  []LeagueMember `gorm:"foreignKey:LeagueID" json:"members"`
	Games    []Game         `gorm:"foreignKey:LeagueID" json:"games"`
}
