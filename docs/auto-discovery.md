# Endpoint Discovery

The agent's observers are responsible for discovering service endpoints.  For
these service endpoints to result in a new monitor instance that is watching
that endpoint, you must apply _discovery rules_ to your monitor configuration.
Every monitor that supports monitoring specific services (i.e. not a static
monitor like the `collectd/cpu` monitor) can be configured with a
`discoveryRule` config option that specifies a rule using a mini rule language.

For example, to monitor a Redis instance that has been discovered by a
container-based observer, you could use the following configuration:

```yaml
  monitors:
    - type: collectd/redis
      discoveryRule: containerImage =~ "redis" && port == 6379
```

## Rule DSL

A rule is an expression that is matched against each discovered endpoint to
determine if a monitor should be active for a particular endpoint. The basic
operators are:

| Operator | Description |
| --- | --- |
| == | Equals |
| != | Not equals |
| < | Less than |
| <= | Less than or equal |
| > | Greater than |
| >= | Greater than or equal |
| =~ | Regex matches |
| !~ | Regex does not match |
| && | And |
| || | Or | 

See [the govaluate
documentation](https://github.com/Knetic/govaluate/blob/master/MANUAL.md) for
all available operators.

The variables available in the expression are dependent on what observer you
are using.  The following two variables are common to all observers:

 - `host` (string): The hostname or IP address of the discovered endpoint
 - `port` (integer): The port number of the discovered endpoint


In addition, these extra functions are provided:

 - `Get(map, key)` - retrieves the value from map with the given key
 - `Contains(map, key)` - returns true if key is inside map, otherwise false

There are no implicit rules built into the agent, so each must be specified
manually in the config file in conjunction with monitor that should monitor the
discovered service.