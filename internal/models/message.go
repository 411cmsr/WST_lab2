package models

import (
	"encoding/xml"
)


type Envelope struct {
	XMLName xml.Name `xml:"http://www.w3.org/2003/05/soap-envelope Envelope"`
	Header  Header   `xml:"Header"`
	Body    Body     `xml:"Body"`
}


type Header struct {

}




