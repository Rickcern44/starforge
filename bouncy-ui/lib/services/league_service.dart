import 'dart:convert';
import 'package:logging/logging.dart';
import '../models/league.dart';
import 'api_service.dart';

class LeagueService {
  static final _log = Logger('LeagueService');

  static Future<List<League>> getLeagues() async {
    try {
      _log.info('Fetching leagues for current user...');
      final response = await ApiService.get('/me/leagues');
      
      _log.fine('Response status: ${response.statusCode}');
      
      if (response.statusCode == 200) {
        final dynamic data = jsonDecode(response.body);
        if (data is List) {
          final leagues = data.map((json) => League.fromMap(json)).toList();
          _log.info('Successfully fetched ${leagues.length} leagues');
          return leagues;
        } else {
          _log.warning('Unexpected data format: expected List, got ${data.runtimeType}');
        }
      } else {
        _log.severe('Failed to fetch leagues. Status: ${response.statusCode}, Body: ${response.body}');
      }
      return [];
    } catch (e, stack) {
      _log.severe('Error fetching leagues', e, stack);
      return [];
    }
  }
}
