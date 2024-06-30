package main

import "testing"

func TestPerson(t *testing.T) {
	id := createPerson(Person{Name: "mando", Age: 34})

	p, ok := datalayer[1]

	if id == 1 && ok && p.Name == "mando" && p.Age == 34 {
		t.Log("person created successfully")
		return
	}

	t.Error("create person failed", id, p, ok, datalayer)
}

func TestReadPerson(t *testing.T) {
	p, err := readPerson(1)

	if err != nil {
		t.Error("read person error", err)
	}

	if p.Name != "mando" || p.Age != 34 {
		t.Error("invalid name or age")
	}

	t.Log("read person successful", p)
}

// func TestUpdatePerson(t *testing.T) {
// 	p, err := updatePerson(1, Person{Name: "kp", Age: 20})

// 	if err != nil {
// 		t.Error("update person err", err)

// 	}

// 	if p.Name != "kp" && p.Age != 20 {
// 		t.Error("Name or age mismatch")
// 	}

// 	t.Log("update person successfull", p)

// }

func TestDeletePerson(t *testing.T){
	err := deletePerson(1)

	if err != nil {
		t.Error("delete person error", err)
	}
	t.Log("delete person successful")
}
