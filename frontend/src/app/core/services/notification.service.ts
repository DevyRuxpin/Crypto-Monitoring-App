import { Injectable } from '@angular/core';
import { WebSocketSubject, webSocket } from 'rxjs/webSocket';
import { environment } from '../../../environments/environment';
import { AuthService } from './auth.service';

@Injectable({
  providedIn: 'root'
})
export class NotificationService {
  private socket$: WebSocketSubject<any>;

  constructor(private authService: AuthService) {
    this.initializeWebSocket();
  }

  private initializeWebSocket() {
    const user = this.authService.currentUserValue;
    if (user) {
      this.socket$ = webSocket({
        url: `${environment.wsUrl}/notifications?token=${user.token}`,
        deserializer: msg => JSON.parse(msg.data)
      });
    }
  }

  public getNotifications() {
    return this.socket$.asObservable();
  }

  public disconnect() {
    if (this.socket$) {
      this.socket$.complete();
    }
  }
}