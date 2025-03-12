import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { BehaviorSubject, Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { User, AuthResponse } from '../models/types';

@Injectable({
    providedIn: 'root'
})
export class AuthService {
    private currentUserSubject: BehaviorSubject<User | null>;
    public currentUser: Observable<User | null>;

    constructor(private http: HttpClient) {
        this.currentUserSubject = new BehaviorSubject<User | null>(
            JSON.parse(localStorage.getItem('currentUser') || 'null')
        );
        this.currentUser = this.currentUserSubject.asObservable();
    }

    public get currentUserValue(): User | null {
        return this.currentUserSubject.value;
    }

    login(email: string, password: string): Observable<User> {
        return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/login`, { email, password })
            .pipe(map(response => {
                const user = response.user;
                user.token = response.token;
                localStorage.setItem('currentUser', JSON.stringify(user));
                this.currentUserSubject.next(user);
                return user;
            }));
    }

    register(email: string, password: string): Observable<User> {
        return this.http.post<AuthResponse>(`${environment.apiUrl}/auth/register`, { email, password })
            .pipe(map(response => {
                const user = response.user;
                user.token = response.token;
                localStorage.setItem('currentUser', JSON.stringify(user));
                this.currentUserSubject.next(user);
                return user;
            }));
    }

    logout(): void {
        localStorage.removeItem('currentUser');
        this.currentUserSubject.next(null);
    }
}
