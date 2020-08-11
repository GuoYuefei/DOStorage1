package util

func GetObjectURL(server, object string) string {
	return "http://" + server + "/object/" + object
}
