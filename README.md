# Rodent

The rodent project intends to provide isolated testing of the trento runner.  The runner itself has a large range of unit tests which help ensure the code works as designed.  The code is also often tested 'end-to-end', providing 'real-world' feedback.  

The project intends to act as bridge between unit and end-to-end testing by allowing testing of the runner in a real environment without deploying the full trento stack.

Rodent requires access to port 8080 of the runner and, by default, will listen on port 8000 for the callbacks, however, the callback can be varied.

## What will be tested?

The runner has four API endpoints.  Three of of the endpoints GET data from the runner, these are:

* /ready
* /health
* /catalog

The final endpoint is '/execute', this endpoint is used to POST a request for a check to be run.  When submitting an execute request the result of the request is not returned immediately.  Instead, the runner performs a callback, posting the data to the trent-web when the job is complete.  The job is the rodent is to ß

## Configuring the runner

When testing the runner with rodent, the runner will need to be configured to send the callbacks to the IP address or hostname of the system where rodent is being executed.  It should be a simple case of running the trento runner with the command line arguments `runner start --callbacks-url <IP or hostname of rodent>`.  If running a custom callback port (anything other than 8000) then the runner will need to be passed for full URL, for example if the callback port needs to be 4000 and the IP address of rodent is 172.16.1.5 the runner command should be `runner start --callbacks-url http://172.16.1.5:4000/api/runner/callbacks`

## Add rodent to a running Trento environment
The github action **Docker Image CI** is creating a new docker image on every pull request merge.
These docker image can be find under

`ghcr.io/trento-project/rodent:latest`

The following yaml file will create a new pod  with the name **trento-server-rodent** inside of trento:

```yaml
apiVersion: v1
kind: Pod
metadata:
   name: trento-server-rodent
   labels:
     app.kubernetes.io/instance: trento-server
     app.kubernetes.io/name: rodent
spec:
   containers:
   - name: rodent-pod
     image: ghcr.io/trento-project/rodent:latest
     ports:
     - containerPort: 4000
     command: [ "/bin/sh", "-c", "--" ]
     args: [ "while true; do sleep 30; done;" ]
```


It can be deployed by running:

`kubectl apply -f create-rodent.yaml`


To redirect the check results from trento-server-web to the new created **trento-server-rodent** pod the following commands are needed:

`kubectl patch svc trento-server-web --type merge -p '{"spec":{"selector":{"app.kubernetes.io/name": "rodent"}}}'` \
`kubectl patch svc trento-server-web --type merge -p '{"spec":{"ports": [{"port":4000,"targetPort":4000}]}}'`


After that change the web GUI is not longer able to execute the checks. It is however now possibe to run
the checks via rodent:

Examples:

`kubectl exec -it trento-server-rodent -- /bin/bash -c  '/rodent -r trento-server-runner ExecuteCheck -t trento-client01 -p gcp -c 373DB8 --callbackPort 4000 --callbackUrl "/api/runner/callback" -v'` \
`kubectl exec -it trento-server-rodent -- /bin/bash -c  '/rodent -r trento-server-runner ExecuteAllChecks -t trento-client01 -p gcp --callbackPort 4000 --callbackUrl "/api/runner/callback" -v'`


The following two commands will set everything back to trento-server-web:

`kubectl patch svc trento-server-web --type merge -p '{"spec":{"ports": [{"port":4000,"targetPort":"http"}]}}'` \
`kubectl patch svc trento-server-web --type merge -p '{"spec":{"selector":{"app.kubernetes.io/name": "web"}}}'`



## Commands

rodent supports the following commands:

* CheckReady - retrieves the endpoint '/ready' and reports the status
* CheckHealth - retrieves the endpoint '/health' and reports the status
* CheckCatalog - retrieves the endpoint '/catalog' and reports the status (optionally list the available checks)
* ExecuteCheck - executes a single check, waits for the check to complete and reports the outcome
* ExecuteAllChecks - executes all checks in the catalog and reports on the outcome
* Version - prints the version of rodent

For most commands the flag `runner` must be set so that rodent knows which runner to communicate with.  Below are some examples of running each command:

```bash
# Check if the runner is ready, the hostname of our runner is 'test-runner'.  Runner is set as a global variable so is required before the 'CheckReady' command and can be expressed as '-r' or '--runner'
./rodent CheckReady -r test-runner
# Output will be something like:
# INFO[0000] The status of endpoint 'http://test-runner:8080/api/ready' is 'true'

# Check health of the runner
./rodent CheckHealth -r test-runner
# Output will be something like
# INFO[0000] The status of endpoint 'http://test-runner:8080/api/health' is 'ok'

# Check health of the runner
./rodent CheckCatalog -r test-runner
# Output will be something like
# INFO[0000] Found 136 checks                             
# INFO[0000] ID        Name
# INFO[0000] 156F64    1.1.1
# INFO[0000] 53D035    1.1.1.runtime
# <SNIP>
# The above IDs can be used to run individual checks.

# Perform a single check.  Now the runner needs to target a host and a provider type.  The following example shows testing a system with the hostname 'test01' running on Azure.
# Here the -t flag is used for the Target.  The -p flag sets the provider and -c sets the check ID
# Also, the callback port is specified as 4000 using the --callbackPort flag, if the callback port required is default 8000 this flag can be omitted.
./rodent ExecuteCheck -r test-runner -t test01 -p azure -c 21FCA6 --callbackPort 4000
# Output will be something like
# INFO[0022] CheckID 21FCA6 state is 'passing'
# Note that running CheckReady, CheckHealth and CheckCatalog will return very quickly but executing checks takes longer!

# Perform all checks.  Again, the runner needs the host to test and a provider, we don't need a check ID when running all checks.  The --callbackPort only needs to be set if the runner is configured 
./rodent ExecuteAllChecks -r test-runner --hostToCheck test01 --provider azure --callbackPort 4000
# Output will be similar to this:
# INFO[0016] CheckID 156F64 state is 'passing'            
# INFO[0022] CheckID 53D035 state is 'passing' 
# INFO[0032] CheckID A1244C state is 'passing' 
# <SNIP>

# In cases where the runner is set to a non-default callbacks url, the ExecuteCheck and ExecuteAllChecks commands can specify a custom URL.  This is done by specifying either --callbackUrl or -u followed by the custom url.
./rodent  ExecuteCheck -r test-runner -t test01 -p azure -c 21FCA6 --callbackPort 4000 --callbackUrl "/api/testurl"

# The Version sub-command flag prints the version of runner.  The version is taken for the git tag and will be applied only when the release is created through github.  If the binary you are using reports 'development' as the version, then you do have probably compiled it yourself.  If you're looking for an offical binary, go to https://github.com/trento-project/rodent/releases
```
