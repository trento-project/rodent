# Rodent

The rodent project intends to provide isolated testing of the trento runner.  The runner itself has a large range of unit tests which help ensure the code works as designed.  The code is also often tested 'end-to-end', providing 'real-world' feedback.  

The project intends to act as bridge between unit and end-to-end testing by allowing testing of the runner in a real environment without deploying the full trento stack.

Rodent requires access to port 8080 of the runner and will listen on port 8000 for the callbacks.

## What will be tested?

The runner has four API endpoints.  Three of of the endpoints GET data from the runner, these are:

* /ready
* /health
* /catalog

The final endpoint is '/execute', this endpoint is used to POST a request for a check to be run.  When submitting an execute request the result of the request is not returned immediately.  Instead, the runner performs a callback, posting the data to the trent-web when the job is complete.  The job is the rodent is to ß

## Configuring the runner

When testing the runner with rodent, the runner will need to be configured to send the callbacks to the IP address or hostname of the system where rodent is being executed.  It should be a simple case of running the trento runner with the command line arguments `runner start --callbacks-url <IP or hostname of rodent>`.  If running a custom callback port (anything other than 8000) then the runner will need to be passed for full URL, for example if the callback port needs to be 4000 and the IP address of rodent is 172.16.1.5 the runner command should be `runner start --callbacks-url http://172.16.1.5:4000/api/runner/callbacks`

## Commands

rodent supports the following commands:

* CheckReady - retrieves the endpoint '/ready' and reports the status
* CheckHealth - retrieves the endpoint '/health' and reports the status
* CheckCatalog - retrieves the endpoint '/catalog' and reports the status (optionally list the available checks)
* ExecuteCheck - executes a single check, waits for the check to complete and reports the outcome
* ExecuteAllChecks - executes all checks in the catalog and reports on the outcome

The global flag `runner` must always be set so that rodent knows which runner to communicate with.  Below are some examples of running each command:

```bash
# Check if the runner is ready, the hostname of our runner is 'test-runner'.  Runner is set as a global variable so is required before the 'CheckReady' command and can be expressed as '-r' or '--runner'
./rodent -r test-runner CheckReady
# Output will be something like:
# INFO[0000] The status of endpoint 'http://test-runner:8080/api/ready' is 'true'

# Check health of the runner
./rodent -r test-runner CheckHealth
# Output will be something like
# INFO[0000] The status of endpoint 'http://test-runner:8080/api/health' is 'ok'

# Check health of the runner
./rodent -r test-runner CheckCatalog
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
./rodent -r test-runner ExecuteCheck -t test01 -p azure -c 21FCA6 --callbackPort 4000
# Output will be something like
# INFO[0022] CheckID 21FCA6 state is 'passing'
# Note that running CheckReady, CheckHealth and CheckCatalog will return very quickly but executing checks takes longer!

# Perform all checks.  Again, the runner needs the host to test and a provider, we don't need a check ID when running all checks.  The --callbackPort only needs to be set if the runner is configured 
./rodent --runner test-runner ExecuteAllChecks --hostToCheck test01 --provider azure --callbackPort 4000
# Output will be similar to this:
# INFO[0016] CheckID 156F64 state is 'passing'            
# INFO[0022] CheckID 53D035 state is 'passing' 
# INFO[0032] CheckID A1244C state is 'passing' 
# <SNIP>
```
