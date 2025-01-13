# homemon

## Project Goals

We need to be able to monitor our internal network services,
where we might need to have an agent running on a linux host
to monitor availability of a service and then report back its
status. Additionally, if we lose power or network connectivity,
our monitoring stack will be unable to alert us of the outage.
To combat this, we will use a host on the internet to bubble
up alerts to us via email or SMS.