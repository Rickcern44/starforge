import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:logging/logging.dart';
import 'package:ui_test/services/auth_service.dart';

class ApiConfig {
  static const String baseUrl = 'http://localhost:3000/api/v1';
}

class ApiService {
  static final _log = Logger('ApiService');

  // Generic GET
  static Future<http.Response> get(
    String path, {
    Map<String, dynamic>? queryParams,
  }) async {
    final uri = _buildUri(path, queryParams: queryParams);
    final token = await AuthService.getAccessToken();

    _log.fine('GET $uri (Token present: ${token != null})');

    return http.get(
      uri,
      headers: {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      },
    );
  }

  // Generic POST
  static Future<http.Response> post(
    String path, {
    Map<String, dynamic>? data,
    Map<String, dynamic>? queryParams,
  }) async {
    final uri = _buildUri(path, queryParams: queryParams);
    final token = await AuthService.getAccessToken();

    _log.fine('POST $uri (Token present: ${token != null})');

    return http.post(
      uri,
      headers: {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      },
      body: data != null ? jsonEncode(data) : null,
    );
  }

  // PATCH
  static Future<http.Response> patch(
    String path, {
    Map<String, dynamic>? data,
    Map<String, dynamic>? queryParams,
  }) async {
    final uri = _buildUri(path, queryParams: queryParams);
    final token = await AuthService.getAccessToken();

    _log.fine('PATCH $uri (Token present: ${token != null})');

    return http.patch(
      uri,
      headers: {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      },
      body: data != null ? jsonEncode(data) : null,
    );
  }

  // PUT
  static Future<http.Response> put(
    String path, {
    Map<String, dynamic>? data,
    Map<String, dynamic>? queryParams,
  }) async {
    final uri = _buildUri(path, queryParams: queryParams);
    final token = await AuthService.getAccessToken();

    _log.fine('PUT $uri (Token present: ${token != null})');

    return http.put(
      uri,
      headers: {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      },
      body: data != null ? jsonEncode(data) : null,
    );
  }

  // DELETE
  static Future<http.Response> delete(
    String path, {
    Map<String, dynamic>? queryParams,
  }) async {
    final uri = _buildUri(path, queryParams: queryParams);
    final token = await AuthService.getAccessToken();

    _log.fine('DELETE $uri (Token present: ${token != null})');

    return http.delete(
      uri,
      headers: {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      },
    );
  }

  static Uri _buildUri(String path, {Map<String, dynamic>? queryParams}) {
    final cleanPath = path.trim().replaceAll(RegExp(r'^/+|/+$'), '');
    final basePath =
        '${ApiConfig.baseUrl.trim().replaceAll(RegExp(r'^/+|/+$'), '')}/$cleanPath';

    if (queryParams != null && queryParams.isNotEmpty) {
      return Uri.parse(basePath)
          .replace(queryParameters: _encodeQueryParams(queryParams));
    }
    return Uri.parse(basePath);
  }

  static Map<String, String> _encodeQueryParams(Map<String, dynamic> params) {
    return params.map((key, value) => MapEntry(key, value.toString()));
  }
}
