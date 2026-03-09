import { browser } from '$app/environment';
import { PUBLIC_API_URL } from '$env/static/public';

const BASE_URL = PUBLIC_API_URL || 'http://localhost:3000/api/v1';

let authToken: string | null = null;

export function setToken(token: string | null) {
  authToken = token;
}

function getToken(): string | null {
  if (authToken) return authToken;
  if (!browser) return null;
  return localStorage.getItem('access_token');
}

export type Fetch = typeof fetch;

async function request(method: string, path: string, data?: any, customFetch?: Fetch, tokenOverride?: string | null) {
  const opts: RequestInit = {
    method,
    headers: {
      'Content-Type': 'application/json'
    }
  };

  const token = tokenOverride !== undefined ? tokenOverride : getToken();
  if (token) {
    (opts.headers as any)['Authorization'] = `Bearer ${token}`;
  }

  if (data) {
    opts.body = JSON.stringify(data);
  }

  const f = customFetch || fetch;
  const response = await f(`${BASE_URL}/${path}`, opts);

  if (!response.ok) {
    const errorBody = await response.text();
    let errorMessage = response.statusText;
    try {
      const parsedError = JSON.parse(errorBody);
      errorMessage = parsedError.error || parsedError.message || errorMessage;
    } catch (e) {
      errorMessage = errorBody || errorMessage;
    }
    console.error(`API Error [${method} ${path}]: ${response.status} ${errorMessage}`);
    throw new Error(errorMessage);
  }

  if (response.status === 204) {
    return null;
  }

  const contentType = response.headers.get('Content-Type');
  if (contentType && contentType.includes('application/json')) {
    return response.json();
  }

  return null;
}

export function get(path: string, customFetch?: Fetch, tokenOverride?: string | null) {
  return request('GET', path, undefined, customFetch, tokenOverride);
}

export function post(path: string, data: any, customFetch?: Fetch, tokenOverride?: string | null) {
  return request('POST', path, data, customFetch, tokenOverride);
}

export function patch(path: string, data: any, customFetch?: Fetch, tokenOverride?: string | null) {
  return request('PATCH', path, data, customFetch, tokenOverride);
}

export { request };
