<h1 align="center">📡 roxy</h1>

Roxy is a feature rich RSS proxy written in Go (no dependencies). Use roxy to handle CORS restrictions and combine multiple RSS feeds into one queryable feed that's available in XML and JSON. 

### Getting started
You can install roxy with this command: `go install github.com/TimoKats/roxy`. Note, I want to add docker and package repositories, so those contributions are welcome!

Next, you can run the roxy cli. It takes two (optional) parameters: "port" (default 2112) and "feeds". Using the feeds parameter, you can load a newsboat url file that will be added to the proxy. 

```console
root@server:~$ roxy -port 2113 -feeds /home/username/.newsboat/urls
2026/04/25 14:34:32 starting roxy...
2026/04/25 14:34:33 added to feed: 'https://www.cbs.nl/nl-nl/rss-feeds/prijzen' economics
2026/04/25 14:34:33 added to feed: 'https://www.cbs.nl/nl-nl/rss-feeds/economie' economics
2026/04/25 14:34:34 added to feed: 'https://news.ycombinator.com/rss' tech
2026/04/25 14:34:34 added to feed: 'https://feeds.nos.nl/nosnieuwstech' tech
2026/04/25 14:34:38 serving on: http://localhost:2113
```

### API
The return feeds are always sorted from newest to oldest and can be filtered using optional parameters. The urls, categories and keywords parameters are all comma seperated lists. The amount parameter is an integer (default 10).

| Endpoint   | Method | Description                         | Parameters                               | Response |
|------------|--------|-------------------------------------|------------------------------------------|----------|
| `/add`     | GET    | Add new feeds to proxy.             | `urls` `category`                        | JSON     |
| `/xml`     | GET    | Retrieve data in XML format.        | `urls` `categories` `keywords` `amount`  | XML      |
| `/json`    | GET    | Retrieve data in JSON format.       | `urls` `categories` `keywords` `amount`  | JSON     |
| `/stat`    | GET    | Get API usage/statistics.           |                                          | JSON     |
| `/refresh` | POST   | Refresh all RSS feeds.              |                                          | JSON     |

For example, this request gets the 20 most recent feeds in the "gaming" category that contain the keywords "RTS" and "RPG" in XML format.

```url
http://localhost:2113/xml?categories=gaming&keywords=rts,rpg&amount=20
```
