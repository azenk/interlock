# Foo

```bash
foo -config=config_file.yml
```

## Configuration File

```yaml
trigger: /path/to/trigger.sh
action: /path/to/action.sh
interval: 60
delay: 10
semaphore:
  baz:
    type: file
    id: id_of_this_instance
    path: /path/to/semaphore/file
  bar:
    type: consul
    id: id_of_this_instance
    key: /path/to/key
    token: consul-acl-token
```

### Trigger

The trigger is an executable that indicates when the action is required by returning a zero exit code. Exiting non-zero indicates that no action is required.

### Action

The action to be performed.  This executable should be idempotent.  Returning zero on success.

### Interval

How frequently the trigger script will be executed, in seconds.

### Delay

How long to wait between a trigger being set and the action being performed, in seconds.  Note: This is a minimum delay.  The delay may be longer.

### Semaphores

Configuration of the semaphore that will be used.
