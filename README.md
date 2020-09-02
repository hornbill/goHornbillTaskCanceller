# Hornbill Task Canceller

## Command Line Arguments

- instance ___string___ : ID of the instance to connect to
- api ___string___ : The API Key to use to authenticate against your Hornbill instance
- taskref ___string___ : Single Task Reference (format: TSK###)
- listfile ___string___ : File name of file containing list of task references (one per line)
- delete ___boolean___ : Set to true to delete the task(s), defaults to false and the cancellation of the task(s)

Please find also included __open-tasks-on-cancelled-requests.report.txt__ as a Service Manager Report which identifies all Task IDs of tasks which are not completed/cancelled which are connected to cancelled tasks.

Sample:

<code>taskCanceller.exe -instance=test_instance -api=123...def -listfile=c:\data\tasks.csv</code>
  
See the [Hornbill Wiki](https://wiki.hornbill.com/index.php/Task_Cancelation_Utility) for more information.
