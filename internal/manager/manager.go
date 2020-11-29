package manager

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/ermos/annotation/parser"
	"github.com/huetify/back/internal/dbm"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Handler struct {}

type Manager struct {
	HTTP 	struct {
		Method 	string
		RequestURI 	string
	}
	User 	struct {
		Token 		string
		Id        	int
		Issuer    	string
		Username 	string
		Roles 		[]string
		Exp		 	int
	}
	Param		map[string]interface{}
	Query 		map[string]string
	Payload 	map[string]interface{}
	annotation parser.API
	ps		httprouter.Params
	DB 		*dbm.Instance
}

func New (r *http.Request, a parser.API, ps httprouter.Params) *Manager {
	m := Manager{}

	m.annotation = a
	m.ps = ps

	m.setHTTP(r)
	m.setQueries(r)

	return &m
}

func (m *Manager) SetDB (db *dbm.Instance) {
	m.DB = db
}

func (m *Manager) CheckRequest(r *http.Request) (status int, err error) {
	err = m._checkAuthorization(r, m.annotation)
	if err != nil {
		return http.StatusUnauthorized, err
	}

	err = m._getParams(m.ps, m.annotation)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if r.Method == "POST" || r.Method == "PUT" {
		ct := strings.Split(r.Header.Get("Content-Type"), ";")

		switch strings.ToLower(ct[0]) {
		case "application/json":
			err = m.getPayloadJSON(r, m.annotation)
			if err != nil {
				return http.StatusBadRequest, err
			}
		default:
			return http.StatusBadRequest, errors.New(ct[0] + " is not supported by this API")
		}
	}


	return http.StatusNoContent, nil
}

func (m *Manager) _checkAuthorization(r *http.Request, a parser.API) error {
	if len(a.Authorization) != 0 {
		token, err := extractToken(r)

		tokenData, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("HUETIFY_JWT_SECRET")), nil
		})
		if err != nil {
			return errors.New("you must be logged in")
		}

		claims, ok := tokenData.Claims.(jwt.MapClaims)
		if !ok || !tokenData.Valid {
			return err
		}

		var stringRoles []string
		for _, userRole := range claims["roles"].([]interface{}) {
			stringRoles = append(stringRoles, userRole.(string))
		}

		var isGranted bool
		for _, authorization := range a.Authorization {
			if strings.ToLower(authorization) == "all" {
				isGranted = true
			}
			for _, userRole := range stringRoles {
				if strings.ToLower(authorization) == strings.ToLower(userRole) {
					isGranted = true
				}
			}
		}
		if !isGranted {
			return errors.New("unauthorized access to the resource")
		}

		m.User.Token = token
		m.User.Issuer = claims["issuer"].(string)
		m.User.Id = int(claims["id"].(float64))
		m.User.Roles = stringRoles
		m.User.Exp = int(claims["exp"].(float64))
	}
	return nil
}

func (m *Manager) _getParams(ps httprouter.Params, a parser.API) error {
	var value interface{}
	var err error
	result := make(map[string]interface{})

	for _, param := range a.Validate.Params {
		value = nil

		for _, p := range ps {
			if p.Key == param.Key {
				value, err = _convert(param.Type, p.Value)
				if err != nil {
					return fmt.Errorf("%s's type is incorrect for this field", param.Key)
				}
			}
		}

		result[param.Key] = value
	}

	m.Param = result

	return nil
}

func _convert(trueType string, value interface{}) (interface{}, error) {
	var valueString string

	switch value.(type) {
	case int:
		valueString = fmt.Sprintf("%d", value.(int))
	case bool:
		valueString = fmt.Sprintf("%t", value.(bool))
	case float64:
		if trueType != "int" {
			valueString = fmt.Sprintf("%2.f", value.(float64))
		}else{
			valueString = fmt.Sprintf("%0.f", value.(float64))
		}
	case string:
		valueString = value.(string)
	default:
		if trueType == "map" {
			marshal, err := json.Marshal(value)
			if err != nil {
				return nil, errors.New("can't parse map type")
			}

			valueString = string(marshal)
		}else{
			return nil, errors.New("type not found")
		}
	}

	switch strings.ToLower(trueType) {
	case "int":
		rInt, err := strconv.Atoi(valueString)
		if err != nil {
			return rInt, errors.New(`Impossible de convertir ` + valueString + ` en int`)
		}
		return rInt, nil
	case "float64":
		rFloat64, err := strconv.ParseFloat(valueString, 64)
		if err != nil {
			return rFloat64, errors.New(`Impossible de convertir ` + valueString + ` en float64`)
		}
		return rFloat64, nil
	case "bool":
		rBool, err := strconv.ParseBool(valueString)
		if err != nil {
			return rBool, errors.New(`Impossible de convertir ` + valueString + ` en bool`)
		}
		return rBool, nil
	case "string", "map":
		return valueString, nil
	default:
		return value, fmt.Errorf("%s's type is not supported", trueType)
	}
}

func (m *Manager) getPayloadJSON(r *http.Request, a parser.API) error {
	var value interface{}
	var err error
	var data map[string]interface{}
	result := make(map[string]interface{})

	if len(a.Validate.Payload) <= 0 {
		return nil
	}

	err = parseBody(r, &data)
	if err != nil {
		return err
	}

	for _, body := range a.Validate.Payload {
		if !body.Nullable && (data[body.Key] == "" || data[body.Key] == nil) {
			return fmt.Errorf("%s's key is required in payload", body.Key)
		}

		if data[body.Key] == "" || data[body.Key] == nil {
			result[body.Key] = nil
			continue
		}

		value, err = _convert(body.Type, data[body.Key])
		if err != nil {
			return err
		}

		result[body.Key] = value
	}

	m.Payload = result

	return nil
}

func parseBody(r *http.Request, v interface{}) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, v)
}


func (m *Manager) setHTTP(r *http.Request) {
	m.HTTP.Method = r.Method
	m.HTTP.RequestURI = r.RequestURI
}

func (m *Manager) setQueries(r *http.Request) {
	list := make(map[string]string)

	split := strings.Split(r.URL.String(), "?")

	if len(split) < 2 {
		m.Query = list
		return
	}

	query := strings.Split(split[1], "&")
	for _, q := range query {
		split := strings.Split(q, "=")
		if len(split) == 1 {
			list[split[0]] = split[0]
		}else{
			list[split[0]] = split[1]
		}
	}

	m.Query = list
	return
}

func extractToken(r *http.Request) (result string, err error) {
	rToken := r.Header.Get("Authorization")
	split := strings.Split(rToken, "Bearer")
	if len(split) != 2 { return "", errors.New("JWT Token is not valid") }
	token := strings.Replace(split[1], ":", "", -1)
	token = strings.Replace(token, " ", "", -1)
	return token, nil
}
