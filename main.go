package filterComments

import (
	"database/sql"
	"github.com/mrtztg/console_backend/src/cons"
	"github.com/mrtztg/console_backend/src/queries/otherQueries"
	"github.com/mrtztg/console_backend/src/regexUtils"
	"strings"
)

func reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func FilterBadWords(content string, db *sql.DB) (string, int, int, []string) {

	var (
		//--make map for save main Content character position and when want to->
		//<- filter inside "for" we wil use
		mkMap = make(map[int]int)
	)

	//content=reverse(content)
	//--black list is list of bad words u want filter
	//var blackList = []string{"کصشر", "خر", "fuck", "مادرجنده", "جنده", "کصکش", "کص"}

	errCondition, blackListWordSlice := otherQueries.SelectWordBlackList(db)
	if errCondition == cons.R_FAILED {
		return "", 0, 0, nil
	}
	blackList := blackListWordSlice
	//useful.PrintErr("RegexQuery.go", 41, blackList)
	//--white list is list of "not bad words" but it's on black list<-
	//<- that is inside middle of bad words for example : "dick" is bad word and surely will filter->
	//<- and u want "dickensian" not going be filter. for not filter u must add "dickensian" to white list.
	//var whiteList = []string{"خربزه", "خرمن","fuck"}

	errCondition, whiteListSlice := otherQueries.SelectWordWhiteList(db)
	if errCondition == cons.R_FAILED {
		return "", 0, 0, nil
	}
	whiteList := whiteListSlice
	//useful.PrintErr("RegexQuery.go", 56, whiteList)
	//file, _ := ioutil.ReadFile("badWords.txt") //-- This file has our bad words. You can add more the word within it

	//blackList := strings.Split(string(file), "\n") // -- I want to separate words in new line

	//-- Length The longest word entire bad word list.
	largestLengthOfBlackList := LongestWordsList(blackList)

	//-- content which has bad words and convert to []rune
	mainContentToRune := []rune(content)

	//--allNumberPosition is the position of words which exist the whitelist. ->
	//<-Also in main content exist
	var joinAllNumberPositions []int

	//-- we want to remove whitelist from main content which our filter process will not see this words
	joinAllNumberPositions = RemoveWordAndSaveItSPosition(mainContentToRune, largestLengthOfBlackList, whiteList, joinAllNumberPositions)

	//--here we cut "notBadWords" from main content and ->
	//<- this variable is our mainContent but we already removed every words which is inside <-
	//<-white list
	var sliceString []rune

	//--add character instead notBadWords(whitelist).  after the filter, we can restore that words
	sliceString = ReplaceCharacter(mainContentToRune, sliceString, joinAllNumberPositions)

	//--main content is without whitelist instead whitelist we replace this "&" character
	mainContentToRune = sliceString

	//--we just save position of true character
	SaveCharacterPosition(mainContentToRune, mkMap)

	//-- Every character will go filter <==
	//==> we want only standard character my means we will remove Symbols like (*&^%$#@!) etc ...

	//--Filtered content of symbols (*&^%$#@!)
	res := regexUtils.RegexOnlyWords.FindAllString(string(mainContentToRune), -1)

	//--We join our content without space
	joinString := strings.Join(res, "")

	//--convert to Lowercase
	joinString = strings.ToLower(joinString)

	//--convert to []rune again
	mainContentWithoutWhiteList := []rune(joinString)

	//-- we want to know filter words position
	var (
		filterWordsPosition           []int
		countOfBadWords, countOfStars int
		listOfBadWords                []string
	)
	//-- Will filter one by one characters
	filterWordsPosition, mainContentToRune, countOfBadWords, countOfStars, listOfBadWords = FilterBlackListWords(mainContentWithoutWhiteList, mainContentToRune,
		largestLengthOfBlackList, blackList, filterWordsPosition, mkMap)

	//--finalContent is first content without any changes
	finalContent := []rune(content)

	FinalContent(joinAllNumberPositions, filterWordsPosition, sliceString, finalContent, mainContentToRune)
	return string(mainContentToRune), countOfBadWords, countOfStars, listOfBadWords
}

func LongestWordsList(wordList [][]rune) int {
	largestLength := 0
	for _, v := range wordList {
		if len(v) > largestLength {
			largestLength = len(v)
		}
	}
	return largestLength
}

func FinalContent(joinAllNumberPositions, filterWordsPosition []int, sliceString, finalContent, mainContentToRune []rune) {
	for i := 0; i < len(mainContentToRune)+len(joinAllNumberPositions); i++ {
		//--replace filtered words in final content using filterWordPosition [] int
		for _, v2 := range filterWordsPosition {
			if i == v2 {
				sliceString[i] = '*'
				continue
			}
		}

		//-- here we restore whitelist
	s:
		for _, v2 := range joinAllNumberPositions {
			if i == v2 {
				sliceString[i] = finalContent[i]
				continue s
			}
		}
	}
}

func FilterBlackListWords(mainContentWithoutWhiteList, mainContentToRune []rune,
	largestLength int, blackList [][]rune, filterWordsPosition []int,
	mkMap map[int]int) ([]int, []rune, int, int, []string) {
	var countOfBadWords, countOfStars int
	var listOfBadWords []string
	var wordToCheck []rune
	mainContentWithoutWhitelistLen := len(mainContentWithoutWhiteList)
	for start := 0; start < mainContentWithoutWhitelistLen; start++ {
		for offset := 1; offset < mainContentWithoutWhitelistLen+1-start && offset < largestLength; offset++ {
			wordToCheck = mainContentWithoutWhiteList[start : start+offset]
		firstFor:
			for _, v := range blackList {
				if offset == len(v) {
					for i := range wordToCheck {
						if wordToCheck[i] != v[i] {
							continue firstFor
						}
					}
					countOfBadWords++
					countOfStars += len(wordToCheck)
					listOfBadWords = append(listOfBadWords, string(v))
					//useful.PrintErr("RegexQuery.go", 138, v)
					for i := start; i < start+offset; i++ {
						filterWordsPosition = append(filterWordsPosition, mkMap[i])
						//--going to filter when words == black list words
						mainContentToRune = replaceAtIndex6(mainContentToRune, '*', mkMap[i])
					}
				}
			}
		}
	}
	return filterWordsPosition, mainContentToRune, countOfBadWords, countOfStars, listOfBadWords
}

func SaveCharacterPosition(mainContentToRune []rune, mkMap map[int]int) {
	j := 0
	for i, v := range mainContentToRune {
		if Index2(string(v)) {
			mkMap[j] = i
			j++
		}
	}
}

func ReplaceCharacter(mainContentToRune, sliceString []rune, joinAllNumberPositions []int) []rune {
n:
	for i, v := range mainContentToRune {
		for _, v2 := range joinAllNumberPositions {
			if i == v2 {
				v = '&'
				sliceString = append(sliceString, v)
				continue n
			}
		}
		sliceString = append(sliceString, v)
	}
	return sliceString
}

func RemoveWordAndSaveItSPosition(mainContentToRune []rune, largestLength int, wordList [][]rune, positionsOfWords []int) []int {
	mainContentToRuneLen := len(mainContentToRune)
	for start := 0; start < mainContentToRuneLen; start++ {
		for offset := 1; offset < mainContentToRuneLen+1-start && offset < largestLength; offset++ {
			wordToCheck := mainContentToRune[start : start+offset]
		firstFor:
			for _, v := range wordList {
				if len(v) == offset {
					for i := range wordToCheck {
						if wordToCheck[i] != v[i] {
							continue firstFor
						}
					}
					for i := start; i < start+offset; i++ {
						positionsOfWords = append(positionsOfWords, i)
					}
				}
			}
		}
	}
	return positionsOfWords
}

//-- Will replace * instead number mkMap[...] (index int )
func replaceAtIndex6(str []rune, replacement rune, index int) []rune {
	tmp := str
	tmp[index] = replacement
	return tmp
}

func replaceAtIndex5(str string, replacement rune, index int) []rune {
	var i int
	var t1 []rune

	tmp := []rune(str)

	for i = 0; i < index && i < len(tmp); i++ {
		t1 = append(t1, tmp[i])
	}
	if index < len(tmp) {
		t1 = append(t1, replacement)
	}
	for i = i + 1; i < len(tmp); i++ {
		t1 = append(t1, tmp[i])
	}

	return t1
}

//--Want to know that (s) character we have, it had, that will be added to map (mkMap) <==
//==> U can add another language characters

func Index2(s string) (boolean bool) {
	content := "abcdefghijklmnopqrstuwxyzABCDEFGHIJKLMNOPQRSTUVWXYZابپتثجچحخدذرزسشصضطظعغفقکگلمنوهی"
	for _, v := range content {
		if string(v) == s {
			return true
		}
	}
	return false
}

//func IndexRune2(s rune) bool {
//	for i := 65; i < 1611; i++ {
//		if rune(i) == s {
//			return true
//		}
//		if i > 90 && i < 97 || i > 122 && i < 1568 {
//			continue
//		}
//	}
//	return false
//}
