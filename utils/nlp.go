package utils

type probabilityMap map[rune]float32

func DetectLanguage(corpus string) (code string, likelyhood float32, err error) {
	buildProbablityMap()
	return "", 0, nil
}

func buildProbablityMap() (pm probabilityMap) {
	pm['A'] = 0.0
	return pm
}

func compareBetweenMaps() (similarity float32) {
	return 0
}
