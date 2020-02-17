package bills

import {

}

// Post represents a Social Media Post type.
type Bill struct {
	UUID             string    `json:"uuid"`
	OriginalFileName string    `json:"OriginalFileName"`
	GeneratedFileName string    `json: "GeneratedFileName"`
}


// The init() function is responsible for initializing the mood state
func init() {
	
}

// NewPost is responsible for creating an instance of the Post type.
func NewBill(uuid string, originalfilename string, generatedfilename string) *Bill {
	return &Bill(UUID: uuid, OriginalFileName: originalfilename, GeneratedFileName: generatedfilename)
}
