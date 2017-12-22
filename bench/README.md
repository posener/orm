# Latest benchmark results

Units: [ns/op (ratio from ORM)]

| Test | ORM | Raw | GORM |
| --- | --- | --- | --- |
| Insert | 21434 (100%) | 14893 (69%) | 61704 (287%) |
| Query | 2956136 (100%) | 2538550 (85%) | 8550167 (289%) |
| QueryLargeStruct | 17244725 (100%) | 18499591 (107%) | 77639139 (450%) |


![graph](./benchmark.png)



Benchmark time: 2017-12-22


#### Compared packages:

- [x] ORM: posener/orm (this package)
- [x] RAW: Direct SQL commands
- [x] GORM: jinzhu/gorm

#### Operations:

- [x] Insert: INSERT operations
- [X] Query: SELECT operations on an object with 2 fields
- [X] QueryLargeStruct: SELECT on an object of ~35 different fields

