# tool-activities-rss

This is a tool for [toolbelt](https://github.com/charlieegan3/toolbelt) which collects activities from a 'well known,
sports-centric, social media platform' and generates an RSS feed item using
[webhook-rss](https://github.com/charlieegan3/tool-webhook-rss) deployed to the same toolbelt.

Example config:

```yaml
tools:
  ...
  activities-rss:
    jobs:
      new-entry:
        schedule: "0 0 5 * * *"
        endpoint: https://...
        email: name@example.com
        password: "secret"
        host: www.example.com
        athlete_id: "1111111"

```
