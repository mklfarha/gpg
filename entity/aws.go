package entity

type AWS struct {
	Region string `json:"region"`
	Key    string `json:"key"`
	Secret string `json:"secret"`
	Bucket string `json:"bucket"`
}
