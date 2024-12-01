# **Test for Go Developers (Mid-level)**

### **Part 1: Theoretical Questions**

1. **What are goroutines, and how do they differ from traditional threads?**
   Goroutines are lightweight thread. They are managed by Go runtime (Go scheduler), very lightweight, starting stack
   size is about 2 mb. Also, goroutines can communicate with each other though channels.
   Threads are managed by operation system, it`s start size is bigger (about 1mb), switching is slower.
2. **Explain the concept of channels in Go and provide an example of how they are used for synchronization.**
   Channels are a safe way for goroutines synchronization. Channel supports read and write operations. Channels can be
   buffered and unbuffered. Buffered channels store specified amount of records. When goroutine writes to channel, if
   buffer is full (or channel is unbuffered), goroutine blocked until other goroutine reads from this channel. Same for
   reading, when goroutine reads from channel, if buffer is empty (or channel is unbuffered), goroutine blocks until
   other goroutine writes to this channel.
   Common case for using channels is waiting for some async operation completion. It can be done by creating a channel
   with empty struct `ch := make(chan struct{})`. Goroutine that process some async operation have to write to this
   channel when processing is complete: `ch <- struct{}{}`. And goroutine that wait for execution to complete have to
   read from this channel: `<-ch`

3. **What are the advantages and disadvantages of Go’s garbage collector compared to manual memory management (e.g., in
   C/C++)?**
   Advantages are:
    - prevents memory leaks and is more safe;
    - makes code more clear and easier to support, especially for large projects;
    - allows to write code faster;
      Because there is no need to manually free memory, less code should be written and fewer mistakes could be done.
      Disadvantages are:
    - less efficient applications. app should be stopped for some time for GC to free memory.Also, GC requires
      additional memory and CPU for tracking links and allocations. GC's work cycles can't be predicted that makes
      harder to predict exact resource consuming and app productivity;
4. **What are slices in Go? How do they differ from arrays, and what are their limitations?**
   Slice is datatype in Go that stores dynamic-size collection of elements of one type. Unlike arrays which have fixed
   size, slices can change their slices during runtime. Slice is a struct which contains pointer to underlying array and
   its capacity. While working with slices, this fact should be kept in mind, because different slices can have pointer
   to same array and changing elements in one slice will cause changing same element in second slice.

5. **How does Go handle error management? How does it differ from exception handling in other languages?**
   Go has special error type (interface with `Error() string` method)for error management. Common practices is: function
   returns error alongside with function results. Then, first way checked if error is nil. If error is not nil, that
   means that function hasn't processed correctly and error should be handled. If error is nil, program can continue
   processing.
   In other programming languages exceptions objects used to notify for processing errors. exceptions can be handled
   with special syntax construction ofter `raise` keyword for initiation and `try-catch-finally` for exception handling.
   Go's approach for error handling is more explicit, it's always easy to understand whether function execution could result
   with error or not. But on the other hand, it requires manually checks for error after each function call and manual
   error propagation, which may be a lot of repetitive code.
   Also, Go has panics - unrecoverable errors that should result in program halting (f.e. `array out of bounds`). It is possible to stop panic propagating with `recover()` function that returns panic's reason. 

### **Part 2: Test Tasks**

1. **Task 1: Implement a Concurrent Worker Pool**
    - **Description**: Write a Go program that launches a fixed number of worker goroutines. These workers should
      process a list of tasks concurrently. Each task should take a random time to complete (simulate it with
      `time.Sleep`), and results should be collected in a synchronized manner using channels.
    - **Expected Output**: Efficient use of goroutines and channels to synchronize task processing, correct handling of
      concurrent tasks, and error management if necessary.
2. **Task 2: REST API Development**
    - **Description**: Develop a simple RESTful API in Go using the `net/http` package or any Go web framework (like
      `Gin`). The API should perform CRUD operations on a "Task" entity (task ID, name, and status).
    - **Expected Output**: Functional RESTful API, proper error handling, clean code, and unit tests for API endpoints.
3. **Task 3: File Processing**
    - **Description**: Write a Go program that reads a large CSV file concurrently and processes each line to extract
      certain information. Ensure proper synchronization to prevent data corruption.
    - **Expected Output**: Efficient file processing using goroutines and channels, and correct handling of concurrency.

### **Part 3: Practical Scenario**

1. **A customer reports that your Go application occasionally crashes when handling large amounts of concurrent traffic.
   How would you troubleshoot and solve this problem?**
  Let`s assume that this case is happened only of prod environment. 
  Firstly, the potential reasons of crash have to be found. For this purpose, I would analyze app logs, traffic logs (some proxy if available), containers logs. Possible problem could be in app processing logic, running out of available resource (cpu, ram, disk memory) on current app's or other apps (other microservices in system) and services (DBs, caches, middlewares).
  Secondly, after potential crash reason was found, it should be reproduced. It may be done with stress test to simulate or reprocess user action which would lead to same load. Then it may be more specific tests for memory leaks, deadlocks or races, goroutines leaks.
  After reproducing and logs re-analyzing, possible solutions are: productivity fixes if problems were in processing logic, or optimization, adding multi-instance support (if possible) and running app in few instances if current amount of instances cant process such amount of requests (data), or adding necessary physical resources (ram, store etc.) to servers, or changing services (dbs, middlewares) configuration to process high load. 