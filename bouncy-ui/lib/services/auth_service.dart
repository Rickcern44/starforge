import 'dart:convert';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../models/user.dart';
import 'api_service.dart';

class AuthService {
  static const _storage = FlutterSecureStorage();
  static const _accessTokenKey = 'access_token';
  static const _userKey = 'current_user';
  
  static User? _cachedUser;

  /// Login - returns true if successful
  static Future<bool> login(String email, String password) async {
    try {
      final response = await ApiService.post('/auth/login', data: {
        'email': email,
        'password': password,
      });

      if (response.statusCode == 200) {
        final data = jsonDecode(response.body);
        await saveAccessTokens(data['token']);
        
        // Fetch user profile immediately after login
        await fetchAndSaveUserProfile();
        return true;
      }
      return false;
    } catch (e) {
      print('Login error: $e');
      return false;
    }
  }

  /// Fetch user profile from API and save locally
  static Future<User?> fetchAndSaveUserProfile() async {
    try {
      final response = await ApiService.get('/users/me');
      if (response.statusCode == 200) {
        final user = User.fromJson(response.body);
        await _storage.write(key: _userKey, value: user.toJson());
        _cachedUser = user;
        return user;
      }
    } catch (e) {
      print('Fetch profile error: $e');
    }
    return null;
  }

  /// Get the current user (from memory or storage)
  static Future<User?> getCurrentUser() async {
    if (_cachedUser != null) return _cachedUser;
    
    final userJson = await _storage.read(key: _userKey);
    if (userJson != null) {
      _cachedUser = User.fromJson(userJson);
      return _cachedUser;
    }
    return null;
  }

  /// Logout
  static Future<void> logout() async {
    await _storage.deleteAll();
    _cachedUser = null;
  }

  static Future<bool> isAuthenticated() async {
    final token = await _storage.read(key: _accessTokenKey);
    return token != null;
  }

  static Future<String?> getAccessToken() async {
    return await _storage.read(key: _accessTokenKey);
  }

  static Future<void> saveAccessTokens(String? accessToken) async {
    if (accessToken != null) {
      await _storage.write(key: _accessTokenKey, value: accessToken);
    }
  }
}
