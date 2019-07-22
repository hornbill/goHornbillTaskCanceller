### Task Canceller for Hornbill

  * -taskref ___string___ ----- 
        Single Task Reference (format: TSK###)
	
  * -instance ___string___ ----- 
        The instance to connect to
	
  * -api ___string___ ----- 
	The API Key
	
  * -listfile ___string___ ----- 
        File name of file containing list of task references - one per line
	
  * -zone ___string___ ----- 
        Override the default Zone the instance sits in (default "eur")

Please find also included __open-tasks-on-cancelled-requests.report.txt__ as a Service Manager Report which identifies all Task IDs of tasks which are not completed/cancelled which are connected to cancelled tasks.

Sample:

  Hornbill_Task_Canceller -instance=test_instance -api=123...def -listfile=c:\data\tasks.csv