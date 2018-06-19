package builders

const (
	// JoinLeft - constant for SQL query builder
	JoinLeft = "LEFT"
	// JoinRight - constant for SQL query builder
	JoinRight = "RIGHT"
	// JoinInner - constant for SQL query builder
	JoinInner = "INNER"
)

const (
	queryTypeSelect = "SELECT"
	queryTypeInsert = "INSERT"
	queryTypeSave   = "SAVE"
	queryTypeUpdate = "UPDATE"
	queryTypeDelete = "DELETE"
	tablePrefix     = "t"
	defaultLimit    = 100
)
