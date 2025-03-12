import { ComponentFixture, TestBed } from '@angular/core/testing';
import { of } from 'rxjs';
import { DashboardComponent } from './dashboard.component';
import { MarketService } from '../../core/services/market.service';
import { PortfolioService } from '../../core/services/portfolio.service';

describe('DashboardComponent', () => {
  let component: DashboardComponent;
  let fixture: ComponentFixture<DashboardComponent>;
  let marketServiceSpy: jasmine.SpyObj<MarketService>;
  let portfolioServiceSpy: jasmine.SpyObj<PortfolioService>;

  beforeEach(async () => {
    const mService = jasmine.createSpyObj('MarketService', ['getPrice']);
    const pService = jasmine.createSpyObj('PortfolioService', ['getPortfolio']);

    await TestBed.configureTestingModule({
      declarations: [ DashboardComponent ],
      providers: [
        { provide: MarketService, useValue: mService },
        { provide: PortfolioService, useValue: pService }
      ]
    }).compileComponents();

    marketServiceSpy = TestBed.inject(MarketService) as jasmine.SpyObj<MarketService>;
    portfolioServiceSpy = TestBed.inject(PortfolioService) as jasmine.SpyObj<PortfolioService>;
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(DashboardComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load portfolio on init', () => {
    const mockPortfolio = [
      { symbol: 'BTC', quantity: 1 },
      { symbol: 'ETH', quantity: 5 }
    ];

    portfolioServiceSpy.getPortfolio.and.returnValue(of(mockPortfolio));
    marketServiceSpy.getPrice.and.returnValue(of(50000));

    fixture.detectChanges();

    expect(portfolioServiceSpy.getPortfolio).toHaveBeenCalled();
    expect(component.portfolio).toEqual(mockPortfolio);
  });
});