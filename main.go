package main

import (
	"employee-management/dao"
	"employee-management/model"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

var ead = dao.EmployeeDAO{}

func apiResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"hello world!"}`))
}

func CreateNewEmployee(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "POST" {
		respondWithError(w, http.StatusBadRequest, "Invalid method")
	}

	var employee model.Employee

	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	if err := ead.Insert(employee); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	} else {
		respondWithJson(w, http.StatusAccepted, map[string]string{
			"message": "Record inserted successfully",
		})
	}
}

func getEmployeesById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "GET" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}

	empId := strings.Split(r.URL.Path, "/")[2]

	employees, err := ead.FindByEmpId(empId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, employees)
}

func deleteEmployeeById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "DELETE" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}
	var reqBody map[string]string

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	empId := reqBody["employee_id"]

	err := ead.DeleteEmployee(empId)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Record deleted successfully",
	})
}

func updateEmployeeById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	if r.Method != "PUT" {
		respondWithError(w, http.StatusBadRequest, "Method not allowed")
	}
	var emp model.Employee
	err := json.NewDecoder(r.Body).Decode(&emp)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
	}

	empId := emp.EmployeeId

	_ = ead.UpdateEmployee(empId, emp)

	respondWithJson(w, http.StatusOK, map[string]string{
		"message": "Record updated successfully",
	})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func init() {
	ead.Server = "mongodb://localhost:27017/"
	ead.Database = "EmployeeDB"
	ead.Collection = "Employee"

	ead.Connect()
}

func main() {
	http.HandleFunc("/", apiResponse)
	http.HandleFunc("/employee/", getEmployeesById)
	http.HandleFunc("/add-employee", CreateNewEmployee)
	http.HandleFunc("/delete-employee", deleteEmployeeById)
	http.HandleFunc("/update-employee", updateEmployeeById)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
