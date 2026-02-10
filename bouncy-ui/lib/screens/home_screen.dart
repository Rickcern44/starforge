import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:intl/intl.dart';
import '../models/game.dart';
import '../services/auth_service.dart';
import '../services/league_service.dart';
import '../models/user.dart';
import '../models/league.dart';
import 'event_details_screen.dart';

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  String selectedLeagueName = 'All Leagues';
  List<League> _leagues = [];
  bool _isLoading = true;
  User? _currentUser;

  @override
  void initState() {
    super.initState();
    _initHome();
  }

  Future<void> _initHome() async {
    final authenticated = await AuthService.isAuthenticated();
    if (!mounted) return;

    if (!authenticated) {
      Navigator.of(context).pushReplacementNamed('/auth/login');
      return;
    }

    // Load user profile and leagues
    final results = await Future.wait([
      AuthService.getCurrentUser(),
      LeagueService.getLeagues(),
    ]);

    setState(() {
      _currentUser = results[0] as User?;
      _leagues = results[1] as List<League>;
      _isLoading = false;
    });

    // Refresh data in background
    _refreshData();
  }

  Future<void> _refreshData() async {
    final updatedUser = await AuthService.fetchAndSaveUserProfile();
    final updatedLeagues = await LeagueService.getLeagues();
    if (mounted) {
      setState(() {
        if (updatedUser != null) _currentUser = updatedUser;
        _leagues = updatedLeagues;
      });
    }
  }

  List<Game> _getFilteredGames() {
    if (selectedLeagueName == 'All Leagues') {
      List<Game> allGames = [];
      for (var league in _leagues) {
        allGames.addAll(league.games);
      }
      allGames.sort((a, b) => a.startTime.compareTo(b.startTime));
      return allGames;
    } else {
      final league = _leagues.firstWhere((l) => l.name == selectedLeagueName,
          orElse: () => League(
              id: '', name: '', isActive: false, members: [], games: []));
      List<Game> leagueGames = List.from(league.games);
      leagueGames.sort((a, b) => a.startTime.compareTo(b.startTime));
      return leagueGames;
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(
        body: Center(child: CircularProgressIndicator()),
      );
    }

    final colorScheme = Theme.of(context).colorScheme;
    final filteredGames = _getFilteredGames();

    return Scaffold(
      backgroundColor: colorScheme.surfaceVariant.withOpacity(0.3),
      appBar: AppBar(
        title: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'League Manager',
              style: GoogleFonts.poppins(
                  fontWeight: FontWeight.w600, fontSize: 18),
            ),
            if (_currentUser != null)
              Text(
                'Welcome, ${_currentUser!.name}',
                style: GoogleFonts.poppins(
                    fontSize: 12, color: colorScheme.onSurfaceVariant),
              ),
          ],
        ),
        centerTitle: false,
        actions: [
          IconButton(
            icon: const Icon(Icons.notifications_outlined),
            onPressed: () {},
          ),
          IconButton(
            icon: const Icon(Icons.logout),
            onPressed: () async {
              await AuthService.logout();
              if (mounted)
                Navigator.of(context).pushReplacementNamed('/auth/login');
            },
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: RefreshIndicator(
        onRefresh: _refreshData,
        child: SingleChildScrollView(
          physics: const AlwaysScrollableScrollPhysics(),
          child: Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                _buildLeagueSelector(colorScheme),
                const SizedBox(height: 24),
                _buildFinancialSummary(colorScheme),
                const SizedBox(height: 24),
                _buildSectionHeader('Upcoming Events', () {}),
                const SizedBox(height: 12),
                _buildUpcomingEvents(filteredGames),
                const SizedBox(height: 24),
                _buildSectionHeader('Recent Activity', () {}),
                const SizedBox(height: 12),
                _buildRecentActivity(colorScheme),
              ],
            ),
          ),
        ),
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => Navigator.pushNamed(context, '/events/create'),
        label: const Text('New Event'),
        icon: const Icon(Icons.add),
      ),
    );
  }

  Widget _buildLeagueSelector(ColorScheme colorScheme) {
    final leagueNames = ['All Leagues', ..._leagues.map((l) => l.name)];

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Select League',
          style: GoogleFonts.poppins(
            fontSize: 14,
            fontWeight: FontWeight.w500,
            color: colorScheme.onSurfaceVariant,
          ),
        ),
        const SizedBox(height: 8),
        SizedBox(
          height: 40,
          child: ListView.separated(
            scrollDirection: Axis.horizontal,
            itemCount: leagueNames.length,
            separatorBuilder: (_, __) => const SizedBox(width: 8),
            itemBuilder: (context, index) {
              final name = leagueNames[index];
              final isSelected = selectedLeagueName == name;
              return FilterChip(
                label: Text(name),
                selected: isSelected,
                onSelected: (selected) {
                  setState(() {
                    selectedLeagueName = name;
                  });
                },
                backgroundColor: colorScheme.surface,
                selectedColor: colorScheme.primaryContainer,
                labelStyle: TextStyle(
                  color: isSelected
                      ? colorScheme.onPrimaryContainer
                      : colorScheme.onSurface,
                  fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
                ),
              );
            },
          ),
        ),
      ],
    );
  }

  Widget _buildFinancialSummary(ColorScheme colorScheme) {
    return InkWell(
      onTap: () => Navigator.pushNamed(context, '/ledger'),
      child: Card(
        elevation: 0,
        color: colorScheme.primary,
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: Column(
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Outstanding Balance',
                        style: GoogleFonts.poppins(
                          color: colorScheme.onPrimary.withOpacity(0.8),
                          fontSize: 14,
                        ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                        '\$30.00',
                        style: GoogleFonts.poppins(
                          color: colorScheme.onPrimary,
                          fontSize: 28,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: colorScheme.onPrimary.withOpacity(0.2),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Icon(Icons.account_balance_wallet,
                        color: colorScheme.onPrimary),
                  ),
                ],
              ),
              const SizedBox(height: 20),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  _buildFinanceStat('Unpaid Charges', '2 Items',
                      Icons.receipt_long, Colors.orangeAccent),
                  const Icon(Icons.arrow_forward_ios,
                      size: 16, color: Colors.white54),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildFinanceStat(
      String label, String value, IconData icon, Color iconColor) {
    return Row(
      children: [
        Icon(icon, size: 16, color: iconColor),
        const SizedBox(width: 4),
        Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(label,
                style: const TextStyle(color: Colors.white70, fontSize: 12)),
            Text(value,
                style: const TextStyle(
                    color: Colors.white, fontWeight: FontWeight.bold)),
          ],
        ),
      ],
    );
  }

  Widget _buildSectionHeader(String title, VoidCallback onSeeAll) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(
          title,
          style: GoogleFonts.poppins(
            fontSize: 18,
            fontWeight: FontWeight.bold,
          ),
        ),
        TextButton(
          onPressed: onSeeAll,
          child: const Text('See all'),
        ),
      ],
    );
  }

  Widget _buildUpcomingEvents(List<Game> games) {
    if (games.isEmpty) {
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 20),
        child: Center(
          child: Text(
            'No upcoming events',
            style:
                TextStyle(color: Colors.grey[600], fontStyle: FontStyle.italic),
          ),
        ),
      );
    }

    return Column(
      children: games.map((game) {
        final league = _leagues.firstWhere((l) => l.id == game.leagueId,
            orElse: () => League(
                id: '',
                name: 'Unknown League',
                isActive: false,
                members: [],
                games: []));

        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          shape:
              RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
          child: ListTile(
            onTap: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) =>
                      EventDetailsScreen(game: game, leagueName: league.name),
                ),
              );
            },
            contentPadding:
                const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            leading: Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Theme.of(context).colorScheme.secondaryContainer,
                borderRadius: BorderRadius.circular(12),
              ),
              child: Icon(Icons.sports_soccer,
                  color: Theme.of(context).colorScheme.onSecondaryContainer),
            ),
            title: Text('Game at ${game.location}',
                style: const TextStyle(fontWeight: FontWeight.bold)),
            subtitle: Text(
                '${league.name} • ${DateFormat('MMM d').format(game.startTime)}'),
            trailing: Text(
              DateFormat('HH:mm').format(game.startTime),
              style: TextStyle(
                fontWeight: FontWeight.bold,
                color: Theme.of(context).colorScheme.primary,
              ),
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildRecentActivity(ColorScheme colorScheme) {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: colorScheme.surface,
        borderRadius: BorderRadius.circular(16),
      ),
      child: Column(
        children: [
          _buildActivityItem('Registration Fee', '+\$50.00', '2 hours ago',
              Icons.payment, Colors.green),
          const Divider(height: 24),
          _buildActivityItem('Equipment Purchase', '-\$120.00', 'Yesterday',
              Icons.shopping_cart, Colors.red),
        ],
      ),
    );
  }

  Widget _buildActivityItem(String title, String amount, String time,
      IconData icon, Color iconColor) {
    return Row(
      children: [
        CircleAvatar(
          backgroundColor: iconColor.withOpacity(0.1),
          child: Icon(icon, color: iconColor, size: 20),
        ),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(title, style: const TextStyle(fontWeight: FontWeight.w500)),
              Text(time,
                  style: const TextStyle(fontSize: 12, color: Colors.grey)),
            ],
          ),
        ),
        Text(
          amount,
          style: TextStyle(
            fontWeight: FontWeight.bold,
            color: amount.startsWith('+') ? Colors.green : Colors.black,
          ),
        ),
      ],
    );
  }
}
