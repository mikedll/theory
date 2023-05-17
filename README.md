
Queues are FIFO, where the next person to be dequeue is on the left.

```
binary semaphore (has value 0 or 1)
sem.wait(), done when entering the critical section, uses 1 capacity
            you can only enter the critical section if you capture the available capacity
            else you block and have to wait until another process calls signal().
            your place in line is preserved.
sem.signal() called when exiting the critical section
             if someone is blocked, wakes them up
             else, increases value from 0 to 1
```

### Problem 1: multiple readers, 1 writer, writer gets priority

#### First solution

reader:

```
bufferSem.wait()
  myThing = buffer.read()
bufferSem.signal()
return myThing
```

writer:

```
given m, where m is something to write
bufferSem.wait()
  buffer.write(m)
bufferSem.signal()
```


Suppose:

```
Reader A enter CS
Reader B tries to enter, blocks (queue: B)
Writer I tries to enter, blocks (queue: B, I)
Reader A exits (queue: B, I)
Reader B enters (queue: I)
Reader B exits (queue: I)
Writer I enters (queue: )
Writer I exits
```

The problem is that Writer I had to wait for Reader B. Writer I should cut in line in front of B, somehow.

#### Second solution

reader:

```
readerSem.wait()
  writerSem.wait()
    myThing = buffer.read()
  writerSem.signal()
readerSem.signal()
return myThing
```

writer:

```
given m, where m is something to write
writerSem.wait()
  buffer.write(m)
writerSem.signal()
```

```
Reader A enter CS
Reader B tries to enter, blocks at readerSem (reader queue: B)
Writer I tries to enter, blocks at writerSem (writer queue: I)
Reader A exits and signals writerSem (reader queue: B)
Writer I captures writerSem (reader queue: B)
Reader A exits and signals readerSem (reader queue: )
Reader B captures readerSem
Reader B calls wait on writerSem, blocks (writer queue: B)
Writer I can write
Writer I signals writerSem (writer queue: )
Reader B captures writerSem
Reader B starts reading
Reader B signals writerSem to 1
Reader B signals readerSem to 1
```


  
### Problem 2: multiple readers, multiple writers, writers have priority

#### First Solution

reader:

```
readerSem.wait()
  writerSem.wait()
    myThing = buffer.read()
  writerSem.signal()
readerSem.signal()
return myThing
```

writer:

```
given m, where m is something to write
writerSem.wait()
  buffer.write(m)
writerSem.signal()
```


```
Reader A enters CS
Reader B waits on readerSem, blocks (reader queue: B)
Writer I waits on writerSem, blocks (writer queue: I)
Writer II waits on writerSem, blocks (writer queue: I, II)
Reader A reads, finishes, signals writerSem (writer queue: II)
Writer I captures writerSem
Reader A signals readerSem (reader queue: )
Reader B captures readerSem
Reader B waits on writerSem, blocks (writer queue: II, B)
Writer III waits on writerSem (writer queue: II, B, III)
Writer I writes, finishes, signals writerSem (writer queue: B, III)
Writer II captures writerSem
Writer II writes, finishes, signals writerSem (writer queue: III)
Reader B captures writerSem
Reader B reads, signals writerSem (writer queue: )
Writer III captures writerSem
Reader B signals readerSem to 1
Writer III writes, finishes, signals writerSem to 1       
```

The problem is that Writer III was not able to cut in front of Reader B.

#### Second solution

reader:

```
readerSem.wait()

  singleReaderSem.wait()
  singleReaderSem.signal()
  
  bufferSem.wait()
    myThing = buffer.read()
  bufferSem.signal()

readerSem.signal()

return myThing
```

writer:

```
given m, where m is something to write

singleReader.wait()
  bufferSem.wait()
    buffer.write(m)
  bufferSem.signal()
singleReader.signal()
```



```
Reader A enters CS
Reader B waits on readerSem, blocks (reader queue: B)
Writer I waits on bufferSem, blocks (buffer queue: I)
Writer II waits on singleReaderSem, blocks (singleReader queue: II)
Reader A reads, finishes, signals bufferSem (buffer queue: )
Writer I captures bufferSem
Reader A signals readerSem (reader queue: )
Reader B captures readerSem
Reader B waits on singleReaderSem (singleReader queue: II, B)
Writer III waits on singleReaderSem (singleReader queue: II, B, III)
Writer I writes, finishes, signals bufferSem to 1
Writer I signals singleReaderSem (singleReader queue: B, III)
Writer II captures singleReaderSem
Writer II captures bufferSem
Writer II writes, finishes, signals bufferSem to 1
Writer II signals singleReaderSem (singleReader queue: III)
Reader B captures singleReaderSem
Reader B signals singleReaderSem (singleReader queue: )
Writer III captures singleReaderSem
```

Writer III and Reader B are racing for bufferSem, which is weird.


  
#### Third Solution (WIP)



reader:

```
readerSem.wait()

  while true
    guardSem.wait()
      myWriters = writers
    guardSem.signal()
    if myWriters == 0
      break
  
  bufferSem.wait()
    myThing = buffer.read()
  bufferSem.signal()

readerSem.signal()

return myThing
```

writer:

```
given m, where m is something to write

guardSem.wait()
  writers += 1
guardSem.signal()

bufferSem.wait()
  buffer.write(m)
bufferSem.signal()

guardSem.wait()
  writers -= 1
guardSem.signal()

```
