# Crypto Monitor API Documentation


## Authentication
All API endpoints except /auth/* require JWT authentication via Bearer token.


### Login
POST /api/v1/auth/login
   json
{
    "email": "user@example.com",
    "password": "password123"
}


### Register
POST /api/v1/auth/register

JSON
{
    "email": "user@example.com",
    "password": "password123"
}


### Portfolio Endpoints
Get Portfolio
GET /api/v1/portfolio

Add to Portfolio
POST /api/v1/portfolio

JSON
{
    "symbol": "BTC",
    "quantity": 1.5
}


### Market Data Endpoints
Get Price
GET /api/v1/market/price/{symbol}

Get Market Data
GET /api/v1/market/data/{symbol}


### Alerts Endpoints
Create Alert
POST /api/v1/alerts

JSON
{
    "symbol": "BTC",
    "targetPrice": 50000,
    "type": "ABOVE"
}