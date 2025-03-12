export interface User {
    id: number;
    email: string;
    token?: string;
}

export interface Portfolio {
    id: number;
    symbol: string;
    quantity: number;
    currentPrice?: number;
    value?: number;
}

export interface MarketData {
    symbol: string;
    price: number;
    change24h: number;
    volume24h: number;
    timestamp: number;
}

export interface Alert {
    id: number;
    symbol: string;
    targetPrice: number;
    type: 'ABOVE' | 'BELOW';
    active: boolean;
}

export interface AuthResponse {
    token: string;
    user: User;
}
