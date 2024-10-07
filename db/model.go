package db

type CompareData struct {
	Id     int64
	Source *Db     `json:"source"`
	Target *Db     `json:"target"`
	Common *Common `json:"common"`
}

type Db struct {
	User string `json:"user"`
	Psw  string `json:"psw"`
	Host string `json:"host"`
	Port string `json:"port"`
	Db   string `json:"db"`
}

type Common struct {
	Path string `json:"path"`
	Ddl  string `json:"ddl"`
	Dml  string `json:"dml"`
	Name string `json:"name"`
}

func DefaultData() CompareData {
	return CompareData{
		Id: -1,
		Source: &Db{
			User: "root",
			Psw:  "root",
			Host: "localhost",
			Port: "3306",
			Db:   "source",
		},
		Target: &Db{
			User: "root",
			Psw:  "root",
			Host: "localhost",
			Port: "3306",
			Db:   "target",
		},
		Common: &Common{
			Path: "output",
			Ddl:  "ddl",
			Dml:  "dml",
			Name: "source => target",
		},
	}
}
