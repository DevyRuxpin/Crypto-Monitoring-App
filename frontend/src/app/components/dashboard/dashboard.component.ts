import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { Portfolio, MarketData } from '../../models/types';
import { PortfolioService } from '../../services/portfolio.service';
import { MarketService } from '../../services/market.service';

@Component({
    selector: 'app-dashboard',
    templateUrl: './dashboard.component.html',
    styleUrls: ['./dashboard.component.scss']
})
export class DashboardComponent implements OnInit, OnDestroy {
    portfolio: Portfolio[] = [];
    marketData: MarketData[] = [];
    totalValue: number = 0;
    private destroy$ = new Subject<void>();

    constructor(
        private portfolioService: PortfolioService,
        private marketService: MarketService
    ) {}

    ngOnInit() {
        this.portfolioService.getPortfolioWithPrices()
            .pipe(takeUntil(this.destroy$))
            .subscribe(portfolio => {
                this.portfolio = portfolio;
                this.calculateTotalValue();
            });

        this.marketService.getMarketUpdates()
            .pipe(takeUntil(this.destroy$))
            .subscribe(data => {
                this.marketData = data;
                this.calculateTotalValue();
            });
    }

    ngOnDestroy() {
        this.destroy$.next();
        this.destroy$.complete();
    }

    private calculateTotalValue() {
        this.totalValue = this.portfolio.reduce((total, position) => {
            return total + (position.value || 0);
        }, 0);
    }

    addPosition(symbol: string, quantity: number) {
        this.portfolioService.addToPortfolio(symbol, quantity)
            .subscribe(
                () => console.log('Position added successfully'),
                error => console.error('Error adding position:', error)
            );
    }

    updatePosition(id: number, quantity: number) {
        this.portfolioService.updatePosition(id, quantity)
            .subscribe(
                () => console.log('Position updated successfully'),
                error => console.error('Error updating position:', error)
            );
    }

    removePosition(id: number) {
        this.portfolioService.removePosition(id)
            .subscribe(
                () => console.log('Position removed successfully'),
                error => console.error('Error removing position:', error)
            );
    }
}
