# **Test for Go Developers (Mid-level)**

### **Part 1: Theoretical Questions**

1. **What are goroutines, and how do they differ from traditional threads?**
2. **Explain the concept of channels in Go and provide an example of how they are used for synchronization.**
3. **What are the advantages and disadvantages of Go’s garbage collector compared to manual memory management (e.g., in C/C++)?**
4. **What are slices in Go? How do they differ from arrays, and what are their limitations?**
5. **How does Go handle error management? How does it differ from exception handling in other languages?**

### **Part 2: Test Tasks**

1. **Task 1: Implement a Concurrent Worker Pool**
    - **Description**: Write a Go program that launches a fixed number of worker goroutines. These workers should process a list of tasks concurrently. Each task should take a random time to complete (simulate it with `time.Sleep`), and results should be collected in a synchronized manner using channels.
    - **Expected Output**: Efficient use of goroutines and channels to synchronize task processing, correct handling of concurrent tasks, and error management if necessary.
2. **Task 2: REST API Development**
    - **Description**: Develop a simple RESTful API in Go using the `net/http` package or any Go web framework (like `Gin`). The API should perform CRUD operations on a "Task" entity (task ID, name, and status).
    - **Expected Output**: Functional RESTful API, proper error handling, clean code, and unit tests for API endpoints.
3. **Task 3: File Processing**
    - **Description**: Write a Go program that reads a large CSV file concurrently and processes each line to extract certain information. Ensure proper synchronization to prevent data corruption.
    - **Expected Output**: Efficient file processing using goroutines and channels, and correct handling of concurrency.

### **Part 3: Practical Scenario**

1. **A customer reports that your Go application occasionally crashes when handling large amounts of concurrent traffic. How would you troubleshoot and solve this problem?**