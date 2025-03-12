import { Component, OnInit, OnDestroy } from '@angular/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';
import { MarketService } from '../../services/market.service';
import { MarketData } from '../../models/types';

@Component({
    selector: 'app-market-data',
    templateUrl: './market-data.component.html',
    styleUrls: ['./market-data.component.scss']
})
export class MarketDataComponent implements OnInit, OnDestroy {
    marketData: MarketData[] = [];
    displayedColumns: string[] = ['symbol', 'price', 'change24h', 'volume24h'];
    private destroy$ = new Subject<void>();

    constructor(private marketService: MarketService) {}

    ngOnInit() {
        this.marketService.getMarketUpdates()
            .pipe(takeUntil(this.destroy$))
            .subscribe(data => {
                this.marketData = data;
            });

        // Subscribe to default symbols
        ['BTC', 'ETH', 'XRP'].forEach(symbol => {
            this.marketService.subscribeToSymbol(symbol);
        });
    }

    ngOnDestroy() {
        this.destroy$.next();
        this.destroy$.complete();
    }
}