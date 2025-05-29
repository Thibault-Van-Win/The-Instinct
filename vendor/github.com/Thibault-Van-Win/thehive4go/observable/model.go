package observable

type ObservableRequest struct {
	DataType         string     `json:"dataType"`
	Data             []string   `json:"data"`
	Message          string     `json:"message"`
	StartDate        int64      `json:"startDate"`
	Attachment       Attachment `json:"attachment,omitempty"`
	TLP              int        `json:"tlp"`
	PAP              int        `json:"pap"`
	Tags             []string   `json:"tags"`
	IOC              bool       `json:"ioc"`
	Sighted          bool       `json:"sighted"`
	SightedAt        int64      `json:"sightedAt"`
	IgnoreSimilarity bool       `json:"ignoreSimilarity"`
	IsZip            bool       `json:"isZip"`
	ZipPassword      string     `json:"zipPassword"`
}

type Attachment struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	ID          string `json:"id"`
}
