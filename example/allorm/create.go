// Autogenerated by github.com/posener/orm
package allorm

func (c *TCreate) String() string {
	// Create statement has a line for each variable with it's name and it's type.
	return `CREATE TABLE all ( int INT, text VARCHAR(100) NOT NULL, bool BOOLEAN, PRIMARY KEY (int) )`
}