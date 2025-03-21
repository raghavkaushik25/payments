# Payments - Money transfering application

## Overview

This project implements a simple concurrent money transfer system in Go. It ensures atomic balance updates while preventing race conditions and overdrafts. Users can transfer money between each other basic HTTP API.

## Features

- Supports concurrent money transfers between users.
- Ensures atomic updates using concurrency-safe data structures.
- Prevents overdrafts (users cannot send more money than they have).
- Handles errors such as:
  - Invalid users
  - Transfers to oneself
  - Insufficient funds
- A simple HTTP API to initiate transactions.

## Setup and Installation

### Installation Steps

1. Clone the repository:
   ```
    git clone https://github.com/raghavkaushik25/payments.git
    cd payments
    go run main.go
2. The application supports 2 endpoints:
    GET /userInfo?userName=XXX
    This endpoint takes userName in the query parameter and responds with
    {"user_name" :"Adam","account_id":"XXXXX","current_balance" : 0}

    POST /transfer
    This endpoint takes a body
    {"from" : "Mark", "to" : "Adam", "Amount" : 12}
    If the transaction goes welll the endpoint responds with
    {"message":"accound Id : 3b5f116d-42a5-49dd-b66f-c30b4e1abe96 has been debited with ammount 1; updated balance is 99","previous_balance":100,"updated_balance":99}
3. To run the tests:
   ```
    cd /bank
    go test -v -count=100 -failfast
