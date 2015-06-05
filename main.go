package main

// This command will allow you to provision, delete, describe, or estimate the cost of the specified CloudFormation template.
//
// Once compiled use the -help flag for details.
// Initital source from http://junctionbox.ca/2015/05/02/golang-aws-cloudformation.html
// https://gist.github.com/nfisher/522c303ef325bd5cf43e

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	awscf "github.com/aws/aws-sdk-go/service/cloudformation"
)

func provisionStack(svc *awscf.CloudFormation, b []byte, params []*awscf.Parameter, stackName string) {
	input := &awscf.CreateStackInput{
		StackName: aws.String(stackName),
		Capabilities: []*string{
			aws.String("CAPABILITY_IAM"),
		},
		OnFailure:        aws.String("DELETE"),
		Parameters:       params,
		TemplateBody:     aws.String(string(b)),
		TimeoutInMinutes: aws.Long(20),
	}
	resp, err := svc.CreateStack(input)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(awsutil.StringValue(resp))
}

func delStack(svc *awscf.CloudFormation, stackName string) {
	input := &awscf.DeleteStackInput{
		StackName: aws.String(stackName),
	}

	_, err := svc.DeleteStack(input)
	if err != nil {
		log.Fatal(err)
	}
	// the log.Println ends up looking like
	// 2015/06/04 16:55:36 {
	//
	// }
	//
	// log.Println(awsutil.StringValue(resp))
}

func descStack(svc *awscf.CloudFormation, stackName string) {
	input := &awscf.DescribeStackEventsInput{
		StackName: aws.String(stackName),
	}
	resp, err := svc.DescribeStackEvents(input)
	if err != nil {
		log.Fatal(err)
	}

	if len(resp.StackEvents) > 0 {
		log.Println(awsutil.StringValue(resp.StackEvents[0]))
	}
}

func cost(svc *awscf.CloudFormation, b []byte, params []*awscf.Parameter) {
	estInput := &awscf.EstimateTemplateCostInput{
		Parameters:   params,
		TemplateBody: aws.String(string(b)),
	}

	cost, err := svc.EstimateTemplateCost(estInput)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*cost.URL)
}

func watch(svc *awscf.CloudFormation, verbose bool, interval int, stackName string) {
	req := &awscf.DescribeStacksInput{StackName: aws.String(stackName)}
	var maxError int
	previousStatus := ""
	var err error
	var firstLoop bool
	for maxError < 3 {
		if !firstLoop {
			time.Sleep(time.Duration(interval) * time.Second)
		} else {
			firstLoop = true
		}

		resp, err := svc.DescribeStacks(req)
		if err != nil {
			if previousStatus == "DELETE_IN_PROGRESS" {
				fmt.Printf("%s Finished\n", time.Now().Format(time.RFC3339))
				return
			}
			fmt.Printf("Error: %s - retrying\n", err)
			maxError++
			continue
		}
		for _, stack := range resp.Stacks {
			if *stack.StackName == stackName {
				if *stack.StackStatus != previousStatus || verbose {
					fmt.Printf("%s %s\n", time.Now().Format(time.RFC3339), *stack.StackStatus)
					previousStatus = *stack.StackStatus
				}
			}
			if strings.HasSuffix(previousStatus, "COMPLETE") {
				fmt.Printf("%s Finished\n", time.Now().Format(time.RFC3339))
				return
			}
		}
	}
	fmt.Printf("Error: %s - giving up\n", err)
}

func main() {
	var templateFile string
	var outputCost bool
	var provision bool
	var desc bool
	var del bool
	var status bool
	var b []byte
	var params []*awscf.Parameter
	var stackName string
	var region string
	var verbose bool
	var interval int
	var prompt bool

	flag.StringVar(&region, "region", "us-west-2", "AWS region to provision script to.")
	flag.StringVar(&templateFile, "template", "", "Template to validate.")
	flag.StringVar(&stackName, "name", "", "Stack name (required).")
	flag.BoolVar(&outputCost, "cost", false, "Output cost URL.")
	flag.BoolVar(&provision, "provision", false, "Provision template.")
	flag.BoolVar(&prompt, "prompt", false, "Prompt for param values")
	flag.BoolVar(&desc, "desc", false, "Describe stack.")
	flag.BoolVar(&del, "del", false, "Delete stack.")
	flag.BoolVar(&status, "watch", false, "")
	flag.BoolVar(&verbose, "v", false, "Verbose output for watch")
	flag.IntVar(&interval, "i", 5, "Polling interval in seconds for watch")
	flag.Parse()

	if stackName == "" {
		fmt.Println("Stack name cannot be empty!")
		flag.Usage()
		return
	}

	config := &aws.Config{Region: region}
	svc := awscf.New(config)

	if outputCost || provision {
		f, err := os.Open(templateFile)
		if err != nil {
			log.Fatal(err)
		}

		b, err = ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}

		input := &awscf.ValidateTemplateInput{
			TemplateBody: aws.String(string(b)),
		}
		resp, err := svc.ValidateTemplate(input)
		if err != nil {
			log.Fatal(err)
		}

		// output the template description
		fmt.Println(awsutil.StringValue(resp.Description))
		params = make([]*awscf.Parameter, len(resp.Parameters))

		// fill out the parameters from the template
		// using either user prompting or defaults
		if prompt {
			stdin := bufio.NewReader(os.Stdin)
			for i, p := range resp.Parameters {
				fmt.Printf("%v (%v): ", awsutil.StringValue(p.Description), awsutil.StringValue(p.DefaultValue))

				// don't care about isMore if someone's typing so much oh well
				b, _, err := stdin.ReadLine()
				if err != nil {
					log.Fatal(err)
				}
				line := string(b)

				params[i] = &awscf.Parameter{
					ParameterKey:     p.ParameterKey,
					UsePreviousValue: aws.Boolean(true),
				}

				if line != "" {
					params[i].ParameterValue = aws.String(line)
				} else {
					params[i].ParameterValue = p.DefaultValue
				}
			}
		} else { // using default values plus any overrides passed on command line
			overrides := make(map[string]string)
			args := flag.Args()
			if len(args)%2 == 0 {
				for i := 0; i < len(args); i++ {
					overrides[args[i]] = args[i+1]
					i++
				}
			} else {
				fmt.Printf("ERROR: If passing template parameters on command line, must be even number of keys and values")
				os.Exit(1)
			}
			for i, p := range resp.Parameters {
				params[i] = &awscf.Parameter{
					ParameterKey:   p.ParameterKey,
					ParameterValue: p.DefaultValue,
				}
				if val, ok := overrides[*p.ParameterKey]; ok {
					params[i].ParameterValue = aws.String(val)
					fmt.Printf("Overriding %s -> %s\n", *p.ParameterKey, val)
					delete(overrides, *p.ParameterKey)
				}
			}
			if len(overrides) != 0 {
				fmt.Println("ERROR: override params passed on command line but not found in template")
				fmt.Println("Unused params listed below")
				for k, v := range overrides {
					fmt.Printf("param %s value %s\n", k, v)
				}
				os.Exit(1)
			}
		}
	}

	if outputCost {
		cost(svc, b, params)
		return
	} else if provision {
		provisionStack(svc, b, params, stackName)
	} else if desc {
		descStack(svc, stackName)
	} else if del {
		delStack(svc, stackName)
	}
	if status {
		watch(svc, verbose, interval, stackName)
	}
}
