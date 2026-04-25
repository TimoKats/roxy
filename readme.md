<h1 align="center">📡 roxy</h1>

Roxy is a feature rich RSS proxy written in Go (no dependencies). Use roxy to handle CORS restrictions and combine multiple RSS feeds into one queryable feed that's available in XML and JSON. 

### Getting started
You can install roxy with this command: `go install github.com/TimoKats/roxy`. Next, you can run the roxy cli. It takes two (optional) parameters: "port" (default 2112) and "feeds". Using the feeds parameter, you can load a newsboat url file that will be added to the proxy. 

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
The return feeds are always sorted from newest to oldest and can be filtered using optional parameters. For example, this request gets the 20 most recent posts in the "gaming" category that contain the keywords "RTS" and "RPG" (in XML format).

```url
http://localhost:2113/xml?categories=gaming&keywords=rts,rpg&amount=20
```



| Endpoint   |  Description                         | Parameters                               | Response |
|------------|--------------------------------------|------------------------------------------|----------|
| `/add`     |  Add new feeds to proxy.             | `urls` `category`                        |          |
| `/xml`     |  Retrieve data in XML format.        | `urls` `categories` `keywords` `amount`  | XML      |
| `/json`    |  Retrieve data in JSON format.       | `urls` `categories` `keywords` `amount`  | JSON     |
| `/stats`   |  Get server information.             |                                          | JSON     |
| `/refresh` |  Refresh all RSS feeds.              |                                          |          |

