# vqs

###  File Format (.bin)

Total size doesn't include header
File Header

<pre>
+--------------------------+
|  Magic Num (0xD1272BAB)  |
+--------------------------+
| total capacity (4 bytes) |
+--------------------------+
|      size (4 bytes)      |
+--------------------------+
| next read off  (4 bytes) |
+--------------------------+
| next write off (4 bytes) |
+--------------------------+
|     data (32 bytes)      |
+--------------------------+
|     data (32 bytes)      |
+--------------------------+
|     data (32 bytes)      |
+--------------------------+
|     data (32 bytes)      |
+--------------------------+
</pre>


## TODO
  - [ ] change http to bin proto

<pre>
+--------------------------+
|     Version(2 bytes)     |
+--------------------------+
|      MAC (64 bytes)      | // TODO: think more about it
+--------------------------+
|     Action (2 bytes)     |
+--------------------------+
| Payload Length (4 bytes) |
+--------------------------+
|         data....         |
+--------------------------+
</pre>
