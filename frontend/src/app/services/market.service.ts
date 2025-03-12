import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject } from 'rxjs';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { environment } from '../../environments/environment';
import { MarketData } from '../models/types';

@Injectable({
    providedIn: 'root'
})
export class MarketService {
    private wsSubject: WebSocketSubject<any> | null = null;
    private marketDataSubject = new BehaviorSubject<MarketData[]>([]);

    constructor(private http: HttpClient) {
        this.connectWebSocket();
    }

    private connectWebSocket() {
        this.wsSubject = webSocket({
            url: environment.wsUrl,
            openObserver: {
                next: () => console.log('WebSocket connected')
            }
        });

        this.wsSubject.subscribe(
            (data) => this.handleWebSocketMessage(data),
            (error) => console.error('WebSocket error:', error),
            () => console.log('WebSocket connection closed')
        );
    }

    private handleWebSocketMessage(data: any) {
        if (data.type === 'MARKET_UPDATE') {
            const currentData = this.marketDataSubject.value;
            const index = currentData.findIndex(item => item.symbol === data.payload.symbol);
            
            if (index !== -1) {
                currentData[index] = { ...currentData[index], ...data.payload };
            } else {
                currentData.push(data.payload);
            }
            
            this.marketDataSubject.next([...currentData]);
        }
    }

    getMarketData(symbol: string): Observable<MarketData> {
        return this.http.get<MarketData>(`${environment.apiUrl}/market/data/${symbol}`);
    }

    getMarketUpdates(): Observable<MarketData[]> {
        return this.marketDataSubject.asObservable();
    }

    subscribeToSymbol(symbol: string): void {
        if (this.wsSubject) {
            this.wsSubject.next({ type: 'SUBSCRIBE', payload: { symbol } });
        }
    }

    unsubscribeFromSymbol(symbol: string): void {
        if (this.wsSubject) {
            this.wsSubject.next({ type: 'UNSUBSCRIBE', payload: { symbol } });
        }
    }
}
