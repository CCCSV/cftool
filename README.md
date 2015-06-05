# cftool
Cloud Formation cli based off gist from github.com/nfisher

Use `-prompt` if you want to be prompted for template param values to use, otherwise the default parameters will be used.  

If `-prompt` is not used and you want to override specific values, pass them as name value pairs on the command line.  

Add `-w` to view status of a stack while it is changing state  
By default `-w` (watch) will check every 5 seconds, but only print an update when it detects an overall cluster status change  

If `-v` is added to watch, it will always print the state every 5 seconds, regardless of it has changed.  



Examples
--
Provisioning a stack interactively and watching until it is completed:

    ./cftool -provision -name test-stack-3 -template test-stack.json -w -prompt
    ... (setting params interactively)
    2015-06-04T16:41:27-07:00 CREATE_IN_PROGRESS
    2015-06-04T16:41:42-07:00 CREATE_COMPLETED
    2015-06-04T16:41:42-07:00 Finished


Provisioning a stack using default params and watching until it is completed:  

     ./cftool -provision -name test-stack-3 -template test-stack.json -w
     2015-06-04T16:41:27-07:00 CREATE_IN_PROGRESS
     2015-06-04T16:41:42-07:00 CREATE_COMPLETED
     2015-06-04T16:41:42-07:00 Finished


Provisioning a stack using overridden params and watching until it is completed:  

     ./cftool -provision -name test-stack-3 -w -template test-stack.json param1 value1 param2 value2
     Overriding param1 -> value1
     Overriding param2 -> value2
     2015-06-04T16:41:27-07:00 CREATE_IN_PROGRESS
     2015-06-04T16:41:32-07:00 CREATE_IN_PROGRESS
     2015-06-04T16:41:42-07:00 CREATE_COMPLETED
     2015-06-04T16:41:42-07:00 Finished


Watching the progress of a stack that you kicked off with another tool:  

     ./cftool -w -v -name test-stack-6 -i 1
     2015-06-04T16:41:27-07:00 DELETE_IN_PROGRESS
     2015-06-04T16:41:28-07:00 DELETE_IN_PROGRESS
     2015-06-04T16:41:29-07:00 DELETE_IN_PROGRESS
     2015-06-04T16:41:31-07:00 DELETE_IN_PROGRESS
     ...
     2015-06-04T16:44:35-07:00 Finished
