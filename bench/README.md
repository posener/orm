# Latest benchmark results

Units: [ns/op (ratio from ORM)]

| Test | ORM | Raw | GORM |
| --- | --- | --- | --- |
| Insert | 24531 (100%) | 26351 (107%) | 79685 (324%) |
| Query | 3928662 (100%) | 3065932 (78%) | 11164356 (284%) |
| QueryLargeStruct | 20737543 (100%) | 20016573 (96%) | 116811220 (563%) |


![graph](./benchmark.png)



Benchmark time: 2017-11-27


#### Compared packages:

- [x] ORM: posener/orm (this package)
- [x] RAW: Direct SQL commands
- [x] GORM: jinzhu/gorm

#### Operations:

- [x] Insert: INSERT operations
- [X] Query: SELECT operations on an object with 2 fields
- [X] QueryLargeStruct: SELECT on an object of ~35 different fields

