package rest

type FormFile struct {
	Name, FileName, ContentType, Path string
	Size                              int
}
