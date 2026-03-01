import { browser } from '$app/environment';
import * as api from './api';
import type { User } from '$lib/models';

class AuthService {
  user = $state<User | null>(null);
  token = $state<string | null>(null);

  constructor() {
    if (browser) {
      const storedToken = localStorage.getItem('access_token') || this.getCookie('access_token');
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

      // If we have a token but no user, fetch profile automatically
      if (this.token && !this.user) {
        this.fetchAndSaveUserProfile();
      }

      $effect.root(() => {
        $effect(() => {
          api.setToken(this.token);
          if (this.token) {
            localStorage.setItem('access_token', this.token);
            this.setCookie('access_token', this.token, 7);
          } else {
            localStorage.removeItem('access_token');
            this.deleteCookie('access_token');
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

  private setCookie(name: string, value: string, days: number) {
    if (!browser) return;
    const date = new Date();
    date.setTime(date.getTime() + (days * 24 * 60 * 60 * 1000));
    const expires = "; expires=" + date.toUTCString();
    document.cookie = name + "=" + (value || "") + expires + "; path=/; SameSite=Lax";
  }

  private getCookie(name: string): string | null {
    if (!browser) return null;
    const nameEQ = name + "=";
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
      let c = ca[i];
      while (c.charAt(0) === ' ') c = c.substring(1, c.length);
      if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
  }

  private deleteCookie(name: string) {
    if (!browser) return;
    document.cookie = name + '=; Max-Age=-99999999; path=/;';
  }

  async login(email: string, password: string): Promise<boolean> {
    try {
      const response = await api.post('auth/login', { email, password });
      if (response && response.token) {
        this.token = response.token;
        api.setToken(response.token);
        this.setCookie('access_token', response.token, 7);
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
        this.setCookie('access_token', response.token, 7);
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
      this.deleteCookie('access_token');
    }
  }
}

export const authService = new AuthService();
