package userstore

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sync"

	"github.com/jakobvarmose/crypta/pathing"
)

type Userstore struct {
	path  string
	mutex sync.Mutex
}

func New(myPath string) (*Userstore, error) {
	c := &Userstore{
		path: myPath,
	}
	err := os.MkdirAll(c.path, 0700)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Userstore) UpdateUser(address string, callback func(obj *pathing.Object) error) error {
	if !regexp.MustCompile("^[0-9A-Za-z]+$").MatchString(address) {
		return errors.New("Invalid characters in address")
	}
	c.mutex.Lock()
	defer c.mutex.Unlock()
	obj, err := c.GetUser(address)
	if err != nil {
		return err
	}
	err = callback(obj)
	if err != nil {
		return err
	}
	buf, err := pathing.Marshal(obj.Value())
	if err != nil {
		return err
	}
	buf = append(buf, '\n')
	filename := path.Join(c.path, address)
	err = ioutil.WriteFile(filename+"~", buf, 0700)
	if err != nil {
		return err
	}
	return os.Rename(filename+"~", filename)
}

func (c *Userstore) GetUser(address string) (*pathing.Object, error) {
	if !regexp.MustCompile("^[0-9A-Za-z]+$").MatchString(address) {
		return nil, errors.New("Invalid characters in address")
	}
	filename := path.Join(c.path, address)
	buf, err := ioutil.ReadFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	if os.IsNotExist(err) {
		return pathing.NewObject(nil), nil
	}
	var obj interface{}
	err = pathing.UnmarshalInto(&obj, buf)
	//err = json.Unmarshal(buf, &obj)
	if err != nil {
		return nil, err
	}
	return pathing.NewObject(obj), nil
}
