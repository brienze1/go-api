package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/brienze1/notes-api/internal/infra/injections"
	"github.com/brienze1/notes-api/internal/integration/persistence/repositories/models"
	"github.com/brienze1/notes-api/test/mocks"
	"github.com/cucumber/godog"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var tags string

func init() {
	flag.StringVar(&tags, "scenarios", "", "tags to run")
}

func TestFeatures(t *testing.T) {
	flag.Parse()

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			InitializeScenario(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features"},
			Tags:     tags,
			Strict:   true,
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

type testUtils struct {
	uri      string
	headers  map[string]string
	client   *http.Client
	response *response
	db       *gorm.DB
	redis    *redis.Client
	sqs      *mocks.SqsClient
}

type response struct {
	status int
	body   interface{}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	_ = os.Setenv("PROFILE", "test")
	_ = os.Setenv("SERVER_PORT", "8082")
	test := &testUtils{
		uri:    "http://localhost:" + os.Getenv("SERVER_PORT"),
		client: &http.Client{},
		db:     mocks.Db(),
		sqs:    mocks.Sqs(),
	}
	injections.Wire().Db = test.db
	injections.Wire().Sqs = test.sqs

	ctx.Before(
		func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
			test.reset()
			test.startApp()
			return ctx, nil
		},
	)

	ctx.After(
		func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
			return nil, nil
		},
	)

	ctx.Given(`^the header is empty$`, test.theHeaderIsEmpty)
	ctx.Given(`^the header contains the key "([^"]*)" with "([^"]*)"$`, test.theHeaderContainsTheKeyWith)
	ctx.Given(`^the "([^"]*)" exists$`, test.theExists)
	ctx.When(`^I call "([^"]*)" "([^"]*)"$`, test.iCall)
	ctx.When(`^I call "([^"]*)" "([^"]*)" with body$`, test.iCallWithBody)
	ctx.Then(`^the status returned should be (\d+)$`, test.theStatusReturnedShouldBe)
	ctx.Then(
		`^the response should contain the field "([^"]*)" equal to "([^"]*)"$`,
		test.theResponseShouldContainTheFieldEqualTo,
	)
	ctx.Then(
		`^the response should contain the field "([^"]*)" list with (\d+) values$`,
		test.theResponseShouldContainTheFieldListWithValues,
	)
	ctx.Then(`^the db should contain the "([^"]*)" with the id "([^"]*)" colum "([^"]*)" equal to "([^"]*)"$`,
		test.theDbShouldContainTheWithTheIdColumEqualTo,
	)
	ctx.Then(`^the db should contain the "([^"]*)" with the "([^"]*)" column value "([^"]*)" colum "([^"]*)" equal to "([^"]*)"$`,
		test.theDbShouldContainTheWithTheColumnColumEqualTo,
	)
	ctx.Then(`^the db should contain the "([^"]*)" with the id "([^"]*)" array colum "([^"]*)" length is (\d+)$`,
		test.theDbShouldContainTheWithTheIdArrayColumLengthIs,
	)
	ctx.Step(`^the db should contain (\d+) objects in the "([^"]*)" table$`, test.theDbShouldContainObjectsInTheTable)
	ctx.Then(
		`^the db should contain the "([^"]*)" with the colum "([^"]*)" equal to "([^"]*)" for the following where parameters$`,
		test.theDbShouldContainTheWithTheColumEqualToForTheFollowingWhereParameters,
	)
	ctx.Then(
		`^the sqs queue "([^"]*)" should have (\d+) messages published$`,
		test.theSqsQueueShouldHaveMessagesPublished,
	)
	ctx.Then(
		`^the sqs queue "([^"]*)" should have a message published with "([^"]*)" field equal to "([^"]*)"$`,
		test.theSqsQueueShouldHaveAMessagePublishedWithFieldEqualTo,
	)
}

var serverInit sync.Once

func (t *testUtils) startApp() {
	serverInit.Do(
		func() {
			go func() {
				err := godotenv.Load("./env/.env.test")
				if err != nil {
					err := godotenv.Load("../../../env/.env.test")
					if err != nil {
						panic(err)
					}
				}
				server, err := injections.Wire().InitializeServer()
				if err != nil {
					panic(err)
				}

				err = server.Serve()
				if err != nil {
					panic(err)
				}
			}()
		},
	)
}

func (t *testUtils) reset() {
	t.headers = make(map[string]string)
	t.sqs.Reset()

	err := mocks.ClearDB(t.db)
	for err != nil {
		err = mocks.ClearDB(t.db)
	}

	//err = mock.ClearRedis(t.redis)
	//if err != nil {
	//	panic(err)
	//}
}

func (t *testUtils) theHeaderIsEmpty() error {
	t.headers = make(map[string]string)

	return nil
}

func (t *testUtils) theHeaderContainsTheKeyWith(key, value string) error {
	t.headers[key] = value

	return nil
}

func (t *testUtils) theExists(entity string, stringJson *godog.DocString) error {
	if entity == "user" {
		var payer models.User
		err := json.Unmarshal([]byte(stringJson.Content), &payer)
		if err != nil {
			return err
		}

		result := t.db.Create(&payer)
		if result.Error != nil {
			return result.Error
		}
	}
	if entity == "note" {
		var bill models.Note
		err := json.Unmarshal([]byte(stringJson.Content), &bill)
		if err != nil {
			return err
		}

		result := t.db.Create(&bill)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func (t *testUtils) iCall(method, url string) error {
	resp, err := t.execute(method, url, nil, t.headers, nil)
	if err != nil {
		return err
	}

	var respMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respMap)
	if err != nil {
		return err
	}

	t.response = &response{
		status: resp.StatusCode,
		body:   respMap,
	}

	return nil
}

func (t *testUtils) iCallWithBody(method, url string, request *godog.DocString) error {
	var req map[string]interface{}
	err := json.Unmarshal([]byte(request.Content), &req)
	if err != nil {
		return err
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := t.execute(method, url, reqByte, t.headers, nil)
	if err != nil {
		return err
	}

	var respMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respMap)
	if err != nil {
		return err
	}

	t.response = &response{
		status: resp.StatusCode,
		body:   respMap,
	}

	return nil
}

func (t *testUtils) theStatusReturnedShouldBe(status int) error {
	return assertEqual(t.response.status, status)
}

func (t *testUtils) theResponseShouldContainTheFieldEqualTo(dotSeparatedField, value string) error {
	field := getFieldValue(t.response.body, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	} else {
		return assertEqual(field, value)
	}
}

func (t *testUtils) theResponseShouldContainTheFieldListWithValues(dotSeparatedField string, quantity int) error {
	field := getFieldValue(t.response.body, dotSeparatedField)

	if ok := field.([]any); ok != nil {
		return assertEqual(len(field.([]any)), quantity)
	}
	return errors.New("field is not a list")
}

func (t *testUtils) theDbShouldContainTheWithTheIdColumEqualTo(
	entityName, id, dotSeparatedField, value string,
) error {
	var entity any
	if entityName == "user" {
		entity = &models.User{Id: id}
		result := t.db.First(entity)
		if result.Error != nil {
			return result.Error
		}
	}
	if entityName == "note" {
		entity = &models.Note{Id: id}
		result := t.db.First(entity)
		if result.Error != nil {
			return result.Error
		}
	}

	field := getFieldValue(entity, dotSeparatedField)

	switch field.(type) {
	case string, int:
		{
			return assertEqual(field, value)
		}
	default:
		jsonByte, err := json.Marshal(field)
		if err != nil {
			return err
		}

		jsonString := strings.Replace(string(jsonByte), "\"", "'", -1)

		return assertEqual(jsonString, value)
	}
}

func (t *testUtils) theDbShouldContainTheWithTheColumnColumEqualTo(
	entityName, key, keyValue, dotSeparatedField, value string,
) error {
	var entity any
	if entityName == "user" {
		entity = &models.User{}

		var result *gorm.DB
		switch key {
		case "name":
			result = t.db.Where("name = ?", keyValue).First(entity)
		default:
			return errors.New("unknown key")
		}

		if result.Error != nil {
			return result.Error
		}
	}
	if entityName == "note" {
		entity = &models.Note{}

		var result *gorm.DB
		switch key {
		case "title":
			result = t.db.Where("title = ?", keyValue).First(entity)
		default:
			return errors.New("unknown key")
		}

		if result.Error != nil {
			return result.Error
		}
	}

	entityJson, err := json.Marshal(entity)
	if err != nil {
		return err
	}

	var entityMap map[string]any
	if err := json.Unmarshal(entityJson, &entityMap); err != nil {
		return err
	}

	field := getFieldValue(entityMap, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	}

	switch field.(type) {
	case string, int:
		{
			return assertEqual(field, value)
		}
	default:
		jsonByte, err := json.Marshal(field)
		if err != nil {
			return err
		}

		jsonString := strings.Replace(string(jsonByte), "\"", "'", -1)

		return assertEqual(jsonString, value)
	}

}

func (t *testUtils) theDbShouldContainTheWithTheIdArrayColumLengthIs(
	entityName, id, dotSeparatedField string,
	length int,
) error {
	var entity any
	if entityName == "notes" {
		entity = &models.Note{UserId: id}
		result := t.db.First(entity)
		if result.Error != nil {
			return result.Error
		}
	}

	field := getFieldValue(entity, dotSeparatedField)

	return assertEqual(len(field.([]interface{})), length)
}

func (t *testUtils) theDbShouldContainObjectsInTheTable(
	quantity int,
	entityName string,
) error {
	if entityName == "users" {
		var entity []*models.User
		result := t.db.Find(&entity)
		if result.Error != nil {
			return result.Error
		}

		return assertEqual(len(entity), quantity)
	}
	if entityName == "notes" {
		var entity []*models.Note
		result := t.db.Find(&entity)
		if result.Error != nil {
			return result.Error
		}

		return assertEqual(len(entity), quantity)
	}
	return errors.New("unknown table")
}

func (t *testUtils) theDbShouldContainTheWithTheColumEqualToForTheFollowingWhereParameters(
	entityName, dotSeparatedField, value string,
	stringJson *godog.DocString,
) error {
	var queryParams interface{}
	_ = json.Unmarshal([]byte(stringJson.Content), &queryParams)
	var entity any

	if entityName == "user" {
		entity = &models.User{}
		result := t.db.Where(queryParams).First(&entity)
		if result.Error != nil {
			return result.Error
		}
	}

	return assertEqual(getFieldValue(entity, dotSeparatedField), value)
}

func (t *testUtils) theSqsQueueShouldHaveMessagesPublished(queue string, quantity int) error {
	err := assertEqual(len(t.sqs.GetMessages(queue)), quantity)
	if err != nil {
		return err
	}
	return nil
}

func (t *testUtils) theSqsQueueShouldHaveAMessagePublishedWithFieldEqualTo(queue, dotSeparatedField, value string) error {
	messages := t.sqs.GetMessages(queue)

	if messages == nil || len(messages) == 0 {
		return errors.New("message is nil")
	}

	var entity map[string]interface{}
	err := json.Unmarshal([]byte(*messages[0]), &entity)
	if err != nil {
		return err
	}

	field := getFieldValue(entity, dotSeparatedField)

	if value == "nil" {
		return assertNull(field)
	} else if value == "not nil" {
		return assertNotNull(field)
	}

	switch field.(type) {
	case string, int:
		{
			return assertEqual(field, value)
		}
	default:
		jsonByte, err := json.Marshal(field)
		if err != nil {
			return err
		}

		jsonString := strings.Replace(string(jsonByte), "\"", "'", -1)

		return assertEqual(jsonString, value)
	}
}

func (t *testUtils) execute(method, url string, request []byte, headers, queryParams map[string]string) (
	*http.Response,
	error,
) {
	var req *http.Request

	if request != nil {
		req, _ = http.NewRequest(method, t.uri+url, bytes.NewReader(request))
	} else {
		req, _ = http.NewRequest(method, t.uri+url, nil)
	}

	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if queryParams != nil {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	return t.client.Do(req)
}

func getFieldValue(object interface{}, dotSeparatedField string) interface{} {
	var field any

	var objectMap map[string]interface{}
	objectJson, _ := json.Marshal(object)
	if err := json.Unmarshal(objectJson, &objectMap); err != nil {
		var listMap []interface{}
		if err := json.Unmarshal(objectJson, &listMap); err != nil {
			panic(err)
		}
		field = listMap
	} else {
		field = objectMap
	}

	fields := strings.Split(dotSeparatedField, ".")

	for _, currentField := range fields {
		if i, err := strconv.Atoi(currentField); err == nil {
			field = field.([]interface{})[i]
		} else {
			field = field.(map[string]interface{})[currentField]
		}
	}

	switch v := field.(type) {
	case int:
		field = strconv.Itoa(v)
	case float64:
		field = strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		field = strconv.FormatBool(v)
	}

	return field
}

func assertNull(val1 interface{}) error {
	if val1 == nil {
		return nil
	}
	val1String, _ := json.Marshal(val1)
	return errors.New(string(val1String) + " should be nil")
}

func assertNotNull(val1 interface{}) error {
	if val1 != nil {
		return nil
	}
	return errors.New("value should not be nil")
}

func assertEqual(current, expected interface{}) error {
	if current == expected {
		return nil
	}
	currentString, _ := json.Marshal(current)
	expectedString, _ := json.Marshal(expected)
	return errors.New(string(currentString) + " should be equal to " + string(expectedString))
}
