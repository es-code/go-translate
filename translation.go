package trans

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var localConfig *appLocalConfig

func CreateAppLocal(defaultLang string) *appLocalConfig {
	if defaultLang != ""{
		defaultLang = "en"
	}
	local:=appLocalConfig{
		defaultLang:        defaultLang,
		supportedLanguages: nil,
		translationsFiles:  nil,
	}
	localConfig = &local
	return &local
}

func (local *appLocalConfig) LoadTranslationsFiles(SupportedLanguages ...string)  {
	filesMap:=make(map[string]map[string]map[string]interface{})
	//check if translations dir is existed
	dirExistOrCreate("translations")

	for i:=0;i<len(SupportedLanguages);i++ {

		langPath:="translations/"+SupportedLanguages[i]
		//check lang dir is not exists
		if dirExistOrCreate(langPath) == false{
			createSampleFile(langPath)
		}

		files, err := ioutil.ReadDir(langPath)
		if err != nil {
			log.Fatal(err)
		}
		fileMap:=make(map[string]map[string]interface{})
		for _, f := range files {
			if f.IsDir() != true{
				content, err := ioutil.ReadFile(langPath+"/"+f.Name())
				if err != nil{
					log.Fatal(err)
				}
				var data map[string]interface{}
				Json:=data
				errJ:=json.Unmarshal(content,&Json)
				if errJ !=nil{
					log.Fatal(errJ)
				}
				fileMap[strings.Split(f.Name(),".json")[0]]=Json
			}
		}
		filesMap[SupportedLanguages[i]]=fileMap
	}
	local.translationsFiles = filesMap
	local.supportedLanguages = SupportedLanguages

}

func T(key string,langOrRequest interface{}) string {

	if stringHasDot(key) == false{
		return key
	}
	fileKeys:=strings.Split(key,".")
	fileMap:=localConfig.translationsFiles[localConfig.getLang(langOrRequest)][fileKeys[0]]
	message,exist:=getMessageValue(1,fileKeys,fileMap)
	if exist == false{
		return key
	}
	return message
}

func getMessageValue(start int, keys []string, fileMap map[string]interface{}) (string,bool)  {
	//check keys has this start key
	if start >= len(keys){
		return "",false
	}
	//if this key value is string
	if message, ok := fileMap[keys[start]].(string); ok {
		return message,true
		//if this key value is other map
	} else if messageMap, ok := fileMap[keys[start]].(map[string]interface{}); ok {
		//get this map message value
		return getMessageValue(start+1,keys,messageMap)
	}else{
		//unexpected value
		return "",false
	}

}

func (local *appLocalConfig) getLang(langOrRequest interface{})  string {

	var language string

	if lang, ok := langOrRequest.(string); ok{

		language = lang
	}else if reflect.TypeOf(langOrRequest).String() == "*http.Request"{

		language =  langOrRequest.(*http.Request).FormValue("lang")
	}

	if language != "" && langIsSupported(local.supportedLanguages,language){

		return language
	}

	return local.defaultLang
}

func langIsSupported(supported []string,lang string) bool  {
	for i:=0;i < len(supported); i++ {
		if lang == supported[i]{
			return true
		}
	}
	return false
}

func stringHasDot(text string)  bool {
	result, _ := regexp.MatchString("\\.", text)
	return result
}

func createSampleFile(path string)  {
	err := ioutil.WriteFile(path+"/sample.json", []byte("{\n  \"email_exists\":\"this email is already registered before\",\n  \"success_register\": \"done register successfully\",\n  \"password_not_confirm\": \"password not matches\"\n}"), 0755)
	if err != nil {
		fmt.Printf("Unable to write sample file: %v", err)
	}
}

func dirExistOrCreate(dir string) bool {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		//create translations dir
		errCreateTranslations := os.Mkdir(dir, 0755)
		if errCreateTranslations != nil {
			log.Fatal(errCreateTranslations)
		}
		return false
	}
	return true
}
