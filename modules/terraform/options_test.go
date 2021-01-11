package terraform

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/require"
)

func TestGetVariablesFromVarFilesAsString(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		number_type = 2
		boolean_type = true
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	stringVal := GetVariableAsStringFromVarFile(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "aws_region")

	boolString := GetVariableAsStringFromVarFile(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "boolean_type")

	numString := GetVariableAsStringFromVarFile(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "number_type")

	require.Equal(t, "us-east-2", stringVal)
	require.Equal(t, "true", boolString)
	require.Equal(t, "2", numString)

}

func TestGetVariablesFromVarFilesAsStringKeyDoesNotExist(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsStringFromVarFileE(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetVariableAsMapFromVarFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	expected := make(map[string]string)
	expected["foo"] = "bar"

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	val := GetVariableAsMapFromVarFile(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "tags")

	require.Equal(t, expected, val)
}

func TestGetVariableAsMapFromVarFileNotMap(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsMapFromVarFileE(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "aws_region")

	require.Error(t, err)
}

func TestGetVariableAsMapFromVarFileKeyDoesNotExist(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsMapFromVarFileE(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetVariableAsListFromVarFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	expected := []string{"item1"}

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	val := GetVariableAsListFromVarFile(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "list")

	require.Equal(t, expected, val)
}

func TestGetVariableAsListNotList(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsListFromVarFileE(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "tags")

	require.Error(t, err)
}

func TestGetVariableAsListKeyDoesNotExist(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())

	testHcl := []byte(`
		aws_region     = "us-east-2"
		aws_account_id = "111111111111"
		tags = {
			foo = "bar"
		}
		list = ["item1"]`)

	WriteFile(t, randomFileName, testHcl)
	defer os.Remove(randomFileName)

	_, err := GetVariableAsListFromVarFileE(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName, "badkey")

	require.Error(t, err)
}

func TestGetAllVariablesFromVarFileENotInVarFiles(t *testing.T) {
	_, err := GetAllVariablesFromVarFileE(t, &Options{
		VarFiles: []string{"filea"},
	}, "fileb")

	//require.Equal(t, "open thisdoesntexist: no such file or directory", err.Error())
	require.Error(t, err)
}

func TestGetAllVariablesFromVarFileEFileDoesNotExist(t *testing.T) {
	_, err := GetAllVariablesFromVarFileE(t, &Options{
		VarFiles: []string{"filea"},
	}, "filea")

	require.Equal(t, "open filea: no such file or directory", err.Error())
}

func TestGetAllVariablesFromVarFileBadFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	testHcl := []byte(`
		thiswillnotwork`)

	err := ioutil.WriteFile(randomFileName, testHcl, 0644)

	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

	defer os.Remove(randomFileName)

	_, err = GetAllVariablesFromVarFileE(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName)

	if err == nil {
		t.FailNow()
	}

	// HCL library could change their error string, so we are only testing the error string contains what we add to it
	require.Regexp(t, fmt.Sprintf("^%s - ", randomFileName), err.Error())

}

func TestGetAllVariablesFromVarFile(t *testing.T) {
	randomFileName := fmt.Sprintf("./%s.tfvars", random.UniqueId())
	testHcl := []byte(`
	aws_region     = "us-east-2"
	`)

	err := ioutil.WriteFile(randomFileName, testHcl, 0644)

	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}

	defer os.Remove(randomFileName)

	val, err := GetAllVariablesFromVarFileE(t, &Options{
		VarFiles: []string{randomFileName},
	}, randomFileName)

	if err != nil {
		t.FailNow()
	}

	expected := make(map[string]interface{})
	expected["aws_region"] = "us-east-2"

	require.Equal(t, expected, val)

}

// Helper function to write a file to the filesystem
// Will immediately fail the test if it could not write the file
func WriteFile(t *testing.T, fileName string, bytes []byte) {
	err := ioutil.WriteFile(fileName, bytes, 0644)

	if err != nil {
		fmt.Println(err.Error())
		t.FailNow()
	}
}
