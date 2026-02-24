import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import '../models/league.dart';
import '../models/user.dart';
import '../services/league_service.dart';
import '../services/auth_service.dart';

class CreateEventScreen extends StatefulWidget {
  const CreateEventScreen({super.key});

  @override
  State<CreateEventScreen> createState() => _CreateEventScreenState();
}

class _CreateEventScreenState extends State<CreateEventScreen> {
  final _formKey = GlobalKey<FormState>();
  String? _selectedLeague;
  DateTime? _selectedDate;
  TimeOfDay? _selectedTime;
  bool _isRecurring = false;
  String _recurrenceInterval = 'Weekly';

  final List<String> recurrenceOptions = [
    'Daily',
    'Weekly',
    'Bi-weekly',
    'Monthly'
  ];

  List<League> _leagues = [];
  List<League> _editableLeagues = [];
  League? _selectedLeagueObj;
  final TextEditingController _priceController = TextEditingController();
  bool _isLeagueAdmin = false;
  @override
  void initState() {
    super.initState();
    _loadLeagues();
  }

  Future<void> _loadLeagues() async {
    try {
      final fetched = await LeagueService.getLeagues();
      final user = await AuthService.getCurrentUser();
      // determine which leagues the current user can create events in (admin/owner)
      final editable = <League>[];
      if (user != null) {
        for (final l in fetched) {
          final matches = l.members.where((m) => m.playerId == user.id);
          if (matches.isNotEmpty) {
            final role = matches.first.role.toLowerCase();
            if (role.contains('admin') || role.contains('owner')) {
              editable.add(l);
            }
          }
        }
      }

      setState(() {
        _leagues = fetched;
        _editableLeagues = editable;
        if (_editableLeagues.length == 1) {
          _selectedLeagueObj = _editableLeagues.first;
          _selectedLeague = _selectedLeagueObj!.id;
          _isLeagueAdmin = true;
        } else {
          // Clear selection if current selection is not editable
          if (_selectedLeague != null &&
              !_editableLeagues.any((e) => e.id == _selectedLeague)) {
            _selectedLeague = null;
            _selectedLeagueObj = null;
            _isLeagueAdmin = false;
          }
        }
      });
    } catch (_) {}
  }

  @override
  void dispose() {
    _priceController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text('Create Event',
            style: GoogleFonts.poppins(fontWeight: FontWeight.w600)),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16.0),
        child: Form(
          key: _formKey,
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              TextFormField(
                decoration: const InputDecoration(
                  labelText: 'Event Title',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.event_note),
                ),
                validator: (value) => value == null || value.isEmpty
                    ? 'Please enter a title'
                    : null,
              ),
              const SizedBox(height: 16),
              if (_editableLeagues.isEmpty) ...[
                const InputDecorator(
                  decoration: InputDecoration(
                    labelText: 'League',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.group),
                  ),
                  child: Text('No leagues available for creating events'),
                ),
              ] else ...[
                DropdownButtonFormField<String>(
                  decoration: const InputDecoration(
                    labelText: 'League',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.group),
                  ),
                  initialValue: _selectedLeague,
                  items: _editableLeagues
                      .map((l) =>
                          DropdownMenuItem(value: l.id, child: Text(l.name)))
                      .toList(),
                  onChanged: _editableLeagues.length == 1
                      ? null
                      : (val) async {
                          setState(() {
                            _selectedLeague = val;
                            _isLeagueAdmin = false;
                            _selectedLeagueObj = null;
                          });

                          if (val == null) return;

                          final matches =
                              _editableLeagues.where((lk) => lk.id == val);
                          if (matches.isNotEmpty) {
                            _selectedLeagueObj = matches.first;
                          }

                          try {
                            final user = await AuthService.getCurrentUser();
                            if (user != null && _selectedLeagueObj != null) {
                              final memberMatches = _selectedLeagueObj!.members
                                  .where((m) => m.playerId == user.id);
                              if (memberMatches.isNotEmpty) {
                                final role =
                                    memberMatches.first.role.toLowerCase();
                                if (role.contains('admin') ||
                                    role.contains('owner')) {
                                  setState(() => _isLeagueAdmin = true);
                                }
                              }
                            }
                          } catch (_) {}
                        },
                  validator: (value) => value == null || value.isEmpty
                      ? 'Please select a league'
                      : null,
                ),
              ],
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: InkWell(
                      onTap: () async {
                        final date = await showDatePicker(
                          context: context,
                          initialDate: DateTime.now(),
                          firstDate: DateTime.now(),
                          lastDate:
                              DateTime.now().add(const Duration(days: 365)),
                        );
                        if (date != null) setState(() => _selectedDate = date);
                      },
                      child: InputDecorator(
                        decoration: const InputDecoration(
                          labelText: 'Date',
                          border: OutlineInputBorder(),
                          prefixIcon: Icon(Icons.calendar_today),
                        ),
                        child: Text(_selectedDate == null
                            ? 'Select Date'
                            : '${_selectedDate!.toLocal()}'.split(' ')[0]),
                      ),
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: InkWell(
                      onTap: () async {
                        final time = await showTimePicker(
                          context: context,
                          initialTime: TimeOfDay.now(),
                        );
                        if (time != null) setState(() => _selectedTime = time);
                      },
                      child: InputDecorator(
                        decoration: const InputDecoration(
                          labelText: 'Time',
                          border: OutlineInputBorder(),
                          prefixIcon: Icon(Icons.access_time),
                        ),
                        child: Text(_selectedTime == null
                            ? 'Select Time'
                            : _selectedTime!.format(context)),
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 16),
              SwitchListTile(
                title: const Text('Recurring Event'),
                subtitle: const Text('Repeat this event automatically'),
                value: _isRecurring,
                onChanged: (bool value) {
                  setState(() {
                    _isRecurring = value;
                  });
                },
              ),
              if (_isRecurring) ...[
                const SizedBox(height: 8),
                DropdownButtonFormField<String>(
                  decoration: const InputDecoration(
                    labelText: 'Repeat Interval',
                    border: OutlineInputBorder(),
                    prefixIcon: Icon(Icons.repeat),
                  ),
                  initialValue: _recurrenceInterval,
                  items: recurrenceOptions
                      .map((opt) =>
                          DropdownMenuItem(value: opt, child: Text(opt)))
                      .toList(),
                  onChanged: (val) =>
                      setState(() => _recurrenceInterval = val!),
                ),
              ],
              const SizedBox(height: 16),
              TextFormField(
                controller: _priceController,
                decoration: const InputDecoration(
                  labelText: 'Price',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.attach_money),
                ),
                keyboardType:
                    const TextInputType.numberWithOptions(decimal: true),
                validator: (value) {
                  if (value == null || value.isEmpty) return null;
                  final parsed = double.tryParse(value);
                  if (parsed == null) return 'Enter a valid number';
                  if (parsed < 0) return 'Price must be 0 or greater';
                  return null;
                },
              ),
              const SizedBox(height: 16),
              TextFormField(
                maxLines: 3,
                decoration: const InputDecoration(
                  labelText: 'Description',
                  border: OutlineInputBorder(),
                  prefixIcon: Icon(Icons.description),
                ),
              ),
              const SizedBox(height: 32),
              if (_isLeagueAdmin) ...[
                SizedBox(
                  width: double.infinity,
                  height: 50,
                  child: ElevatedButton(
                    style: ElevatedButton.styleFrom(
                      shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(12)),
                    ),
                    onPressed: () {
                      if (_formKey.currentState!.validate()) {
                        Navigator.pop(context);
                      }
                    },
                    child: Text('Create Event',
                        style: GoogleFonts.poppins(
                            fontSize: 16, fontWeight: FontWeight.w600)),
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }
}
