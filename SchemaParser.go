package main
import "bufio"
import "os"
import "fmt"
import "strconv"
import "regexp"

const (
	PROP_TYPE_INVALID uint8 = iota
	PROP_TYPE_STRING = iota
	PROP_TYPE_STRING_NO_CACHE = iota
	PROP_TYPE_PROTECTEDSTRING_0 = iota
	PROP_TYPE_PROTECTEDSTRING_1 = iota
	PROP_TYPE_PROTECTEDSTRING_2 = iota
	PROP_TYPE_PROTECTEDSTRING_3 = iota
	PROP_TYPE_ENUM = iota
	PROP_TYPE_BINARYSTRING = iota
	PROP_TYPE_PBOOL = iota
	PROP_TYPE_PSINT = iota
	PROP_TYPE_PFLOAT = iota
	PROP_TYPE_PDOUBLE = iota
	PROP_TYPE_UDIM = iota
	PROP_TYPE_UDIM2 = iota
	PROP_TYPE_RAY = iota
	PROP_TYPE_FACES = iota
	PROP_TYPE_AXES = iota
	PROP_TYPE_BRICKCOLOR = iota
	PROP_TYPE_COLOR3 = iota
	PROP_TYPE_COLOR3UINT8 = iota
	PROP_TYPE_VECTOR2 = iota
	PROP_TYPE_VECTOR3_SIMPLE = iota
	PROP_TYPE_VECTOR3_COMPLICATED = iota
	PROP_TYPE_VECTOR2UINT16 = iota
	PROP_TYPE_VECTOR3UINT16 = iota
	PROP_TYPE_CFRAME_SIMPLE = iota
	PROP_TYPE_CFRAME_COMPLICATED = iota
	PROP_TYPE_INSTANCE = iota
	PROP_TYPE_TUPLE = iota
	PROP_TYPE_ARRAY = iota
	PROP_TYPE_DICTIONARY = iota
	PROP_TYPE_MAP = iota
	PROP_TYPE_CONTENT = iota
	PROP_TYPE_SYSTEMADDRESS = iota
	PROP_TYPE_NUMBERSEQUENCE = iota
	PROP_TYPE_NUMBERSEQUENCEKEYPOINT = iota
	PROP_TYPE_NUMBERRANGE = iota
	PROP_TYPE_COLORSEQUENCE = iota
	PROP_TYPE_COLORSEQUENCEKEYPOINT = iota
	PROP_TYPE_RECT2D = iota
	PROP_TYPE_PHYSICALPROPERTIES = iota
)

type StaticPropertySchema struct {
	Name string
	Type uint8
	TypeString string
	InstanceSchema *StaticInstanceSchema
}

type StaticInstanceSchema struct {
	Name string
	Properties []StaticPropertySchema
}

type StaticSchema struct {
	Instances []StaticInstanceSchema
	Properties []StaticPropertySchema
}

func parseInstSchema(filename string) ([]StaticInstanceSchema, error) {
	var countInstances uint16

	propFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	file := bufio.NewReader(propFile)
	_, err = fmt.Fscanf(file, "%d\n", &countInstances)
	if err != nil {
		return nil, err
	}
	instances := make([]StaticInstanceSchema, countInstances)

	propMatcher := regexp.MustCompile(`\s+(\d+)\s+'([a-zA-Z0-9 _]+)'\s+(\w+)\s*`)

	for i := 0; i < int(countInstances); i++ {
		var instanceName string
		var countProperties uint16
		_, err = fmt.Fscanf(file, "%s\n", &instanceName)
		if err != nil {
			return nil, err
		}
		_, err = fmt.Fscanf(file, "%d\n", &countProperties)
		if err != nil {
			return nil, err
		}
		instance := StaticInstanceSchema{
			instanceName,
			make([]StaticPropertySchema, countProperties),
		}
		for j := 0; j < int(countProperties); j++ {
			var propertyType int
			var propertyName string
			var propertyTypeName string
			line, err := file.ReadString('\n')
			matches := propMatcher.FindStringSubmatch(line)
			propertyType, _ = strconv.Atoi(matches[1])
			propertyName = matches[2]
			propertyTypeName = matches[3]
			if err != nil {
				return nil, err
			}
			instance.Properties[j] = StaticPropertySchema{propertyName, uint8(propertyType), propertyTypeName, &instance}
		}
		instances[i] = instance
	}
	return instances, nil
}

func parsePropSchema(filename string, schema []StaticInstanceSchema) ([]StaticPropertySchema, error) {
	var countProperties uint16

	propFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	file := bufio.NewReader(propFile)
	_, err = fmt.Fscanf(file, "%d\n", &countProperties)
	if err != nil {
		return nil, err
	}
	props := make([]StaticPropertySchema, countProperties)

	propMatcher := regexp.MustCompile(`(\d+)\s+'([a-zA-Z0-9 _]+)'\s+(\w+)\s+(\d+)\s*`)

	for i := 0; i < int(countProperties); i++ {
		var propertyType int
		var propertyName string
		var propertyTypeName string
		var classID int
		line, err := file.ReadString('\n')
		matches := propMatcher.FindStringSubmatch(line)
		propertyType, _ = strconv.Atoi(matches[1])
		propertyName = matches[2]
		propertyTypeName = matches[3]
		classID, _ = strconv.Atoi(matches[4])
		if err != nil {
			return nil, err
		}
		props[i] = StaticPropertySchema{propertyName, uint8(propertyType), propertyTypeName, &schema[classID]}
	}
	return props, nil
}

func ParseStaticSchema(instanceFilename string, propertyFilename string) (*StaticSchema, error) {
	schema := &StaticSchema{}
	var err error
	schema.Instances, err = parseInstSchema(instanceFilename)
	if err != nil {
		return schema, err
	}
	schema.Properties, err = parsePropSchema(propertyFilename, schema.Instances)

	return schema, err
}