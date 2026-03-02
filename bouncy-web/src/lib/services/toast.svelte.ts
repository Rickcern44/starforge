export type ToastType = 'success' | 'error' | 'info';

export interface Toast {
  id: string;
  message: string;
  type: ToastType;
  duration: number;
}

class ToastService {
  toasts = $state<Toast[]>([]);

  show(message: string, type: ToastType = 'info', duration: number = 3000) {
    const id = Math.random().toString(36).substring(2, 9);
    const toast: Toast = { id, message, type, duration };
    
    this.toasts.push(toast);

    setTimeout(() => {
      this.remove(id);
    }, duration);
  }

  success(message: string, duration?: number) {
    this.show(message, 'success', duration);
  }

  error(message: string, duration?: number) {
    this.show(message, 'error', duration);
  }

  info(message: string, duration?: number) {
    this.show(message, 'info', duration);
  }

  remove(id: string) {
    this.toasts = this.toasts.filter((t) => t.id !== id);
  }
}

export const toastService = new ToastService();
