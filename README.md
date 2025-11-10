# Subscription Service API

REST API CRUDL service based on Clean Architecture for managing user subscriptions and calculating their total costs over a selected period.  
Built with **PostgreSQL** (with migrations for database initialization).

## Features

**CRUDL operations** for subscriptions:

- **Create** – add a new subscription
- **Read** – get a single subscription by ID
- **Update** – modify an existing subscription
- **Delete** – remove a subscription
- **List** – retrieve all subscriptions with optional filters

Each subscription record includes:

- Service name
- Monthly cost (in RUB)
- User ID (UUID)
- Start date (month & year)
- Optional end date

**Summary endpoint:**

- Calculate the total cost of subscriptions for a given period
- Filter by user ID and/or service name

## Tech Stack

* **Language:** Go 1.24
* **Web Framework:** Gin
* **Database:** PostgreSQL
* **Database Driver:** jackc/pgx/v5
* **Containerization:** Docker

## Implementation Details

* Clean Architecture
* Database migrations
* Data validation
* Pagination
* Context handling
