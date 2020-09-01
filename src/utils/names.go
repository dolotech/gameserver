package utils

import (
	"gameserver/utils/log"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

var nameDictionary []string

func LoadNames(file string) error {
	_, e := os.Stat(file)

	if e != nil {
		log.Error(" ioutil.ReadFile(dicFile) ", e)
		return e
	}

	data, e := ioutil.ReadFile(file)
	if e != nil || len(data) <= 0 {
		log.Error(" ioutil.ReadFile(file) ", e)
		return e
	}
	nameDictionary = strings.Split(string(data), "\n")
	return nil
}

func GetName() string {
	l := len(nameDictionary)
	if l > 0 {
		return nameDictionary[rand.Intn(l)]
	} else {
		log.Error("name len is 0")
	}
	return ""

}
