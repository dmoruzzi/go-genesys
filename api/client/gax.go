// Copyright © 2018 Antoine GIRARD <antoine.girard@sapk.fr>
package client

//go:generate go run ../generate/client.go

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/sapk/go-genesys/api/object"
)

//Login log the client on the GAX instance linked
func (c *Client) Login(user, pass string) (*object.LoginResponse, error) {
	req, err := c.newRequest("POST", "session/login", object.LoginRequest{Username: user, Password: pass, IsPasswordEncrypted: false})
	if err != nil {
		return nil, err
	}
	_, err = c.do(req, nil)

	//Check logged user
	req, err = c.newRequest("GET", "user/info", nil)
	if err != nil {
		return nil, err
	}
	var u object.LoginResponse
	_, err = c.do(req, &u)
	return &u, err
}

//UpdateObject Update a object. The object could be a json string or a go object
func (c *Client) UpdateObject(objType, objID string, v interface{}) (*http.Response, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("cfg/objects/%s/%s", objType, objID), v)
	if err != nil {
		return nil, err
	}
	return c.do(req, v)
}

//PostObject Create a object. The object could be a json string or a go object
func (c *Client) PostObject(v interface{}) (*http.Response, error) {
	req, err := c.newRequest("POST", "cfg/objects", v)
	if err != nil {
		return nil, err
	}
	return c.do(req, v)
}

//ListObject Return all the object of a specific type
func (c *Client) ListObject(objType string, v interface{}) (*http.Response, error) {
	req, err := c.newRequest("GET", "cfg/objects", nil)
	if err != nil {
		return nil, err
	}
	//req.URL.RawQuery = "brief=false&type=" + objType
	parameters := url.Values{}
	parameters.Add("brief", "false")
	parameters.Add("type", objType)
	req.URL.RawQuery = parameters.Encode()

	return c.do(req, v)
}

//GetObjectByID retrieve object with an ID and a type
func (c *Client) GetObjectByID(objType, objID string, v interface{}) (*http.Response, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("cfg/objects/%s/%s", objType, objID), nil)
	if err != nil {
		return nil, err
	}
	return c.do(req, v)
}

//GetObjectByName retrieve object with an name and a type
func (c *Client) GetObjectByName(objType, objName string) (map[string]interface{}, *http.Response, error) {
	req, err := c.newRequest("GET", "cfg/objects", nil)
	if err != nil {
		return nil, nil, err
	}

	//req.URL.RawQuery = "brief=false&type=" + objType + "&name=" + objName
	parameters := url.Values{}
	parameters.Add("brief", "false")
	parameters.Add("type", objType)
	parameters.Add("name", objName)
	req.URL.RawQuery = parameters.Encode()

	var objList []map[string]interface{}
	resp, err := c.do(req, &objList)
	if err != nil {
		return nil, resp, err
	}
	if len(objList) == 0 {
		return nil, resp, errors.New("Object not found")
	}
	//We have at least one object so fill it
	if len(objList) > 1 {
		return objList[0], resp, errors.New("Multiple object matched")
	}
	return objList[0], resp, nil
}

//TODO http://host:8080/gax/api/cfg/tree/CfgApplication/104/path
//TODO http://host:8080/gax/api/cfgobjects/search?type=CfgPerson&iscasesensitive=false&name=test
//TODO http://host:8080/gax/api/cfg/objects?brief=true&filters=folderid%3D102+AND+subtype%3DCFGFolder+OR+folderid%3D102+AND+subtype%3DCFGApplication&type=CfgFolder
//TODO http://host:8080/gax/api/scs/applications hosts alarms solutions
//TODO http://host:8080/gax/api/cfg/appmetadata/102
