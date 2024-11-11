**Task 1: Implement a Concurrent Worker Pool**
- **Description**: Write a Go program that launches a fixed number of worker goroutines. These workers should process a list of tasks concurrently. Each task should take a random time to complete (simulate it with `time.Sleep`), and results should be collected in a synchronized manner using channels.
- **Expected Output**: Efficient use of goroutines and channels to synchronize task processing, correct handling of concurrent tasks, and error management if necessary.
