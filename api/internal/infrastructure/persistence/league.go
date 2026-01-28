package persistence

type League struct {
	Base

	Name     string         `gorm:"not null;uniqueIndex:ux_league_name"`
	IsActive bool           `gorm:"default:true"`
	Members  []LeagueMember `gorm:"foreignKey:LeagueID"`
	Games    []Game         `gorm:"foreignKey:LeagueID"`
}
