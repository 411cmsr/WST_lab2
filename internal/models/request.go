package models

type AddPersonRequest struct {
	Name      string `xml:"Name"`
	Surname   string `xml:"Surname"`
	Age       int    `xml:"Age"`
	Email     string `xml:"Email"`
	Telephone string `xml:"Telephone"`
}

type DeletePersonRequest struct {
	ID int `xml:"ID"`
}


type UpdatePersonRequest struct {
	ID        uint    `xml:"ID"`
	Name      string `xml:"Name"`
	Surname   string `xml:"Surname"`
	Age       int    `xml:"Age"`
	Email     string `xml:"Email"`
	Telephone string `xml:"Telephone"`
}

type GetPersonRequest struct {
	ID uint `xml:"ID"`
}

type GetAllPersonsRequest struct{}

type SearchPersonRequest struct {
	Query string `xml:"Query"`
}
type Body struct {
	AddPerson      *AddPersonRequest      `xml:"AddPerson,omitempty"`
	DeletePerson   *DeletePersonRequest   `xml:"DeletePerson,omitempty"`
	UpdatePerson   *UpdatePersonRequest   `xml:"UpdatePerson,omitempty"`
	GetPerson      *GetPersonRequest      `xml:"GetPerson,omitempty"`
	GetAllPersons  *GetAllPersonsRequest  `xml:"GetAllPersons,omitempty"`
	SearchPerson   *SearchPersonRequest   `xml:"SearchPerson,omitempty"`
}