# Storage architecture

Message are stored in the system in a (yet unspecified) byte format. The messages are saved to a
file system on local storage. When new messages are written to the system (`Enqueue` operation),
new messages are appended to the file. The service then stores the location of individual messages
in primary─index, that maps `messageID`s to byte locations in the file. The file layout looks like
the following:

```
┌─────────────┬───────────┬───────────┬───────────┐
│ File header │ message 1 │ message 2 │ message 3 │
└─────────────┴───────────┴───────────┴───────────┘
```

The file header contains information about the file such as, which queue the file is for. The
messages are laid out like the following: 

```
┌─────────────────┬────────────────┬─────────────────┐
│     4-bytes     │    4-bytes     │     n-bytes     │
│ Message Version │ Message Length │ Message Content │
└─────────────────┴────────────────┴─────────────────┘
```

## Indices

To avoid scanning the entire file then looking for a specific message, indices need to be
implemented. These indices are kept in memory. For the operation of the application two indices
are required: (1) primary index and (2) timeout index.

### Primary index

The primary index maps `messageID`s to byte locations in the file. The index will be implemented
as a linked, balanced binary search tree. The nodes will store two values: (1) the `messageID`
used to find the nodes and (2) the byte location in the file (`uint64`).

### Timeout index

> Each message has a delay (timeout) associated that tells the system when the message should 
> be returned upon request. Messages of which the timeout does not lie in the past may not be
> returned by the system!

The timeout index is a (priority) queue that associates the timeout of the message with the
`messageID` or the byte location in the file.  
