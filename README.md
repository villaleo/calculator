# 📟 Calculator API

A lightweight, stateless REST API to perform basic arithmetic operations.

## 🚀 Features

- [x] Perform **addition, subtraction, multiplication, and division** on two numbers.
- [x] Sum an array of numbers.
- [x] Supports floating point numbers.
- [x] Handles errors gracefully.
- [x] Middleware attaches a request ID to each incoming request.

### 👷🏽‍♂️ Future Features

- [ ] Add in rate limiter to prevent misuse of the API.
- [ ] Add in token authentication to prevent unauthorized users from using the API.
- [ ] Add in a database to keep track of all of the calculations that have taken place.
- [ ] Create an associated http client that can work with the Calculator API.
- [ ] Create a frontend that makes use of the API.

## 📡 Base URL

<http://localhost:8080>

## 📖 API Endpoints

### Add Two Numbers

**Endpoint:**

```http
POST /add
```

**Request Body:**

```json
{
  "x": 10.34,
  "y": 45
}
```

**Response:**

```json
{
  "interpretation": "10.34 + 45",
  "result": 55.34
}
```

### Subtract Two Numbers

**Endpoint:**

```http
POST /subtract
```

**Request Body:**

```json
{
  "x": 10.34,
  "y": 45
}
```

**Response:**

```json
{
  "interpretation": "10.34 - 45",
  "result": -34.66
}
```

### Multiply Two Numbers

**Endpoint:**

```http
POST /multiply
```

**Request Body:**

```json
{
  "x": 10.34,
  "y": 45
}
```

**Response:**

```json
{
  "interpretation": "10.34 * 45",
  "result": 465.3
}
```

### Divide Two Numbers

**Endpoint:**

```http
POST /divide
```

**Request Body:**

```json
{
  "x": 10,
  "y": 2
}
```

**Response:**

```json
{
  "interpretation": "10 / 2",
  "result": 5
}
```

### Sum an Array of Numbers

**Endpoint:**

```http
POST /sum
```

**Request Body:**

```json
{
  "numbers": [2, 5, -10, 32.5, -8]
}
```

**Response:**

```json
{
  "interpretation": "2 + 5 - 10 + 32.5 - 8",
  "result": 21.5
}
```
