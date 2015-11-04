# Deis Cluster Monitoring

Deis has the ability to stand up a component to provide a base level of monitoring for the cluster. It makes use of Prometheus, AlertManager, Node Exporter and Cadvisor to provide a base level of monitoring for Deis component uptime, system metrics, etcd health and cluster status.

The system can alert via Slack webhook, Amazon SNS, PagerDuty and plain old email. To enable alerting, set the config for the alert method.

## Slack alerting
deisctl config monitor set slack_webhook_url=http://slack/somewhere
deisctl config monitor set slack_channel=robots

## UI
The prometheus UI can be viewed at... wait how do we protect access?
