# Latest benchmark results

Units: [ns/op (ratio from ORM)]

| Test | ORM | Raw | GORM |
| --- | --- | --- | --- |
| Insert | 22110 (100%) | 15660 (70%) | 63408 (286%) |
| Query | 2949922 (100%) | 3031721 (102%) | 10064958 (341%) |
| QueryLargeStruct | 21340732 (100%) | 20527749 (96%) | 88145729 (413%) |


![graph](./benchmark.png)



Benchmark time: 2018-01-27


#### Compared packages:

- [x] ORM: posener/orm (this package)
- [x] RAW: Direct SQL commands
- [x] GORM: jinzhu/gorm

#### Operations:

- [x] Insert: INSERT operations
- [X] Query: SELECT operations on an object with 2 fields
- [X] QueryLargeStruct: SELECT on an object of ~35 different fields

