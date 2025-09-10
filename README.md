# Subscription Service API

A simple REST API for managing user subscriptions and calculating their total costs over a selected period.  
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

- **Backend:** REST API (framework of your choice)
- **Database:** PostgreSQL
- **Migrations:** Included for database initialization  
