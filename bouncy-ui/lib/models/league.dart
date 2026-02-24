import 'dart:convert';

import 'game.dart';
import 'league_member.dart';

class League {
  final String id;
  final String name;
  final bool isActive;
  final List<LeagueMember> members;
  final List<Game> games;

  League({
    required this.id,
    required this.name,
    required this.isActive,
    required this.members,
    required this.games,
  });

  factory League.fromMap(Map<String, dynamic> map) {
    return League(
      id: map['ID'] ?? '',
      name: map['Name'] ?? '',
      isActive: map['IsActive'] ?? false,
      members: List<LeagueMember>.from(
        ((map['Members'] ?? map['members']) ?? [])
            .map((x) => LeagueMember.fromMap(x)),
      ),
      games: List<Game>.from(
        ((map['Games'] ?? map['games']) ?? []).map((x) => Game.fromMap(x)),
      ),
    );
  }

  factory League.fromJson(String source) => League.fromMap(json.decode(source));
}
