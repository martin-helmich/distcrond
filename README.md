distcrond: Distributed cronjob
==============================

Disclaimer
----------

This project is **purely experimental**! It started as a private programming exercise for learning the Go programming
language (it's awesome!) and is still highly unstable. The APIs may change without any notice; I will not care about
backwards compatibility at this point and I will **not** give any assurance that this program works as it is intended.

This means that *distcrond* may change or delete arbitrary files in your computer, flood your filesystem with random log
data, delete your computer, steal your credit card data and buy me a new boat or cause green spots on your skin. You have
been warned.

Motivation
----------

I was looking for a service that can schedule and execute periodic jobs over an entire cluster of hosts, especially ones
that should be run on *any* host, but not on *every* host. While its easy to distribute crontabs with configuration
management utilities like [Salt](http://saltstack.com) (that's awesome, too), software that
can schedule and load-balance the execution of these "run-anywhere-but-not-everywhere" jobs is surprisingly rare.

This is basically a subset of what [Chronos](https://github.com/mesos/chronos) is doing, which is built on Apache Mesos
and appears to be a very well engineered piece of software. However, it was way too complex and
enterprisey for my requirements (plus, I have a slight aversion to anything Java-related).

Features
--------

### Already there

- Define arbitrary jobs. A job is simply a shell command that is run on a server. Just like a cron.
- Define high-resolution intervals. Go timers work with nanosecond precision, so if you want a job to be run every
  5 ns, that's fine (although that'd be crazy. Just because you can, doesn't mean you should).
- Remote execution. *distcrond* executes shell commands over SSH with public key authentication. No agent of any kind on
  the target node required.
- Reporting. Currently, the output of each job run is stored in an Elasticsearch index where it can be further
  processed.
- Rest API for management purposes.

### Planned

- Notification options (ever had a cronjob that had been failing for months and you didn't notice it?)
- Better definition of schedules
- Fuzzyfication of schedules ("run once a day, but I don't care when!") with assurance of uniform distribution.
- Alternate remote execution engines that do not require SSH access (maybe using a *Salt* runner or a custom agent)
- More storage backends for job execution reports (like for example MongoDB)

Installation
------------

No binary package (yet).

1. To install *distcrond*, first [install Go](https://golang.org/doc/install).

2. Then, switch to your `$GOPATH` and run:

    ```
    > go install github.com/martin-helmich/distcrond
    ```

3. Go should build the `distcrond` binary for you.

4. Done (sorry, I know you'd hoped for downloading, `./configure`-ing and `make`-ing all kinds of
   dependency libraries).

Usage
-----

### Running distcrond

Run the daemon:

    ./distcrond

If you call the command without any parameters, *distcrond* will search for job definitions in `/etc/distcron/jobs.d/`
and for nodes in `/etc/distcron/nodes.d`. *distcrond* will also expect that an Elasticsearch server listens at
localhost at port 9200.

You can change this behaviour using appropriate command-line options:

    ./distcrond --jobsDirectory=/foo/jobs --nodesDirectory=/foo/nodes --esHost=elasticsearch.host --esPort=9200

Please note that there is *no way* to start *distcrond* as an actual daemon. This is a different concern that should
not be handled by the service itself. Use a decent init system like systemd for that. If you dont't have systemd, use
[Supervisor](http://supervisord.org).

### Defining nodes

Nodes are defined as JSON files (one per node) in your node configuration directory (usually, `/etc/distcron/nodes.d`):

Example for remote SSH nodes:

```json
# file: node_remote.json
{
    "roles": ["role1", "role2"],
    "connection_type": "ssh",
    "connection_options": {
        "ssh_user": "mhelmich",
        "ssh_host": "your.remote.host:22",
        "ssh_private_key_file": "/home/mhelmich/.ssh/id_rsa"
    }
}
```

Example for local node:

```json
# file: node_local.json
{
    "roles": ["role1", "role2"],
    "connection_type": "local"
}
```

### Defining jobs

Jobs are also defined as JSON files (one per job) in your job configuration directory (usually, `/etc/distcrond/jobs.d`):

```json
# file: job1.json
{
    "description": "Do something!",
    "owners": [
        {
            "name": "Martin Helmich",
            "email_address": "kontakt@martin-helmich.de"
        }
    ],
    "policy": {
        "hosts": "any",
        "roles": ["web", "db"]
    },
    "schedule": {
        "interval": "10m",
        "reference": "03:10"
    },
    "command": ["/usr/bin/rm", "-r", "-f", "/"]
}
```
