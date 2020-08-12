package util

func GetObjectURL(server, object string) string {
	return "http://" + server + "/objects/" + object
}
