package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

func mutantFound(msj string) {
	//fmt.Println("You don't seem to have a hooman DNA")
	//fmt.Println(msj)
	//os.Exit(0)
}

func checksq(dna []string) bool {
	if len(dna) != len(dna[0]) {
		mutantFound("nxn")
		return true
	}
	return false
}

func checkGene(gene string) bool {

	// validate they're valid characters
	// TO DO: should we take lowercase as valid input with strings.ToUpper(str)?
	allels := []string{"A", "T", "C", "G"}
	for i := 0; i < len(gene); i++ {

		if !stringInSlice(string(gene[i]), allels) {
			mutantFound("non char")
			return true
		}

	}
	return false

}

func checkConcurrent(gene string) bool {

	// TO DO: benchmark if running an array comparison is chepear on cpu and time
	if strings.Count(gene, "AAAA") > 0 {
		mutantFound(gene + "A")
		return true
	}
	if strings.Count(gene, "TTTT") > 0 {
		mutantFound(gene + "T")
		return true
	}
	if strings.Count(gene, "CCCC") > 0 {
		mutantFound(gene + "C")
		return true
	}
	if strings.Count(gene, "GGGG") > 0 {
		mutantFound(gene + "G")
		return true
	}

	return false
}

func linearSearch(dna []string) bool {
	length := len(dna)
	for i := 0; i < length; i++ {
		var v string
		for c := 0; c < length; c++ {
			//TO DO: check io.Writer performance
			v = v + string([]rune(dna[c])[i])
		}
		if checkConcurrent(v) == true {
			return true
		}
	}
	return false
}

// search left to right
func diagSearch(dna []string, y string) bool {

	length := len(dna)
	for i := 0; i < length; i++ {
		var diag string
		var v int

		// avoid shorter than possibly mutant chains
		if length-i < 4 {
			break
		}

		for c := 0; c < length; c++ {
			// go left
			v = c + i

			// avoid going out of bounds
			if v >= length {
				break
			}

			if y == "hoz" {
				diag = diag + string([]rune(dna[c])[v])
			}
			if y == "ver" {
				diag = diag + string([]rune(dna[v])[c])
			}

		}
		if checkConcurrent(diag) == true {
			return true
		}
		//fmt.Println(checkConcurrent(diag))
	}
	return false
}

// search right to left
func diagSearchRTL(dna []string, y string) bool {

	length := len(dna)
	var reverse int

	for i := 0; i < length; i++ {
		var diag string
		var v int
		var count int
		var dwn int

		count = 0
		for c := length - 1; c >= 0; c-- {
			v = c - i
			dwn = count + i

			if y == "ver" {
				diag = diag + string([]rune(dna[dwn])[c])
			}

			if y == "hoz" {
				diag = diag + string([]rune(dna[count])[v])
			}
			count++
			reverse++
			dwn++
			if v <= 0 {
				break
			}
		}

		checkConcurrent(diag)
		if count < 4 {
			break
		}
	}
	return false

}

func isMutant(dna []string) bool {
	mutantBoolean := false

	// CHECK ITS NxN
	if checksq(dna) == true {
		mutantBoolean = true
	}

	if mutantBoolean == false {
		// CHECK EACH ROW LINEAR
		for i := 0; i < len(dna); i++ {
			if checkGene(dna[i]) == true || checkConcurrent(dna[i]) == true {
				mutantFound("non char or linear")
				mutantBoolean = true
			}
		}
	}

	// CHECK DOWN LEFT TO RIGHT
	if mutantBoolean == false {
		if linearSearch(dna) == true {
			mutantBoolean = true
		}
	}

	// CHECK DIAGONALLY LEFT TO RIGHT
	if mutantBoolean == false {
		if diagSearch(dna, "hoz") == true {
			mutantBoolean = true
		}
	}

	// CHECK DIAGONALLY DOWN RIGHT TO LEFT
	if mutantBoolean == false {
		if diagSearch(dna, "ver") == true {
			mutantBoolean = true
		}
	}

	// CHECK RIGHT TO LEFT
	if mutantBoolean == false {
		if diagSearchRTL(dna, "hoz") == true {
			mutantBoolean = true
		}
	}

	// CHECK DIAGONALLY DOWN RIGHT TO LEFT
	if mutantBoolean == false {
		if diagSearchRTL(dna, "ver") == true {
			mutantBoolean = true
		}
	}

	return mutantBoolean
}

func dynamoStore(dna []string, table string) {

}

var statusCode int

func lambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "You're human",
		StatusCode: 200,
	}, nil
}

func main() {

	//no mutante
	//dna := []string{"ATGCGC", "CAGTGC", "TTATTT", "AGACGG", "GCGTCA", "TCACTG"}

	// meli dna
	// dna := []string{
	// 	"ATGCGA",
	// 	"CAGTGC",
	// 	"TTATTT",
	// 	"AGACGG",
	// 	"GCGTCA",
	// 	"TCACTG"}

	//mutant diag
	// dna := []string{
	// 	"ATAATA",
	// 	"CATTGC",
	// 	"TTATTT",
	// 	"ATACTG",
	// 	"GCTTCT",
	// 	"TCATTG"}

	// mutant horiz
	//dna := []string{"ATGCGT", "TAATAT", "TTATGT", "AGAAGG", "CCCCTA", "TCACTG"}

	// mutant vert
	//dna := []string{"ATGCGT", "TAATAT", "TTATGT", "AGAAGG", "CTACTA", "TCCCTG"}

	//mutant last diag
	dna := []string{"ATGCGT", "TAATAT", "TTGTGT", "TGAGGG", "TTACGA", "CCGCTG"}

	// TO DO: Get dna from request
	if isMutant(dna) == true {
		statusCode = 403
		dynamoStore(dna, "mutant")
	} else {
		statusCode = 200
		dynamoStore(dna, "human")
	}

	fmt.Println(statusCode)

	// SETTING UP AWS
	// Create Lambda service client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	client := lambda.New(sess, &aws.Config{Region: aws.String("us-west-2")})

	lambda.Start(lambdaHandler)

}
