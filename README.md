# Blog agreGator

[RSS](https://en.wikipedia.org/wiki/RSS) feed aggregator build as a part of [Boot.dev](https://www.boot.dev) course.

Features:
* Add RSS feeds from across the internet to be collected
* Store the collected posts in a PostgreSQL database
* Follow and unfollow RSS feeds that other users have added
* View summaries of the aggregated posts in the terminal, with a link to the full post

Requirements:
* PostgreSQL
* Go
* Goose

Setup:
* Install Goose

```go install github.com/pressly/goose/v3/cmd/goose@latest```
* Create the database - In psql or your preferred client:

```CREATE DATABASE gator;```
* Configure Gator - Create a config file ".gatorconfig.json" in your home directory 
```
{
  "db_url": "connection_string_goes_here",
  "current_user_name": "username_goes_here"
}
```
* Run migrations - From the sql/schema/ directory:
```
goose postgres "postgres://user:password@localhost:5432/gator" up
```
* Install Gator:
```
go install github.com/Khaz713/gator@latest
```

Commands:
* register <userName> - registers a new user with a given name
* login <userName> - changes current user to other exiting one with a given name
* users - lists all users
* addfeed <title> <url> - adds new feed with given title and url
* feeds - lists all the feeds
* follow <url> - make the current user a follower of the feed with give url
* following - lists all the feeds followed by the current user
* unfollow <url> - make the current user stop following the feed with give url
* agg <timeBetweenReqs> - starts collecting posts from feeds until stoped, with a time between each request(1s, 1m, 1h etc.)
* browse <limit> - prints latest <limit> posts from aggregated feeds