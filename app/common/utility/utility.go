package utility

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gorm.io/datatypes"
	"hintword.com/api/app/common/constants"
)

var lastPrimaryId int64

func UniqueId(index ...int) int64 {
	random := rand.Intn(99)
	ok := random
	if random < 10 {
		ok = ok + 10
	}
	unique := strconv.Itoa(int(time.Now().UnixMicro())) + strconv.Itoa(ok)
	i, _ := strconv.ParseInt(unique, 10, 64)

	if len(index) > 0 && index[0] != 0 {
		idString := fmt.Sprintf("%v", i)
		withoutSuffix := idString[:len(idString)-len(strconv.Itoa(index[0]))]
		final := fmt.Sprintf("%v%v", withoutSuffix, index[0])
		return int64(StringToInt(final))
	}

	if lastPrimaryId == i {
		time.Sleep(time.Second)
		unique = strconv.Itoa(int(time.Now().UnixMicro())) + strconv.Itoa(ok)
		j, _ := strconv.ParseInt(unique, 10, 64)
		return j
	}
	lastPrimaryId = i
	return i
}

func MinifyXml(xmlData []byte) (*bytes.Buffer, error) {
	minifiedXML := bytes.NewBuffer(nil)

	// Decode the XML data into a structured format
	decoder := xml.NewDecoder(bytes.NewReader(xmlData))
	encoder := xml.NewEncoder(minifiedXML)
	encoder.Indent("", "")

	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			// Write the start element to the minified XML buffer
			err = encoder.EncodeToken(se)
			if err != nil {
				fmt.Println("Error encoding XML token:", err)
				return minifiedXML, err
			}
		case xml.EndElement:
			// Write the end element to the minified XML buffer
			err = encoder.EncodeToken(se)
			if err != nil {
				fmt.Println("Error encoding XML token:", err)
				return minifiedXML, err
			}
		case xml.CharData:
			// Write the character data (content) to the minified XML buffer
			err = encoder.EncodeToken(xml.CharData(bytes.TrimSpace(se)))
			if err != nil {
				fmt.Println("Error encoding XML token:", err)
				return minifiedXML, err
			}
		}

	}

	err := encoder.Flush()
	if err != nil {
		fmt.Println("Error encoding XML token:", err)
		return minifiedXML, err
	}

	return minifiedXML, nil

}

func CrifGenderMap(value string) string {
	switch value {
	case "MALE":
		return "G01"
	case "FEMALE":
		return "G02"
	default:
		return "G03"
	}
}

func SliceStringToString(a []string) string {
	b := ""
	for _, v := range a {
		if len(b) > 0 {
			b += " "
		}
		b += v
	}

	return b
}

func PincodeFromAddress(value string) string {
	pattern := `[0-9]{6}`
	regex := regexp.MustCompile(pattern)
	return regex.FindString(value)
}

func GetMatchingStateCode(stateCodes []string, targetString string) string {

	// Extract the state code from the target string
	targetStateCode := targetString[len(targetString)-2:]

	// Check if the target state code is present in the list
	for _, code := range stateCodes {
		if code == targetStateCode {
			return code
		}
	}

	return ""
}

func StructToMap(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(data)
	typ := reflect.TypeOf(data)

	if val.Kind() != reflect.Struct {
		return result
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		result[fieldName] = field.Interface()
	}

	return result
}

func ParseDate(s string) time.Time {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}
	}

	return t
}

func GenerateCode(first string, second string, third string, fourth string) string {
	return strings.TrimSpace(fmt.Sprintf("%v%v%v%v", first, second, third, fourth))
}

func StringToInt(value string) int {
	id, _ := strconv.Atoi(value)
	/*if err != nil {
		fmt.Printf("string conversion failedL %v\n", err)
	}*/
	return id
}

func GetRunningSeries(series string, lastDigit int) (data string) {

	defer func() {
		if r := recover(); r != nil {
			data = fmt.Sprintf("%s%s", strings.Repeat("0", lastDigit-1), "1")
		}
	}()

	lastDigitsStr := series[len(series)-lastDigit:]

	// Convert the last four digits to an integer
	lastDigits, err := strconv.Atoi(lastDigitsStr)

	if err != nil {
		fmt.Println("Error converting last four digits to integer:", err)
		return ""
	}

	// Increment the last four digits
	lastDigits++

	newLastFourDigitsStr := fmt.Sprintf("%0"+strconv.Itoa(lastDigit)+"d", lastDigits)

	return newLastFourDigitsStr
}

func Unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	var list []string
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func CleanURLString(input string) string {
	// Convert to lowercase
	result := strings.ToLower(input)

	// Remove special characters
	re := regexp.MustCompile("[^a-z0-9]+")
	result = re.ReplaceAllString(result, "-")

	// Remove leading and trailing hyphens
	result = strings.Trim(result, "-")

	return result
}

func CleanString(keyword string) string {

	//keyword = strings.ReplaceAll(keyword, "  ", " ")
	//keyword = strings.ReplaceAll(keyword, "(", "")
	keyword = strings.ReplaceAll(keyword, "&lt;", "<")
	keyword = strings.ReplaceAll(keyword, "&gt;", ">")
	keyword = strings.ReplaceAll(keyword, "&quot;", `"`)
	keyword = strings.ReplaceAll(keyword, "&nbsp;", "")
	keyword = strings.ReplaceAll(keyword, " ", "")
	keyword = strings.ReplaceAll(keyword, "\"", "'")
	keyword = strings.ReplaceAll(keyword, "\n", "")

	return strings.TrimSpace(keyword)
}

func MapToStruct(inputMap map[string]interface{}, resultStruct interface{}) error {
	resultValue := reflect.ValueOf(resultStruct).Elem()

	for key, value := range inputMap {
		structField := resultValue.FieldByName(key)

		if !structField.IsValid() {
			return fmt.Errorf("field %s not found in the struct", key)
		}

		if !structField.CanSet() {
			return fmt.Errorf("cannot set value for field %s", key)
		}

		mapValue := reflect.ValueOf(value)
		if structField.Type() != mapValue.Type() {
			return fmt.Errorf("type mismatch for field %s", key)
		}

		structField.Set(mapValue)
	}

	return nil
}

func MapToJSON(inputMap map[string]interface{}) (string, error) {
	jsonData, err := json.Marshal(inputMap)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// JSONToStruct converts JSON string to a struct
func JSONToStruct(jsonData string, resultStruct interface{}) error {
	err := json.Unmarshal([]byte(jsonData), resultStruct)
	return err
}

func ContainInList(inputList []string, value string) (result bool) {
	for _, item := range inputList {
		if reflect.DeepEqual(value, item) {
			return true
		}
	}
	return false
}

func StructToJSON(inputStruct interface{}) (datatypes.JSON, error) {
	jsonData, err := json.Marshal(inputStruct)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	return jsonData, nil
}

func GetRolePriority(role string) int {
	switch role {
	case "PRIMARY":
		return 1
	case "CO_APPLICANT":
		return 2
	case "GUARANTOR":
		return 3
	default:
		return 4 // For any other roles
	}
}
func CheckValueInList(inputList []interface{}, value interface{}) (result bool) {
	for _, item := range inputList {
		if reflect.DeepEqual(value, item) {
			return true
		}
	}
	return false
}

func CalculatePercentage(x int, y int) float64 {
	if y == 0 {
		return 0
	}
	return toFixed(float64(x)/float64(y)*float64(100), 2)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func IsExpired(expirationTime time.Time) bool {
	// Compare the current time with the expiration time
	return time.Now().After(expirationTime)
}

func ConvertPlainTextToJSON(plainTextPayload string) (map[string]interface{}, error) {
	// Parse the plain text payload into a map
	values, err := parsePlainText(plainTextPayload)
	if err != nil {
		return nil, err
	}

	// Convert the map to a JSON object
	jsonPayload, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	// Unmarshal the JSON into a map[string]interface{}
	var result map[string]interface{}
	err = json.Unmarshal(jsonPayload, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func parsePlainText(plainTextPayload string) (map[string][]string, error) {
	// You may need to implement your own parsing logic based on your payload format
	// In this example, we assume a URL-encoded format similar to key1=value1&key2=value2
	// You may need to adjust this based on your actual payload format
	values, err := url.ParseQuery(plainTextPayload)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func CleanConditionS(value string) string {
	value = strings.ReplaceAll(value, " ", "_")
	return strings.ToUpper(value)
}

func CleanCountryCode(inputString, countryCode string) string {

	// Check if the input string starts with the country code
	if strings.HasPrefix(inputString, countryCode) {
		// Remove the country code from the beginning of the string
		result := inputString[len(countryCode):]
		return result
	} else {
		// If the input string doesn't start with the country code, do nothing
		fmt.Println("No country code found in the input string.")
	}

	return inputString
}

func TimestampToTime(stamp string) time.Time {
	timestampMilliseconds := int64(StringToInt(stamp))
	return time.Unix(0, timestampMilliseconds*int64(time.Millisecond))
}

func LastMonth(month int) time.Time {
	previous := time.Now().AddDate(0, month, 0).UTC()
	return time.Date(previous.Year(), previous.Month(), 1, previous.Hour(), previous.Minute(), previous.Second(),
		previous.Nanosecond(), previous.Location())
}

func ValidatePlatform(platform string) string {
	switch platform {
	case constants.SystemPartnerPortal, constants.SystemPartnerMobile:
		return constants.USER_TYPE_CHANNEL
	case constants.SystemEmployeePortal, constants.SystemEmployeeMobile:
		return constants.USER_TYPE_EMPLOYEE
	default:
		return ""
	}
}

func LoanCode(input string) string {
	words := strings.Fields(input)
	var result string

	if len(words) == 1 {
		return strings.ToUpper(words[0])
	}

	for _, word := range words {
		if len(word) > 0 {
			result += strings.ToUpper(string(word[0]))
		}
	}

	return result
}

func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
