package cloud

type Cloud interface {
	Write(object, filePath string) error
	Read(object string) ([]byte, error)
	MakePublic(object string) (string, error)
	Delete(object string) error
}
var cloudClient Cloud

func Client() Cloud {
	// TODO use DI (google wire)
	return GCloudLazyInit()
}