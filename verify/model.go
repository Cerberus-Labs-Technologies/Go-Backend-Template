package verify

type Verify struct {
	Id       int    `json:"id" db:"ID"`
	UserId   int    `json:"userId" db:"userId"`
	Platform string `json:"platform" db:"platform"`
	Identity string `json:"identity" db:"identity"`
}
