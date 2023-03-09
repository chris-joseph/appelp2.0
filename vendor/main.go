package main

import (
	"C"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slices"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const exampleJSON = `
{"web-app": {
	"servlet": [   
	  {
		"servlet-name": "cofaxCDS",
		"servlet-class": "org.cofax.cds.CDSServlet",
		"init-param": {
		  "configGlossary:installationAt": "Philadelphia, PA",
		  "configGlossary:adminEmail": "ksm@pobox.com",
		  "configGlossary:poweredBy": "Cofax",
		  "configGlossary:poweredByIcon": "/images/cofax.gif",
		  "configGlossary:staticPath": "/content/static",
		  "templateProcessorClass": "org.cofax.WysiwygTemplate",
		  "templateLoaderClass": "org.cofax.FilesTemplateLoader",
		  "templatePath": "templates",
		  "templateOverridePath": "",
		  "defaultListTemplate": "listTemplate.htm",
		  "defaultFileTemplate": "articleTemplate.htm",
		  "useJSP": false,
		  "jspListTemplate": "listTemplate.jsp",
		  "jspFileTemplate": "articleTemplate.jsp",
		  "cachePackageTagsTrack": 200,
		  "cachePackageTagsStore": 200,
		  "cachePackageTagsRefresh": 60,
		  "cacheTemplatesTrack": 100,
		  "cacheTemplatesStore": 50,
		  "cacheTemplatesRefresh": 15,
		  "cachePagesTrack": 200,
		  "cachePagesStore": 100,
		  "cachePagesRefresh": 10,
		  "cachePagesDirtyRead": 10,
		  "searchEngineListTemplate": "forSearchEnginesList.htm",
		  "searchEngineFileTemplate": "forSearchEngines.htm",
		  "searchEngineRobotsDb": "WEB-INF/robots.db",
		  "useDataStore": true,
		  "dataStoreClass": "org.cofax.SqlDataStore",
		  "redirectionClass": "org.cofax.SqlRedirection",
		  "dataStoreName": "cofax",
		  "dataStoreDriver": "com.microsoft.jdbc.sqlserver.SQLServerDriver",
		  "dataStoreUrl": "jdbc:microsoft:sqlserver://LOCALHOST:1433;DatabaseName=goon",
		  "dataStoreUser": "sa",
		  "dataStorePassword": "dataStoreTestQuery",
		  "dataStoreTestQuery": "SET NOCOUNT ON;select test='test';",
		  "dataStoreLogFile": "/usr/local/tomcat/logs/datastore.log",
		  "dataStoreInitConns": 10,
		  "dataStoreMaxConns": 100,
		  "dataStoreConnUsageLimit": 100,
		  "dataStoreLogLevel": "debug",
		  "maxUrlLength": 500}},
	  {
		"servlet-name": "cofaxEmail",
		"servlet-class": "org.cofax.cds.EmailServlet",
		"init-param": {
		"mailHost": "mail1",
		"mailHostOverride": "mail2"}},
	  {
		"servlet-name": "cofaxAdmin",
		"servlet-class": "org.cofax.cds.AdminServlet"},
   
	  {
		"servlet-name": "fileServlet",
		"servlet-class": "org.cofax.cds.FileServlet"},
	  {
		"servlet-name": "cofaxTools",
		"servlet-class": "org.cofax.cms.CofaxToolsServlet",
		"init-param": {
		  "templatePath": "toolstemplates/",
		  "log": 1,
		  "logLocation": "/usr/local/tomcat/logs/CofaxTools.log",
		  "logMaxSize": "",
		  "dataLog": 1,
		  "dataLogLocation": "/usr/local/tomcat/logs/dataLog.log",
		  "dataLogMaxSize": "",
		  "removePageCache": "/content/admin/remove?cache=pages&id=",
		  "removeTemplateCache": "/content/admin/remove?cache=templates&id=",
		  "fileTransferFolder": "/usr/local/tomcat/webapps/content/fileTransferFolder",
		  "lookInContext": 1,
		  "adminGroupID": 4,
		  "betaServer": true}}],
	"servlet-mapping": {
	  "cofaxCDS": "/",
	  "cofaxEmail": "/cofaxutil/aemail/*",
	  "cofaxAdmin": "/admin/*",
	  "fileServlet": "/static/*",
	  "cofaxTools": "/tools/*"},
   
	"taglib": {
	  "taglib-uri": "cofax.tld",
	  "taglib-location": "/WEB-INF/tlds/cofax.tld"}}}
`

func main() {
	print(C.GoString(ParseJson(C.CString(exampleJSON), C.CString("TestClass"))))
}

func getDartTypeList() []string {
	return []string{"int", "double", "String", "bool", "dynamic"}
}

//export ParseJson
func ParseJson(jsons *C.char, name *C.char) *C.char {
	var sb strings.Builder
	var jsonText = C.GoString(jsons)
	var className = C.GoString(name)
	valid := json.Valid([]byte(jsonText))
	if !valid {
		return C.CString("Invalid JSON")
	}
	jsonParse(jsonText, className, &sb)
	return C.CString(sb.String())
}

func jsonParse(str string, name string, sb *strings.Builder) {
	var output map[string]any
	var params = make(map[string]string)
	name = stringValidator(name)
	json.Unmarshal([]byte(str), &output)
	for key, val := range output {
		params[stringValidator(key)] = getType(val, key, sb, params)
	}
	classBuilder(sb, name, params)
}
func getType(val any, key string, sb *strings.Builder, params map[string]string) string {
	k := reflect.ValueOf(val).Kind().String()

	switch k {

	case "bool":
		return "bool"
	case "float32", "float64":
		res := strconv.FormatFloat(val.(float64), 'f', 6, 64)
		v, _ := strconv.ParseFloat(res, 64)
		if val == float64(int(v)) {
			return "int"
		}
		return "double"
	case "map":
		empData, err := json.Marshal(val)
		if err != nil {
			fmt.Println(err.Error())
			return "Map<dynamic,dynamic>"

		}
		caser := cases.Title(language.Und)
		jsonStr := string(empData)
		if _, ok := params[key]; ok {
			return stringValidator(caser.String(strings.ToLower(key)))
		}
		jsonParse(jsonStr, caser.String((strings.ToLower(key))), sb)
		return stringValidator((caser.String(strings.ToLower(key))))

	case "slice":
		j := val.([]any)
		//TODO add logic to sepereate double and int by iterating all values
		t := getType(j[0], key, sb, params)

		return "List<" + t + ">"
	case "string":

		_, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", val))
		if err == nil {
			return "DateTime"
		}

		return "String"
	default:
		return "dynamic"
	}
}
func stringValidator(name string) string {
	name = strings.ReplaceAll(name, "_", "")
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, ":", "")
	name = strings.ReplaceAll(name, ",", "")
	name = strings.ReplaceAll(name, ".", "")
	return name
}
func classBuilder(sb *strings.Builder, name string, params map[string]string) {
	sb.WriteString("\n")
	sb.WriteString("\n")
	sb.WriteString(`class ` + name + ` {`)
	sb.WriteString("\n")
	for k, v := range params {
		sb.WriteString("  final " + v + "? " + k + ";")
		sb.WriteString("\n")
	}
	sb.WriteString("  " + makeConstructor(name, params))
	sb.WriteString("\n")
	sb.WriteString("  " + makeFromJson(name, params))
	sb.WriteString("\n")
	sb.WriteString("  " + makeToJson(name, params))
	sb.WriteString("\n")
	sb.WriteString("  " + makeCopyWith(name, params))
	sb.WriteString("\n")
	sb.WriteString("}")
	sb.WriteString("\n")

}

func makeConstructor(name string, params map[string]string) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("const " + name + "({")
	for k := range params {
		sb.WriteString(" this." + k + ",")

	}
	sb.WriteString("});")
	return sb.String()

}

func makeFromJson(name string, params map[string]string) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("factory " + name + ".fromJson(Map<String, dynamic> json)=>" + "\n")
	sb.WriteString(name + "(")
	sb.WriteString("\n")
	for k, v := range params {
		if strings.Contains(v, "List") {
			sb.WriteString(makeListFromJson(k, v))
			continue
		}
		if strings.Contains(v, "DateTime") {
			sb.WriteString(" " + k + ":" + "DateTime.parse(json[\"" + k + "\"])" + ",")
			continue
		}
		if !slices.Contains(getDartTypeList(), v) {
			sb.WriteString(" " + k + ":" + v + ".fromJson(json[\"" + k + "\"])" + ",")
			sb.WriteString("\n")
			continue
		}
		sb.WriteString(" " + k + " : " + "json[\"" + k + "\"]" + ",")
		sb.WriteString("\n")
	}

	sb.WriteString(");")
	return sb.String()

}

func makeListFromJson(k string, v string) string {
	j := strings.Split(v, "<")
	j = j[1:]
	j = strings.Split(strings.Join(j, ""), ">")
	typ := strings.Join(j, "")
	if slices.Contains(getDartTypeList(), typ) {
		return " " + k + " : " + "json[\"" + k + "\"]" + "," + "\n"
	}

	return " " + k + " : " + v + ".from(json[\"" + k + "\"].map((x) => " + typ + ".fromJson(x)))" + "," + "\n"

}

func makeListToJson(k string, v string) string {
	j := strings.Split(v, "<")
	j = j[1:]
	j = strings.Split(strings.Join(j, ""), ">")
	typ := strings.Join(j, "")
	if slices.Contains(getDartTypeList(), typ) {
		return " \"" + k + "\" : " + k + "," + "\n"
	}

	return " \"" + k + "\" : " + k + "==null?null:List<dynamic>" + ".from(" + k + "!.map((x) => " + "x.toJson()))" + "," + "\n"

}

func makeCopyWith(name string, params map[string]string) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString(name + " copyWith({" + "\n")
	for k, v := range params {

		sb.WriteString(" " + v + "? " + k + ",")
		sb.WriteString("\n")

	}
	sb.WriteString("})=> " + name + "(\n")
	for k := range params {

		sb.WriteString(" " + k + " : " + k + "?? this." + k + ",")
		sb.WriteString("\n")

	}
	sb.WriteString(");\n")
	return sb.String()
}

func makeToJson(name string, params map[string]string) string {
	var sb strings.Builder

	sb.WriteString("\n")
	sb.WriteString("Map<String, dynamic> toJson(){" + "\n")
	sb.WriteString("var map= {" + "\n")

	for k, v := range params {
		if strings.Contains(v, "List") {
			sb.WriteString(makeListToJson(k, v))
			continue
		}
		if strings.Contains(v, "DateTime") {
			sb.WriteString(" \"" + k + "\" : " + k + "?.toIso8601String()" + ",")
			continue
		}
		if !slices.Contains(getDartTypeList(), v) {
			sb.WriteString(" \"" + k + "\" : " + k + "?.toJson()" + ",")
			sb.WriteString("\n")
			continue
		}
		sb.WriteString(" \"" + k + "\" : " + k + ",")
		sb.WriteString("\n")

	}
	sb.WriteString("};\n")
	sb.WriteString("map.removeWhere((key, value) => value == null);\n")
	sb.WriteString("return map;\n")
	sb.WriteString("}\n")
	return sb.String()
}
