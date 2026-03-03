import { browser } from '$app/environment';
import * as api from './api';
import type { User } from '$lib/models';

class AuthService {
  user = $state<User | null>(null);
  token = $state<string | null>(null);
  initialized: Promise<void> = browser ? Promise.resolve() : new Promise(() => {});

  constructor() {
    if (browser) {
      const storedToken = localStorage.getItem('access_token');
      const storedUser = localStorage.getItem('current_user');

      if (storedToken) {
        this.token = storedToken;
        api.setToken(storedToken);
      }
      if (storedUser) {
        try {
          this.user = JSON.parse(storedUser);
        } catch (e) {
          console.error('Failed to parse stored user', e);
        }
      }

      // If we have a token, always fetch profile to ensure it is still valid
      if (this.token) {
        this.initialized = this.fetchAndSaveUserProfile().then(() => {});
      } else {
        this.initialized = Promise.resolve();
      }

      $effect.root(() => {
        $effect(() => {
          api.setToken(this.token);
          if (this.token) {
            localStorage.setItem('access_token', this.token);
          } else {
            localStorage.removeItem('access_token');
          }
        });

        $effect(() => {
          if (this.user) {
            localStorage.setItem('current_user', JSON.stringify(this.user));
          } else {
            localStorage.removeItem('current_user');
          }
        });
      });
    }
  }

  async login(email: string, password: string): Promise<boolean> {
    try {
      const response = await api.post('auth/login', { email, password });
      if (response && response.token) {
        this.token = response.token;
        api.setToken(response.token);
        await this.fetchAndSaveUserProfile();
        return true;
      }
      return false;
    } catch (error) {
      console.error('Login error:', error);
      return false;
    }
  }

  async register(name: string, email: string, password: string): Promise<boolean> {
    try {
      const response = await api.post('auth/register', { name, email, password });
      if (response && response.token) {
        this.token = response.token;
        api.setToken(response.token);
        await this.fetchAndSaveUserProfile();
        return true;
      }
      return true;
    } catch (error) {
      console.error('Registration error:', error);
      return false;
    }
  }

  async fetchAndSaveUserProfile(): Promise<User | null> {
    try {
      const profile: User = await api.get('users/me');
      this.user = profile;
      return profile;
    } catch (error) {
      console.error('Fetch profile error:', error);
      this.logout();
      return null;
    }
  }

  logout() {
    this.token = null;
    api.setToken(null);
    this.user = null;
    if (browser) {
      localStorage.removeItem('access_token');
      localStorage.removeItem('current_user');
    }
  }
}

export const authService = new AuthService();
