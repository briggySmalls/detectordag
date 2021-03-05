# Connection lambdas

Part of the detector dag state is the current connection status of any given dag.
This is in order reassure a user that their dag is working correctly, or to notify them when they cannot expect a power status update (because it is disconnected).

This is implemented using the AWS IoT [lifecycle events](https://docs.aws.amazon.com/iot/latest/developerguide/life-cycle-events.html):

- $aws/events/presence/connected/_clientId_
- $aws/events/presence/disconnected/_clientId_

These signals are debounced for a few reasons:

- Spurious disconnection events can be sent immediately _after_ a connection event, in the case where a device's connection is interrupted and then returns.
- A user is unlikely to be interested in geniune disconnect/reconnect events if they only account for short periods of downtime (e.g. "lost for 10 minutes").

# The solution

1. All connected/disconnected events are listened to by the `listener` lambda function and the connection status and time are saved as a "transient" status in the device shadow.
1. If the new status differs from the "current" status then the `handler` lambda function is scheduled to be executed 15 minutes later.
1. The `handler` lambda function checks if the "transient" status was the same that triggered it's execution, and if it is the "current" status is updated.

