package interfaces

type GameAttendanceRepository interface {
	Add(gameID, userID string) error
	Remove(gameID, userID string) error
}
