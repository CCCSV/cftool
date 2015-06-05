# cftool
Cloud Formation cli based off gist from github.com/nfisher

Added -default flag to use default values when provisioning a stack  

Added -watch and -v to view status of a stack while it is changing state  
By default watch will check every 5 seconds, but only print an update when it detects an overall cluster status change  
If -v is added to watch, it will always print the state every 5 seconds  

Watching the progress of a stack that you kicked off with another tool
    ./cftool -watch -v -name test-stack-6 -i 1
    2015-06-04T16:41:27-07:00 DELETE_IN_PROGRESS
    2015-06-04T16:41:28-07:00 DELETE_IN_PROGRESS
    2015-06-04T16:41:29-07:00 DELETE_IN_PROGRESS
    2015-06-04T16:41:31-07:00 DELETE_IN_PROGRESS
    ...
    2015-06-04T16:44:35-07:00 Finished


    ./cftool -provision -name test-stack-3 -watch -template test-stack.json
    ... (answering params interactively)
    2015-06-04T16:41:27-07:00 CREATE_IN_PROGRESS
    2015-06-04T16:41:32-07:00 CREATE_IN_PROGRESS
    2015-06-04T16:41:42-07:00 CREATE_COMPLETED
    2015-06-04T16:41:42-07:00 Finished
