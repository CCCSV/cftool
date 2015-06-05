# cftool
Cloud Formation cli based off gist from github.com/nfisher

Added -default flag to use default values when provisioning a stack  

Added -watch and -v to view status of a stack while it is changing state

   ./cftool -watch -v -name test-stack-6 -i 1
   2015-06-04T16:41:27-07:00 DELETE_IN_PROGRESS
   2015-06-04T16:41:28-07:00 DELETE_IN_PROGRESS
   2015-06-04T16:41:29-07:00 DELETE_IN_PROGRESS
   2015-06-04T16:41:31-07:00 DELETE_IN_PROGRESS
   ...
   2015-06-04T16:44:35-07:00 Finished

   ./cftool -provision -name test-stack-3 -watch -template test-stack.json
   ... (answering params)
