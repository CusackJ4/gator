# Blog Aggregator Guide
A simple RSS retriever

## Requirements
User must have Postgres and Go installed.

## Installation
Use: go install github.com/cusackj4/gator

## Config
> Create config file: ~/.gatorconfig.json
> Add db location: {"db_url": "postgres://username:@localhost:5432/database?sslmode=disable"}
    <username> is a placeholder for your home folder. 

## Commands
1. Register a user: `gator register <name>`
ex: `gator register paul` registers a user named paul.
Duplicate names are not allowed.

2. Login a user: `gator login <name>`
Logs in a registered user as the currently logged-in user. ex: `gator login paul`.

3. Reset database: `gator reset`
Drops all rows from the database. 

4. Retrieve names of registered users: `gator users`
Also provides the name of the currently logged-in user.

5. Add feed: `gator addfeed <feedname> <url>`
Add a feed to the RSS aggregator. ex: `gator addfeed techcrunch https://techcrunch.com/feed/`

6. Fetch feeds: `gator agg <timeString>`
Ex: `gator agg 30s` - Scrapes the content of the next feed every 30s.

7. View feeds `gator feeds`
View the feeds and who added them.

```bash
Username: jill
The feed name is: TechCrunch
The feed url is: https://techcrunch.com/feed/
Username: jill
The feed name is: HackerNews
The feed url is: https://news.ycombinator.com/rss
Username: paul
The feed name is: Quanta
The feed url is: https://api.quantamagazine.org/feed
Username: jill
The feed name is: Boot.dev
The feed url is: https://www.boot.dev/blog/index.xml
```
'

8. Follow a feed `gator follow <url>`
Ex: `gator follow https://techcrunch.com/feed/`
    If `felix` is the logged-in user, produces the following output:
        `Feed Name: TechCrunch, User Name: felix`

9. See what feeds the logged-in user is following: `gator following`
Ex, if `jill` is the logged-in user, she might see:
```bash
gator following 
jill's feeds:
 - Boot.dev
 - TechCrunch
 - HackerNews
 ```

10. Unfollow a feed `gator unfollow <url>`
ex: `gator unfollow https://techcrunch.com/feed/`
    Using `gator following` will now show that jill (see example 9) no longer follows techcrunch. 

11. Browse your feeds. `gator browse <limit>`
The limit parameter determines how many feeds will be fetched (most recent first)
Ex:
```bash
gator browse 3
Feed: HackerNews
Title: Any Color You Like: NIST Scientists Create 'Any Wavelength' Lasers
Publish Date: 2026-04-18 20:54:17 +0000 +0000
URL: https://www.nist.gov/news-events/news/2026/04/any-color-you-nist-scientists-create-any-wavelength-lasers-tiny-circuits
Description: <a href="https://news.ycombinator.com/item?id=47819453">Comments</a>
Feed: HackerNews
Title: Optimizing Ruby Path Methods
Publish Date: 2026-04-18 20:42:29 +0000 +0000
URL: https://byroot.github.io/ruby/performance/2026/04/18/faster-paths.html
Description: <a href="https://news.ycombinator.com/item?id=47819369">Comments</a>
Feed: HackerNews
Title: PostgreSQL production incident caused by transaction ID wraparound
Publish Date: 2026-04-18 20:34:31 +0000 +0000
URL: https://www.sqlservercentral.com/articles/i-too-have-a-production-story-a-downtime-caused-by-postgres-transaction-id-wraparound-problem
Description: <a href="https://news.ycombinator.com/item?id=47819305">Comments</a>
```



