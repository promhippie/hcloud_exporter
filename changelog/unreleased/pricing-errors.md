Change: Improve pricing error handling

So far we always existed the scraping if there have been any kind of error while
parsing the metric values, from now on we are logging an error but continue to
provide the remaining metrics to avoid loosing unrelated metrics.

https://github.com/promhippie/hcloud_exporter/pull/240
