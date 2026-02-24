import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:intl/intl.dart';
import '../models/game.dart';
import '../models/game_attendance.dart';
import '../services/game_service.dart';
import '../services/auth_service.dart';

class EventDetailsScreen extends StatefulWidget {
  final Game game;
  final String leagueName;

  const EventDetailsScreen({
    super.key,
    required this.game,
    required this.leagueName,
  });

  @override
  State<EventDetailsScreen> createState() => _EventDetailsScreenState();
}

class _EventDetailsScreenState extends State<EventDetailsScreen> {
  String? _selectedStatus;
  final TextEditingController _commentController = TextEditingController();
  bool _isUpdating = false;
  late Game _game;

  @override
  void dispose() {
    _commentController.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _game = widget.game;
  }

  Future<void> _handleUpdateAttendance() async {
    if (_selectedStatus == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Please select a status first!')),
      );
      return;
    }

    setState(() {
      _isUpdating = true;
    });

    final statusMap = {'Yes': 0, 'No': 1, 'Tentative': 2};
    final status = statusMap[_selectedStatus!];

    final success = await GameService.updateAttendance(
      gameId: _game.id,
      status: status!,
      comment: _commentController.text,
    );

    if (!mounted) return;

    setState(() {
      _isUpdating = false;
    });

    // If update succeeded, fetch the canonical game and refresh UI
    if (success) {
      try {
        final fetched = await GameService.getGameById(_game.id);
        if (fetched != null) {
          setState(() {
            _game = fetched;
          });
        }
      } catch (_) {}
    }

    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(success
            ? 'Attendance updated successfully!'
            : 'Failed to update attendance.'),
        backgroundColor: success ? Colors.green : Colors.red,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      body: CustomScrollView(
        slivers: [
          SliverAppBar(
            expandedHeight: 200,
            pinned: true,
            flexibleSpace: FlexibleSpaceBar(
              title: Text(
                'Game at ${_game.location}',
                style: GoogleFonts.poppins(
                  fontWeight: FontWeight.bold,
                  color: Colors.white,
                  shadows: [
                    const Shadow(blurRadius: 10, color: Colors.black45)
                  ],
                ),
              ),
              background: Container(
                decoration: BoxDecoration(
                  gradient: LinearGradient(
                    begin: Alignment.topCenter,
                    end: Alignment.bottomCenter,
                    colors: [colorScheme.primary, colorScheme.primaryContainer],
                  ),
                ),
                child: Center(
                  child: Icon(
                    Icons.event,
                    size: 80,
                    color: colorScheme.onPrimary.withOpacity(0.3),
                  ),
                ),
              ),
            ),
          ),
          SliverToBoxAdapter(
            child: Padding(
              padding: const EdgeInsets.all(16.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildInfoCard(context, colorScheme),
                  const SizedBox(height: 24),
                  _buildCheckInSection(colorScheme),
                  const SizedBox(height: 24),
                  Text(
                    'Attendees (${_game.attendance.length})',
                    style: GoogleFonts.poppins(
                        fontSize: 18, fontWeight: FontWeight.bold),
                  ),
                  const SizedBox(height: 12),
                  _buildAttendeesList(colorScheme),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCheckInSection(ColorScheme colorScheme) {
    return Card(
      elevation: 0,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
        side: BorderSide(color: colorScheme.outlineVariant),
      ),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Your Status',
              style: GoogleFonts.poppins(
                  fontWeight: FontWeight.bold, fontSize: 16),
            ),
            const SizedBox(height: 12),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _statusButton('Yes', Icons.check_circle_outline, Colors.green),
                _statusButton('No', Icons.highlight_off, Colors.red),
                _statusButton('Tentative', Icons.help_outline, Colors.orange),
              ],
            ),
            const SizedBox(height: 16),
            TextField(
              controller: _commentController,
              decoration: InputDecoration(
                hintText: 'Add a comment...',
                border:
                    OutlineInputBorder(borderRadius: BorderRadius.circular(12)),
                contentPadding:
                    const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
              ),
            ),
            const SizedBox(height: 12),
            SizedBox(
              width: double.infinity,
              child: ElevatedButton(
                onPressed: _isUpdating ? null : _handleUpdateAttendance,
                child: _isUpdating
                    ? const SizedBox(
                        height: 20,
                        width: 20,
                        child: CircularProgressIndicator(
                            strokeWidth: 2, color: Colors.white),
                      )
                    : const Text('Update Attendance'),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _statusButton(String label, IconData icon, Color color) {
    final isSelected = _selectedStatus == label;
    return InkWell(
      onTap: () => setState(() => _selectedStatus = label),
      child: Column(
        children: [
          Icon(
            icon,
            color: isSelected ? color : Colors.grey,
            size: 32,
          ),
          const SizedBox(height: 4),
          Text(
            label,
            style: TextStyle(
              color: isSelected ? color : Colors.grey,
              fontWeight: isSelected ? FontWeight.bold : FontWeight.normal,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAttendeesList(ColorScheme colorScheme) {
    if (_game.attendance.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(20.0),
          child: Text('No one has checked in yet',
              style: TextStyle(
                  color: Colors.grey[600], fontStyle: FontStyle.italic)),
        ),
      );
    }

    return ListView.separated(
      shrinkWrap: true,
      physics: const NeverScrollableScrollPhysics(),
      itemCount: _game.attendance.length,
      separatorBuilder: (_, __) => const SizedBox(height: 8),
      itemBuilder: (context, index) {
        final attendance = _game.attendance[index];

        return Card(
          elevation: 0,
          color: colorScheme.surfaceContainerHighest.withOpacity(0.1),
          shape:
              RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
          child: Padding(
            padding: const EdgeInsets.all(8.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                ListTile(
                  contentPadding: EdgeInsets.zero,
                  leading: CircleAvatar(
                    backgroundColor: colorScheme.secondaryContainer,
                    child: const Icon(Icons.person, size: 20),
                  ),
                  title: Text('User: ${attendance.userId.substring(0, 8)}...',
                      style: const TextStyle(fontWeight: FontWeight.w500)),
                  subtitle:
                      Text(attendance.checkedIn ? 'Checked In' : 'Pending'),
                  trailing: Icon(
                    attendance.checkedIn
                        ? Icons.check_circle
                        : Icons.help_outline,
                    color: attendance.checkedIn ? Colors.green : Colors.orange,
                  ),
                ),
                if (attendance.checkInComment.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(
                        left: 56.0, bottom: 8.0, right: 16.0),
                    child: Container(
                      padding: const EdgeInsets.all(8),
                      decoration: BoxDecoration(
                        color: colorScheme.surface,
                        borderRadius: BorderRadius.circular(8),
                        border: Border.all(
                            color: colorScheme.outlineVariant.withOpacity(0.5)),
                      ),
                      child: Text(
                        attendance.checkInComment,
                        style: TextStyle(
                            fontSize: 13,
                            color: colorScheme.onSurfaceVariant,
                            fontStyle: FontStyle.italic),
                      ),
                    ),
                  ),
              ],
            ),
          ),
        );
      },
    );
  }

  Widget _buildInfoCard(BuildContext context, ColorScheme colorScheme) {
    return Card(
      elevation: 0,
      color: colorScheme.surfaceContainerHighest.withOpacity(0.3),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(16.0),
        child: Column(
          children: [
            _buildInfoRow(Icons.calendar_today, 'Date',
                DateFormat('MMM d, y').format(_game.startTime)),
            const Divider(height: 24),
            _buildInfoRow(Icons.access_time, 'Time',
                DateFormat('HH:mm').format(_game.startTime)),
            const Divider(height: 24),
            _buildInfoRow(Icons.group, 'League', widget.leagueName),
            const Divider(height: 24),
            _buildInfoRow(Icons.location_on, 'Location', _game.location),
          ],
        ),
      ),
    );
  }

  Widget _buildInfoRow(IconData icon, String label, String value) {
    return Row(
      children: [
        Icon(icon, size: 20, color: Colors.grey[600]),
        const SizedBox(width: 12),
        Expanded(
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Text(label,
                  style: const TextStyle(fontSize: 12, color: Colors.grey)),
              Text(
                value,
                style: const TextStyle(fontWeight: FontWeight.w500),
                overflow: TextOverflow.ellipsis,
                maxLines: 2,
              ),
            ],
          ),
        ),
      ],
    );
  }
}
