### Task Canceller for Hornbill

  -taskref string
        Single Task Reference (format: TSK###)
  -instance string
        The instance to connect to
  -api string
	The API Key
  -listfile string
        File name of file containing list of task references - one per line
  -zone string
        Override the default Zone the instance sits in (default "eur")

Please find also included __open-tasks-on-cancelled-requests.report.txt__ as a Service Manager Report which identifies all Task IDs of tasks which are not completed/cancelled which are connected to cancelled tasks.

Sample:

  Hornbill_Task_Canceller -instance=test_instance -api=123...def -listfile=c:\data\tasks.csv