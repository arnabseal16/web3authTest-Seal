# Problem statement - Create a deployment which can scale, provided the pods don't restart

## Solution:

1. Kubernetes by it's thesis requires pods to be restarted, because the replication controller will find a change in the pod hash and roll out a new pod to be restarted.
2. This solution uses a patch using go-client to add resources to be added in for the said deployment. Not an elegant solution but is the closest implementation that can be achieved

## Future Improvements:

1. Roll it out as an operator with an listener attached to ksm or metrics server to be used as a trigger. ref: https://github.com/arnabseal16/sniff-n-fix . This is an old project of mine used to scale resouces which in other ways would be impossible to scale.
