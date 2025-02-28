# RSS Feed Aggregator (CLI-Based)

## Overview

This is a command-line-based RSS feed aggregator that allows users to register, log in, follow feeds, and fetch feed updates. The application uses PostgreSQL as the database backend and periodically retrieves the latest articles from followed feeds.

## Features

- User authentication (registration & login)
- Follow and unfollow RSS feeds
- Fetch the latest articles from followed feeds
- Periodic automatic feed updates
- CLI-based interaction

## Prerequisites

- Go 1.18+
- PostgreSQL
- RSS feed URLs

## Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/piglitch/blog_aggregator.git
   cd blog_aggregator
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Set up PostgreSQL database:
   ```sh
   psql -U your_user -d your_database -f schema.sql
   ```
4. Update the database connection details in the application.

## Usage

### Register a new user

```sh
./gator register <user_name>
```

### Login

```sh
./gator login <user_name>
```

### Reset database

```sh
./gator reset
```

### Add an RSS feed

```sh
./gator addfeed <feed_url>
```

### List followed feeds

```sh
./gator feeds
```

### Fetch latest articles from followed feeds

```sh
./gator following
```

### Unfollow a feed

```sh
./gator unfollow <feed_url>
```

## Configuration

Ensure that your database credentials and other settings are correctly configured in the application before running it.





