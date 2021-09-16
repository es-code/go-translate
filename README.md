# go-translate
Localization package for GOLANG using json files.

# Installation
Install the package with the command

`go get github.com/es-code/go-translate`
<hr>

## Usage :

```
func main() {
        //create app localization with default language
	Local := trans.CreateAppLocal("en")
	
	//load translations files from folders 
	//if this folders dosn't exists "go-translate" will create them and create sample json file
	Local.LoadTranslationsFiles("en","ar")
	
	//translate to en 
	fmt.Println(trans.T("sample.email_exists","en"))
	
	//translate to ar 
	fmt.Println(trans.T("sample.email_exists","ar"))

}	
```
### Using in Request
detect language from "lang" parameter sent at request "query-string" or form value
```

func Login(w http.ResponseWriter, r *http.Request){
    fmt.Println(trans.T("sample.email_exists",r))
}

//translate to en
http://localhost/login?lang=en

//translate to ar
http://localhost/login?lang=ar
```


