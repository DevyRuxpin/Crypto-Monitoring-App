import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable, BehaviorSubject, combineLatest } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Portfolio } from '../models/types';
import { MarketService } from './market.service';

@Injectable({
    providedIn: 'root'
})
export class PortfolioService {
    private portfolioSubject = new BehaviorSubject<Portfolio[]>([]);
    public portfolio$ = this.portfolioSubject.asObservable();

    constructor(
        private http: HttpClient,
        private marketService: MarketService
    ) {
        this.loadPortfolio();
    }

    private loadPortfolio() {
        this.http.get<Portfolio[]>(`${environment.apiUrl}/portfolio`)
            .subscribe(
                portfolio => this.portfolioSubject.next(portfolio),
                error => console.error('Failed to load portfolio:', error)
            );
    }

    getPortfolioWithPrices(): Observable<Portfolio[]> {
        return combineLatest([
            this.portfolio$,
            this.marketService.getMarketUpdates()
        ]).pipe(
            map(([portfolio, marketData]) => {
                return portfolio.map(position => ({
                    ...position,
                    currentPrice: marketData.find(m => m.symbol === position.symbol)?.price,
                    value: position.quantity * (marketData.find(m => m.symbol === position.symbol)?.price || 0)
                }));
            })
        );
    }

    addToPortfolio(symbol: string, quantity: number): Observable<Portfolio> {
        return this.http.post<Portfolio>(`${environment.apiUrl}/portfolio`, { symbol, quantity })
            .pipe(
                map(newPosition => {
                    const currentPortfolio = this.portfolioSubject.value;
                    this.portfolioSubject.next([...currentPortfolio, newPosition]);
                    return newPosition;
                })
            );
    }

    updatePosition(id: number, quantity: number): Observable<Portfolio> {
        return this.http.put<Portfolio>(`${environment.apiUrl}/portfolio/${id}`, { quantity })
            .pipe(
                map(updatedPosition => {
                    const currentPortfolio = this.portfolioSubject.value;
                    const index = currentPortfolio.findIndex(p => p.id === id);
                    if (index !== -1) {
                        currentPortfolio[index] = updatedPosition;
                        this.portfolioSubject.next([...currentPortfolio]);
                    }
                    return updatedPosition;
                })
            );
    }

    removePosition(id: number): Observable<void> {
        return this.http.delete<void>(`${environment.apiUrl}/portfolio/${id}`)
            .pipe(
                map(() => {
                    const currentPortfolio = this.portfolioSubject.value;
                    this.portfolioSubject.next(
                        currentPortfolio.filter(p => p.id !== id)
                    );
                })
            );
    }
}
