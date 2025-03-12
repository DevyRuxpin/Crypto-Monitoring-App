import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 20 },
        { duration: '1m', target: 50 },
        { duration: '30s', target: 100 },
        { duration: '30s', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'],
        http_req_failed: ['rate<0.01'],
    },
};

const BASE_URL = 'http://localhost:8080';
let authToken = '';

export function setup() {
    const loginRes = http.post(`${BASE_URL}/api/v1/auth/login`, {
        email: 'test@example.com',
        password: 'test123',
    });
    
    check(loginRes, {
        'logged in successfully': (resp) => resp.json('token') !== '',
    });
    
    return { token: loginRes.json('token') };
}

export default function(data) {
    const headers = {
        'Authorization': `Bearer ${data.token}`,
        'Content-Type': 'application/json',
    };

    // Test portfolio endpoint
    const portfolioRes = http.get(`${BASE_URL}/api/v1/portfolio`, { headers });
    check(portfolioRes, {
        'portfolio status was 200': (r) => r.status === 200,
    });
    sleep(1);

    // Test market data endpoint
    const marketRes = http.get(`${BASE_URL}/api/v1/market/price/BTC`, { headers });
    check(marketRes, {
        'market data status was 200': (r) => r.status === 200,
    });
    sleep(1);

    // Test adding to portfolio
    const addPortfolioRes = http.post(`${BASE_URL}/api/v1/portfolio`, JSON.stringify({
        symbol: 'BTC',
        quantity: 0.1
    }), { headers });
    check(addPortfolioRes, {
        'add portfolio status was 200': (r) => r.status === 200,
    });
    sleep(1);
}

export function teardown(data) {
    // Cleanup if needed
}