package main

// Hamed

import (
	"fmt"
	"io/ioutil"
	"strings"
	"regexp"
	"time"
)

func main() {
	
	content := `my Fuck suck f . u. c. k salam khoobi ? f u c king are u ok?  and f  uck pussy   سکس(sex)`
	
	result := FilterBadWords(content)
	fmt.Println("main (17):::", result)
}

func FilterBadWords(content string) string {
	startTime := time.Now().UnixNano() //-- I want to know How long does it take
	
	var (
		mkMap      = make(map[int]int)
		mkSliceMap = make([]map[int]int, 0)
		j          = 0
	)
	
	file, _ := ioutil.ReadFile("text.txt") //-- This file has our bad words. You can add more the word within it
	
	split := strings.Split(string(file), "\n") // -- I want to separate words in new line
	
	largestLength := 0
	for _, v := range split {
		if len(v) > largestLength {
			largestLength = len(v)
		}
	} //-- Length The longest word.
	
	//-- Our content which has bad words
	
	//-- A bit hard to Understand but don't worry here we want to know what is our characters position <==
	//==> we can achieve from numbers. so alternative * bad words (into for)
	for i, v := range content {
		if Index(string(v)) {
			mkMap[j] = i
			j++
		}
	}
	mkSliceMap = append(mkSliceMap, mkMap) //-- append to slice map
	
	//-- Every character will go filter <==
	//==> we want only standard character my means we will remove Symbols like (*&^%$#@!) etc ...
	regex, _ := regexp.Compile(`[A-Za-zا-ی]`)
	
	res := regex.FindAllString(content, -1) //--Filtered content of symbols (*&^%$#@!)
	
	joinString := strings.Join(res, "") //--We join our content without space
	
	joinString = strings.ToLower(joinString) // --Note here we will do convert every word to Lowercase
	
	//-- Will filter one by one characters
	for start := 0; start < len(joinString); start++ {
		for offset := 1; offset < (len(joinString)+1-start) && offset < largestLength; offset++ {
			wordToCheck := joinString[start: start+offset]
			//fmt.Println("main (319):::", wordToCheck)
			condition := strings.Contains("", wordToCheck)
			for _, v := range split {
				if wordToCheck == v {
					condition = true
				}
				if condition {
					for i := start; i < start+offset; i++ {
						content = replaceAtIndex2(content, '*', mkMap[i])
					}
				}
			}
		}
	}
	
	fmt.Println("main (297):::", time.Now().UnixNano()-startTime, "Milisecond long does it take (depended Your CPU ;) )")
	return content
}

//-- Will replace * instead number mkMap (index int )
func replaceAtIndex2(str string, replacement rune, index int) string {
	return str[:index] + string(replacement) + str[index+1:]
}

//--Want to know that (s) character we have, it had, that will be added to map (mkMap) <==
//==> U can add another language characters
func Index(s string) (boolean bool) {
	content := "abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUVWXYZابپتثجچحخدذرزسشصضطظعغفقکگلمنوهی"
	for _, v := range content {
		if string(v) == s {
			return true
		}
	}
	return false
}
