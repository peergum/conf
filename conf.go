package conf

import (
	"errors"
	"flag"
	"fmt"
	"github.com/kylelemons/go-gypsy/yaml"
	"os"
	"path/filepath"
	"strings"
)

const confDir = "/usr/local/conf"

var (
	appName   string
	confFlags = flag.NewFlagSet("confFlags", flag.ContinueOnError)
	env       = confFlags.String("env", "local", "environment (prod|staging|dev|local...)")
	Conf      *YamlConf
)

type YamlConf struct {
	File *yaml.File
	Root yaml.Node
	Conf interface{}
}

func init() {
	var err error
	appName = filepath.Base(os.Args[0])

	var myFlags []string
	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			if confFlags.Lookup(arg) != nil {
				myFlags = append(myFlags, arg)
			}
		}
		confFlags.Parse(myFlags)
	}

	Conf, err = ReadFile(confDir + "/" + appName + "-" + *env + ".yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetConf() *YamlConf {
	return Conf
}

func ReadFile(name string) (*YamlConf, error) {
	if name == "" {
		name = os.Args[0] + "-" + *env + ".yml"
	}
	config, err := yaml.ReadFile(name)
	if err != nil {
		return nil, err
	}
	var theconf = YamlConf{
		File: config,
		Root: config.Root,
	}
	theconf.Conf, err = theconf.Read()
	if err != nil {
		return &theconf, err
	}
	return &theconf, nil
}

func (theconf *YamlConf) Read() (interface{}, error) {
	var err error
	theconf.Conf, err = ReadNode(theconf.Root)
	return theconf.Conf, err
}

func ReadNode(element yaml.Node) (interface{}, error) {
	return FindNode(element, "")
}

func FindNode(element interface{}, path string) (interface{}, error) {
	var err error
	var newelem interface{}
	var paths []string = strings.Split(path, ".")
	//fmt.Println("Paths:", paths, "Element:", element)
	switch element.(type) {
	case yaml.Map:
		if path == "" {
			//fmt.Println("top=", element)
			return element, nil
		}
		elements := (map[string]yaml.Node)(element.(yaml.Map))
		newpath := strings.Join(paths[1:], ".")
		//fmt.Println("newpath=", newpath)
		for name, value := range elements {
			if paths[0] == name {
				//fmt.Println("Path match:", paths[0])
				return FindNode(value, newpath)
			}
		}
		return nil, errors.New("Path " + path + " not found")
	case yaml.Scalar:
		newelem = string(element.(yaml.Scalar))
	case yaml.List:
		if path == "" {
			return element, nil
		}
		elements := ([]yaml.Node)(element.(yaml.List))
		newelem = []interface{}{}
		for _, value := range elements {
			var elem interface{}
			elem, err = FindNode(value, path)
			if err != nil {
				return nil, err
			}
			newelem = append(newelem.([]interface{}), elem)
		}
	default:
		return nil, errors.New("Don't know what to do with" + fmt.Sprintf("%v", element))
	}
	//fmt.Println("Found:", newelem)
	return newelem, nil
}

func GetString(path string) (string, error) {
	element, err := FindNode(Conf.Conf, path)
	if err == nil {
		return "", err
	}
	if result, ok := element.(string); ok {
		return result, nil
	}
	return "", fmt.Errorf("Hey Dummy! that's not a string.")
}

func GetList(element interface{}) ([]yaml.Node, error) {
	switch element.(type) {
	case yaml.List:
		return []yaml.Node(element.(yaml.List)), nil
	}
	return nil, errors.New("Not a list")
}

func (config *YamlConf) GetStringsMap(path string) (stringMap map[string]string, err error) {
	var element interface{}
	element, err = FindNode(config.Conf, path)
	//fmt.Println("element=", element, "err=", err)
	if err != nil {
		return stringMap, err
	}
	//var ok bool
	stringMap = make(map[string]string)
	if interfaceMap, ok := element.(yaml.Map); ok {
		//fmt.Println("map:", interfaceMap)
		for i, v := range interfaceMap {
			//fmt.Printf("%#V=%#v\n", i, v)
			vString, okString := v.(yaml.Scalar)
			if !okString {
				return nil, fmt.Errorf("Hey! that's not a map of strings.")
			}
			stringMap[i] = string(vString)
		}
		return stringMap, nil
	}
	return nil, fmt.Errorf("Hey! that's not a map of strings.")

}
