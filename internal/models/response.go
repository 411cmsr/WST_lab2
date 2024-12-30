package models

type GetAllPersonsResponse struct {
	Persons []Person `xml:"persons"`
}
type GetPersonResponse struct {
    Person Person `xml:"Person"` 
}


type DeleteResponse struct {
	Status bool `xml:"status"`
}

type SearchPersonResponse struct {
    Persons []Person `xml:"Persons"`
}


type AddPersonResponse struct {
    ID uint `xml:"ID"`
}

type UpdatePersonResponse struct {
    Status bool `xml:"status"`
}