# ec2-network-monitor

This is a simple example of monitoring ec2 network metrics that are being presented via ethtool.  This is not yet supported by the datadog agent so this is a temporary workaround.

For more information see [Monitoring EC2 Network Performance ENA](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/monitoring-network-performance-ena.html).

## Rebuild

```
helm package helm/
```

## Install

```
helm repo add ec2nm https://jfreeland.github.io/ec2-network-monitor
helm install -n kube-system ec2-network-monitor ec2nm/ec2-network-monitor
```
