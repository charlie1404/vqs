# VQS (Virtual Queue System)

## Intro

Todo!!

### File Format (.bin)

File Header 16KB

| Data Stored              | Bytes             | Comments                                                                                                                                        |
| ------------------------ | ----------------- | ----------------------------------------------------------------------------------------------------------------------------------------------- |
| Magic Num `(0x01535156)` | 4                 | VQS (`0x56`, `0x51`, `0x53`) Followed By `0x01` (v1) <br/>in little endianness                                                                  |
| Capacity left            | 4                 | in bytes                                                                                                                                        |
| Items Count              | 4                 |                                                                                                                                                 |
| Next Read Offset         | 4                 |                                                                                                                                                 |
| Next Write Offset        | 4                 |                                                                                                                                                 |
| Delay in Seconds         | 2                 |                                                                                                                                                 |
| Max Message Size         | 4                 |                                                                                                                                                 |
| Message Retention Period | 4                 |                                                                                                                                                 |
| Message Wait in Seconds  | 2                 | long poll duration                                                                                                                              |
| Visibility Timeout       | 2                 |                                                                                                                                                 |
| Tags                     | 16KB - 34 = 16350 | <table><tbody><tr><td>TagName</td><td>Null Terminated String</td></tr><tr><td>TagValue</td><td>Null Terminated String</td></tr></tbody></table> |

TODO: Later look how to resize files, or use multiple files
TODO: Page size increase benchmark
TODO: monitoring capablity for page faults
TODO call msync to guarantee disk persist // unix.Msync()

<!--
<pre>
+--------------------------+-----------+
| Date Stored              |   bytes   |
+--------------------------+-----------+
| Magic Num  (0x01535156)  |     4     | VQS (`0x56`, `0x51`, `0x53`) Followed By `0x01` (v1) <br/>in little endianness
+--------------------------+-----------+
| Capacity Left            |     4     |
+--------------------------+-----------+
| Item Count               |     4     |
+--------------------------+-----------+
| Next Read Offset         |     4     |
+--------------------------+-----------+
| Next Write Offset        |     4     |
+--------------------------+-----------+
| Delay in Seconds         |     2     |
+--------------------------+-----------+
| Max Message Size         |     4     |
+--------------------------+-----------+
| Message Retention Period |     4     |
+--------------------------+-----------+
| Message Wait in Seconds  |     2     | // long poll duration
+--------------------------+-----------+
| Visibility Timeout       |     2     |
+--------------------------+-----------+
| Tags                     |   16350   |
+--------------------------+-----------+
</pre>

Tags Structure

<pre>
+----------+----------------------------+
| TagName  |   Null Terminated String   |
+----------+----------------------------+
| TagValue |   Null Terminated String   |
+----------+----------------------------+
</pre>
-->
