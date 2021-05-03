package model
type User struct {
	Id		 int        `json:"id" `
	Username string	    `json:"username" db:"username"`
	Password string		`json:"password" db:password`
	Email    string		`json:"email" db:email`
	Gender   string		`json:"gender"`
	Avatar 	 string		`json:"avatar"`
	Desc 	 string		`json:"desc"`
	//gorm.Model
	Song 	 []Song 	`gorm:"many2many:user_song;"`
}


