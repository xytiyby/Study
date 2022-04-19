package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
)

type Student struct {
	ObjectType string `json:"docType"`
	Name       string `json:"Name"`      //姓名
	EntityID   string `json:"EntityID"`  //身份证号
	StudentID  string `json:"StudentID"` //学号
	Password   string `json:"Password"`  //密码
	Major      string `json:"Major"`     //专业
}
type StudentChaincode struct {
}

func (t *StudentChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}
func (t *StudentChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// 获取用户意图
	fun, args := stub.GetFunctionAndParameters()
	if fun == "addStudent" {
		return t.addStudent(stub, args) // 添加信息
	}
	else if fun == "queryByEntityIdAndName" {
		return t.queryByEntityIdAndName(stub, args) // 根据证书编号及姓名查询信息
	}

return shim.Error("指定的函数名称错误")

}
const DOC_TYPE = "stuObj"
func PutStudent(stub shim.ChaincodeStubInterface, stu Student) ([]byte, bool) {

	stu.ObjectType = DOC_TYPE

	b, err := json.Marshal(stu)
	if err != nil {
		return nil, false
	}

	// 保存edu状态
	err = stub.PutState(stu.EntityID, b)
	if err != nil {
		return nil, false
	}

	return b, true
}
func GetEduInfo(stub shim.ChaincodeStubInterface, entityID string) (Student, bool) {
	var stu Student
	// 根据身份证号码查询信息状态
	b, err := stub.GetState(entityID)
	if err != nil {
		return stu, false
	}

	if b == nil {
		return stu, false
	}

	// 对查询到的状态进行反序列化
	err = json.Unmarshal(b, &stu)
	if err != nil {
		return stu, false
	}

	// 返回结果
	return stu, true
}
func (t *StudentChaincode) addStudent(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}

	var stu Student
	err := json.Unmarshal([]byte(args[0]), &stu)
	if err != nil {
		return shim.Error("反序列化信息时发生错误")
	}

	// 查重: 身份证号码必须唯一
	_, exist := GetEduInfo(stub, stu.EntityID)
	if exist {
		return shim.Error("要添加的身份证号码已存在")
	}

	_, bl := PutStudent(stub, stu)
	if !bl {
		return shim.Error("保存信息时发生错误")
	}

	err = stub.SetEvent(args[1], []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("信息添加成功"))
}
// 根据指定的查询字符串实现富查询
func getEduByQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}

		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		bArrayMemberAlreadyWritten = true
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil

}
func (t *StudentChaincode) queryByEntityIdAndName(stub shim.ChaincodeStubInterface, args []string) peer.Response {

	if len(args) != 2 {
		return shim.Error("给定的参数个数不符合要求")
	}
	studentID := args[0]
	name := args[1]

	// 拼装CouchDB所需要的查询字符串(是标准的一个JSON串)
	// queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"eduObj\", \"CertNo\":\"%s\"}}", CertNo)
	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"%s\", \"StudentID\":\"%s\", \"Name\":\"%s\"}}", DOC_TYPE, studentID, name)

	// 查询数据
	result, err := getEduByQueryString(stub, queryString)
	if err != nil {
		return shim.Error("根据证书编号及姓名查询信息时发生错误")
	}
	if result == nil {
		return shim.Error("根据指定的证书编号及姓名没有查询到相关的信息")
	}
	return shim.Success(result)
}