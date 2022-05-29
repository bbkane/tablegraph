package datatables

type Column struct {
	Title string `json:"title"`
}

type ColumnDef struct {
	ClassName string `json:"className"`
	Targets   string `json:"targets"`
}

type DataTable struct {
	Data       [][]string  `json:"data"`
	Columns    []Column    `json:"columns"`
	ColumnDefs []ColumnDef `json:"columnDefs"`
	PageLength int         `json:"pageLength"`
}
