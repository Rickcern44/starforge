class LeagueMember {
  final String leagueId;
  final String playerId;
  final String role;
  final DateTime joinedAt;

  LeagueMember({
    required this.leagueId,
    required this.playerId,
    required this.role,
    required this.joinedAt,
  });

  factory LeagueMember.fromMap(Map<String, dynamic> map) {
    return LeagueMember(
      leagueId: map['LeagueID'] ?? '',
      playerId: map['PlayerID'] ?? '',
      role: map['Role']?.toString() ?? '',
      joinedAt: DateTime.tryParse(map['JoinedAt'] ?? '') ?? DateTime.now(),
    );
  }
}