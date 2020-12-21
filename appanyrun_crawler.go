package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// ws endpoint
var endpoints = [...]string{
	"wss://app.any.run/sockjs/158/r2jz998p/websocket",
	"wss://app.any.run/sockjs/479/eokzh54x/websocket",
	"wss://app.any.run/sockjs/937/thitlatz/websocket",
	"wss://app.any.run/sockjs/222/3_5u81il/websocket",
}

const (
	// handshake message
	connectMsg = `["{\"msg\":\"connect\",\"version\":\"1\",\"support\":[\"1\",\"pre2\",\"pre1\"]}"]`

	// public tasks
	//publicTasksUrlFormat        = `["{\"msg\":\"sub\",\"id\":\"vTWcZmngJ49BLcmsr\",\"name\":\"publicTasks\",\"params\":[%d,%d,{\"isPublic\":true,\"hash\":\"\",\"runtype\":[],\"verdict\":[],\"ext\":[],\"tag\":\"%s\",\"significant\":false,\"ip\":\"\",\"fileHash\":\"\",\"mitreId\":\"\",\"sid\":0,\"skip\":0}]}"]`
	//publicTasksCounterUrlFormat = `["{\"msg\":\"method\",\"method\":\"publicTasksCounter\",\"params\":[{\"isPublic\":true,\"hash\":\"\",\"runtype\":[],\"verdict\":[],\"ext\":[],\"tag\":\"%s\",\"significant\":false,\"ip\":\"\",\"fileHash\":\"\",\"mitreId\":\"\",\"sid\":0,\"skip\":0}],\"id\":\"4\"}"]`

	// PE EXE public tasks
	publicTasksUrlFormat        = `["{\"msg\":\"sub\",\"id\":\"qDA2CKe3Km4N9MPAE\",\"name\":\"publicTasks\",\"params\":[%d,%d,{\"isPublic\":true,\"hash\":\"\",\"runtype\":[],\"verdict\":[],\"ext\":[\"0\"],\"ip\":\"\",\"domain\":\"\",\"fileHash\":\"\",\"mitreId\":\"\",\"sid\":0,\"significant\":false,\"tag\":\"%s\",\"skip\":0}]}"]`
	publicTasksCounterUrlFormat = `["{\"msg\":\"method\",\"method\":\"publicTasksCounter\",\"params\":[{\"isPublic\":true,\"hash\":\"\",\"runtype\":[],\"verdict\":[],\"ext\":[\"0\"],\"ip\":\"\",\"domain\":\"\",\"fileHash\":\"\",\"mitreId\":\"\",\"sid\":0,\"significant\":false,\"tag\":\"%s\",\"skip\":0}],\"id\":\"7\"}"]`
	publicTasksDoneMsg          = `{"msg":"ready","subs":["qDA2CKe3Km4N9MPAE"]}`

	// File, PE EXE public tasks
	//publicTasksUrlFormat        = `["{\"msg\":\"sub\",\"id\":\"2ZYiGdXjJxY4QBvtb\",\"name\":\"publicTasks\",\"params\":[%d,%d,{\"isPublic\":true,\"hash\":\"\",\"runtype\":[\"1\"],\"verdict\":[],\"ext\":[\"0\"],\"ip\":\"\",\"domain\":\"\",\"fileHash\":\"\",\"mitreId\":\"\",\"sid\":0,\"significant\":false,\"tag\":\"%s\",\"skip\":0}]}"]`
	//publicTasksCounterUrlFormat = `["{\"msg\":\"method\",\"method\":\"publicTasksCounter\",\"params\":[{\"isPublic\":true,\"hash\":\"\",\"runtype\":[\"1\"],\"verdict\":[],\"ext\":[\"0\"],\"ip\":\"\",\"domain\":\"\",\"fileHash\":\"\",\"mitreId\":\"\",\"sid\":0,\"significant\":false,\"tag\":\"%s\",\"skip\":0}],\"id\":\"7\"}"]`
	//publicTasksDoneMsg          = `{"msg":"ready","subs":["2ZYiGdXjJxY4QBvtb"]}`

	// process tree
	processUrlFormat = `["{\"msg\":\"sub\",\"id\":\"ojEf2kD8Qo8Nt8aCg\",\"name\":\"process\",\"params\":[{\"taskID\":{\"$type\":\"oid\",\"$value\":\"%s\"},\"status\":100,\"important\":true}]}"]`
	processDoneMsg   = `{"msg":"ready","subs":["ojEf2kD8Qo8Nt8aCg"]}`

	// Mitre ATT&CK Mapping
	allIncidentsUrlFormat = `["{\"msg\":\"sub\",\"id\":\"xhR3rXWu4M8X6xFow\",\"name\":\"allIncidents\",\"params\":[{\"$type\":\"oid\",\"$value\":\"%s\"}]}"]`
	allIncidentsDoneMsg   = `{"msg":"ready","subs":["xhR3rXWu4M8X6xFow"]}`

	taskExistsUrlFormat = `["{\"msg\":\"sub\",\"id\":\"L6La59ezwZEf9qP2F\",\"name\":\"taskexists\",\"params\":[\"%s\"]}"]`
	taskExistsDoneMsg   = `["{\"msg\":\"ready\",\"subs\":[\"L6La59ezwZEf9qP2F\"]}"]`

	singleTaskUrlFormat = `["{\"msg\":\"sub\",\"id\":\"mkdKdJqprjPj98Z2e\",\"name\":\"singleTask\",\"params\":[{\"$type\":\"oid\",\"$value\":\"%s\"},true]}"]`
	singleTaskDoneMsg   = `["{\"msg\":\"ready\",\"subs\":[\"mkdKdJqprjPj98Z2e\"]}"]`
)

func getPublicTasksUrl(tag string, numOfEntries, index int) string {
	return fmt.Sprintf(publicTasksUrlFormat, numOfEntries, index, tag)
}

func getPublicTasksCounterUrl(tag string) string {
	return fmt.Sprintf(publicTasksCounterUrlFormat, tag)
}

func getProcessUrl(taskId string) string {
	return fmt.Sprintf(processUrlFormat, taskId)
}

func getAllIncidentsUrl(taskId string) string {
	return fmt.Sprintf(allIncidentsUrlFormat, taskId)
}

func getTaskExistsUrl(taskUuid string) string {
	return fmt.Sprintf(taskExistsUrlFormat, taskUuid)
}

func getSingleTaskUrl(taskId string) string {
	return fmt.Sprintf(singleTaskUrlFormat, taskId)
}

func sendAll(conn *websocket.Conn, msg string) error {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return err
	}
	return nil
}

func readAll(conn *websocket.Conn) (string, error) {
	_, bytes, err := conn.ReadMessage()
	if err != nil {
		return "", err
	}
	if bytes[0] == 'a' && bytes[1] == '[' { // if the message in format a[payload]
		msg, err := strconv.Unquote(string(bytes[2 : len(bytes)-1]))
		if err != nil {
			return "", err
		}
		return msg, nil
	}
	return string(bytes), nil
}

func dumpProcessTree(conn *websocket.Conn, taskId string) ([]*Process, error) {
	processes := make([]*Process, 0)
	if err := sendAll(conn, getProcessUrl(taskId)); err != nil {
		return nil, err
	}
	for { // receive all
		rawProc := new(RawProcess)
		msg, err := readAll(conn)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(msg), &rawProc); err != nil {
			log.Println(msg)
			return nil, err
		}
		if rawProc.Fields.Pid == 0 && msg == processDoneMsg {
			break
		}
		processes = append(processes, NewProcess(rawProc))
	}
	return processes, nil
}

func dumpAllIncidents(conn *websocket.Conn, taskId string) ([]*Incident, error) {
	incidents := make([]*Incident, 0)
	if err := sendAll(conn, getAllIncidentsUrl(taskId)); err != nil {
		return nil, err
	}
	for { // receive all
		rawIncident := new(RawIncident)
		msg, err := readAll(conn)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(msg), &rawIncident); err != nil {
			return nil, err
		}
		if rawIncident.Collection == "" && msg == allIncidentsDoneMsg {
			break
		}
		incidents = append(incidents, NewIncident(rawIncident))
	}
	return incidents, nil
}

func dumpToFile(fileName string, bytes []byte) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func dumpTask(conn *websocket.Conn, malwareTag string, task *Task) error {
	processes, err := dumpProcessTree(conn, task.ID)
	if err != nil {
		return err
	}
	incidents, err := dumpAllIncidents(conn, task.ID)
	if err != nil {
		return err
	}
	mainObject := task.Fields.Public.Objects.MainObject
	processData := &ProcessData{
		Name:      mainObject.Names.Basename,
		Md5:       mainObject.Hashes.Md5,
		UUID:      task.Fields.UUID,
		Processes: processes,
		Incidents: incidents,
	}
	bytes, err := json.MarshalIndent(processData, "", " ")
	if err != nil {
		return err
	}
	taskFileName := fmt.Sprintf("%s/%s.json", malwareTag, getTaskUrl(task))
	if err := dumpToFile(taskFileName, bytes); err != nil {
		return err
	}
	return nil
}

func NewAppAnyClient() *websocket.Conn {
	reqHeader := make(http.Header)
	reqHeader.Add("Host", "app.any.run")
	reqHeader.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36")
	reqHeader.Add("Origin", "https://app.any.run")

	//rand.Seed(time.Now().Unix())
	conn, _, err := websocket.DefaultDialer.Dial(endpoints[rand.Intn(len(endpoints))], reqHeader)
	if err != nil {
		log.Fatal(err)
	}
	conn.ReadMessage()
	conn.ReadMessage()

	// connect
	sendAll(conn, connectMsg)
	msg, err := readAll(conn)
	if err != nil {
		log.Fatal(err)
	}
	if !strings.Contains(msg, "connected") {
		log.Fatal("connection to app.any failed")
	}
	return conn
}

//GetNumOfTasks
func countTasksByTag(conn *websocket.Conn, malwareTag string) (int, error) {
	// count public tasks by tag
	var countResult Result
	if err := sendAll(conn, getPublicTasksCounterUrl(malwareTag)); err != nil {
		return 0, err
	}
	// a["{\"msg\":\"updated\",\"methods\":[\"4\"]}"]
	conn.ReadMessage()
	msg, err := readAll(conn)
	if err != nil {
		return 0, err
	}
	if err := json.Unmarshal([]byte(msg), &countResult); err != nil {
		return 0, err
	}
	return countResult.Result.Count, nil
}

//GetTasks
func getTasksByTag(conn *websocket.Conn, malwareTag string, numOfEntries, index int) ([]*Task, error) {
	tasks := make([]*Task, 0)
	// get public tasks
	if err := sendAll(conn, getPublicTasksUrl(malwareTag, numOfEntries, index)); err != nil {
		return nil, err
	}
	for { // receive all
		task := new(Task)
		msg, err := readAll(conn)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(msg), &task); err != nil {
			log.Println(msg)
			return nil, err
		}
		if task.Collection == "" && msg == publicTasksDoneMsg {
			break
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func getTaskIdentity(task *Task) string {
	mainObject := task.Fields.Public.Objects.MainObject
	format := "name: %s, MD5: %s, TaskID: %s"
	switch mainObject.Type {
	case "file":
		return fmt.Sprintf(format, mainObject.Names.Basename, mainObject.Hashes.Md5, mainObject.Task.Value)
	case "url":
		return fmt.Sprintf(format, mainObject.Names.URL, mainObject.Hashes.Md5, mainObject.Task.Value)
	}
	return "unknown"
}

func getTaskUrl(task *Task) string {
	return task.Fields.UUID
}

func crawlTasks(malwareTag string, taskIndex, numOfTasks int) {
	conn := NewAppAnyClient()

	taskCount, err := countTasksByTag(conn, malwareTag)
	if err != nil {
		conn.Close()
		log.Println("cannot count tasks, ", err)
		return
	}
	log.Printf("Number of tasks for %s: %d\n", malwareTag, taskCount)
	conn.Close()
	if numOfTasks <= 0 || numOfTasks > taskCount {
		numOfTasks = taskCount
	}
	log.Printf("Start crawling %d tasks\n", numOfTasks)

	var counter = taskIndex
	for i := taskIndex; i < numOfTasks; i += 50 {
		conn := NewAppAnyClient()
		tasks, err := getTasksByTag(conn, malwareTag, 50, i)
		if err != nil {
			conn.Close()
			log.Println("cannot get tasks, ", err)
			return
		}
		os.Mkdir(malwareTag, os.ModePerm)
		conn.Close()

		for _, task := range tasks {
			log.Println(counter, getTaskIdentity(task))
			counter++
			conn := NewAppAnyClient()
			if err := dumpTask(conn, malwareTag, task); err != nil {
				log.Println("cannot dump task, ", err)
				conn.Close()
				return
			}
			conn.Close()
		}
	}
}

type TaskExistsResult struct {
	Msg        string `json:"msg"`
	Collection string `json:"collection"`
	ID         string `json:"id"`
	Fields     struct {
		TaskID       string `json:"taskId"`
		TaskObjectID struct {
			Type  string `json:"$type"`
			Value string `json:"$value"`
		} `json:"taskObjectId"`
	} `json:"fields"`
}

func crawlTaskByUUID(outDirPath, taskUuid string) error {
	conn := NewAppAnyClient()
	// check existence and get internal id
	var result TaskExistsResult
	if err := sendAll(conn, getTaskExistsUrl(taskUuid)); err != nil {
		return fmt.Errorf("in sendAll: %s", err)
	}
	msg, err := readAll(conn)
	if err != nil {
		return fmt.Errorf("in readAll: %s", err)
	}
	if err := json.Unmarshal([]byte(msg), &result); err != nil {
		return fmt.Errorf("in Unmarshal: %s", err)
	}
	conn.ReadMessage()

	// get process tree and incidents
	taskId := result.Fields.TaskObjectID.Value
	processes, err := dumpProcessTree(conn, taskId)
	if err != nil {
		return err
	}
	incidents, err := dumpAllIncidents(conn, taskId)
	if err != nil {
		return err
	}
	// task information
	var taskInfo *Task
	if err := sendAll(conn, getSingleTaskUrl(taskId)); err != nil {
		return fmt.Errorf("in sendAll: %s", err)
	}
	msg, err = readAll(conn)
	if err != nil {
		return fmt.Errorf("in readAll: %s", err)
	}
	if err := json.Unmarshal([]byte(msg), &taskInfo); err != nil {
		return fmt.Errorf("in Unmarshal: %s", err)
	}
	// save
	mainObject := taskInfo.Fields.Public.Objects.MainObject
	processData := &ProcessData{
		Name:      mainObject.Names.Basename,
		Md5:       mainObject.Hashes.Md5,
		UUID:      taskInfo.Fields.UUID,
		Processes: processes,
		Incidents: incidents,
	}
	bytes, err := json.MarshalIndent(processData, "", " ")
	if err != nil {
		return err
	}
	if _, err := os.Stat(outDirPath); os.IsNotExist(err) {
		if err = os.Mkdir(outDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create dir for saving: %s", err)
		}
	}
	taskFileName := fmt.Sprintf("%s/%s.json", outDirPath, getTaskUrl(taskInfo))
	if err := dumpToFile(taskFileName, bytes); err != nil {
		return err
	}
	return nil
}
