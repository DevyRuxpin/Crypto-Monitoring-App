import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Portfolio } from '../models/portfolio.model';

@Injectable({
  providedIn: 'root'
})
export class PortfolioService {
  constructor(private http: HttpClient) {}

  getPortfolio(): Observable<Portfolio[]> {
    return this.http.get<Portfolio[]>(`${environment.apiUrl}/portfolio`);
  }

  addToPortfolio(symbol: string, quantity: number): Observable<any> {
    return this.http.post(`${environment.apiUrl}/portfolio`, { symbol, quantity });
  }

  removeFromPortfolio(id: number): Observable<any> {
    return this.http.delete(`${environment.apiUrl}/portfolio/${id}`);
  }
}