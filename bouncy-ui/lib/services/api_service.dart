import 'dart:convert';
import 'package:http/http.dart' as http;
import 'auth_service.dart';

class ApiConfig {
  // make this better in the future wih app settings or env
  static const String baseUrl = 'http://localhost:3000/api/v1';
}

class ApiService {
  // Generic GET with path building and query params
  static Future<http.Response> get(
    String path, {
    Map<String, dynamic>? queryParams,
  }) async {
    final uri = _buildUri(path, queryParams: queryParams);
    final token = await AuthService.getAccessToken();

    return http.get(
      uri,
      headers: {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      },
    );
  }

  // Generic POST with smart body handling
  static Future<http.Response> post(
    String path, {
    Map<String, dynamic>? data,
    Map<String, dynamic>? queryParams,
  }) async {
    final uri = _buildUri(path, queryParams: queryParams);
    final token = await AuthService.getAccessToken();

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

    return http.delete(
      uri,
      headers: {
        'Content-Type': 'application/json',
        if (token != null) 'Authorization': 'Bearer $token',
      },
    );
  }

  // Smart URI builder
  static Uri _buildUri(
    String path, {
    Map<String, dynamic>? queryParams,
  }) {
    // Clean path (remove leading/trailing slashes)
    final cleanPath = path.trim().replaceAll(RegExp(r'^/+|/+$'), '');

    // Build base + path
    final basePath =
        '${ApiConfig.baseUrl.trim().replaceAll(RegExp(r'^/+|/+$'), '')}/$cleanPath';

    // Add query params
    if (queryParams != null && queryParams.isNotEmpty) {
      return Uri.parse(basePath)
          .replace(queryParameters: _encodeQueryParams(queryParams));
    }

    return Uri.parse(basePath);
  }

  // Handle nested query params and proper encoding
  static Map<String, String> _encodeQueryParams(Map<String, dynamic> params) {
    return params.map((key, value) => MapEntry(
          key,
          value.toString(),
        ));
  }
}
