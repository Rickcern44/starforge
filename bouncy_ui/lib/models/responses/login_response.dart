class LoginResponse {
  final String token;
  final String tokenType;

  LoginResponse({required this.token, required this.tokenType});

  factory LoginResponse.fromJson(Map<String, dynamic> json) {
    return LoginResponse(
      token: json['token'],
      tokenType: json['userId'],
    );
  }
}
