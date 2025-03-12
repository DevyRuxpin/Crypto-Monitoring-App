import { TestBed } from '@angular/core/testing';
import { HttpClientTestingModule, HttpTestingController } from '@angular/common/http/testing';
import { MarketService } from './market.service';
import { environment } from '../../environments/environment';

describe('MarketService', () => {
    let service: MarketService;
    let httpMock: HttpTestingController;

    beforeEach(() => {
        TestBed.configureTestingModule({
            imports: [HttpClientTestingModule],
            providers: [MarketService]
        });

        service = TestBed.inject(MarketService);
        httpMock = TestBed.inject(HttpTestingController);
    });

    afterEach(() => {
        httpMock.verify();
    });

    it('should retrieve market data', (done) => {
        const mockData = {
            symbol: 'BTC',
            price: 50000,
            change24h: 2.5,
            volume24h: 1000000
        };

        service.getMarketData('BTC').subscribe(data => {
            expect(data).toEqual(mockData);
            done();
        });

        const req = httpMock.expectOne(`${environment.apiUrl}/market/data/BTC`);
        expect(req.request.method).toBe('GET');
        req.flush(mockData);
    });
});
