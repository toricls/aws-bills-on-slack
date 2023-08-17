# aws-bills-on-slack
A monkey increases your visibility on AWS bills via Slack, running on AWS Lambda on a daily basis 🐒

![aws-bills-on-slack](./aws-bills-on-slack.png)

## Prerequisites

- Slack Incoming Webhook
- An AWS account with AWS Organizations enabled
  - `aws-bills-on-slack` uses [`toricls/acos`](https://github.com/toricls/acos) as its dependency so you may find more details about the prerequisites in the [`toricls/acos`'s prerequisites](https://github.com/toricls/acos/blob/main/README.md#prerequisites) to run.

## Installation

First clone this repository to your local machine or somewhere. You can use any tools whatever you want like AWS SAM, AWS CloudFormation, or Terraform, to deploy `aws-bills-on-slack` on AWS. See the following sections for more details.

### Install using AWS SAM (tested and recommended)

**For the first deployment:**

> **Note**  
> You'll be prompted to enter some parameters after running the `sam deploy` command. See the [template.yaml's "Parameters" section](./deploy/aws-sam/template.yaml) for more info about the parameters.

```shell
cd /path/to/your/aws-bills-on-slack
sam build --template ./deploy/aws-sam/template.yaml
sam deploy --guided --stack-name aws-bills-on-slack --capabilities CAPABILITY_IAM # `--profile` and/or `--region` options may be also needed
```

_We'd strongly recommend you to set `./deploy/aws-sam/samconfig.toml` for the `SAM configuration file` parameter that is prompted by the `sam deploy` command. This will make your subsequent deployments much easier._

<details>
  <summary>See full example output of `sam deploy`</summary>

```bash
$ sam deploy --guided --stack-name aws-bills-on-slack --capabilities CAPABILITY_IAM

Configuring SAM deploy
======================

	Looking for config file [samconfig.toml] :  Not found

	Setting default arguments for 'sam deploy'
	=========================================
	Stack Name [aws-bills-on-slack]:
	AWS Region [ap-northeast-1]:
	Parameter MessageText [Here's the daily bill on our AWS accounts:]:
	Parameter OuId []: <YOUR-OU-ID-HERE>
	Parameter SlackIncomingWebhookUrl []: <YOUR-SLACK-INCOMING-WEBHOOK-URL-HERE>
	Parameter CronScheduleString [cron(0 0 * * ? *)]: cron(0 0 * * ? *)
	#Shows you resources changes to be deployed and require a 'Y' to initiate deploy
	Confirm changes before deploy [y/N]: y
	#SAM needs permission to be able to create roles to connect to the resources in your template
	Allow SAM CLI IAM role creation [Y/n]: Y
	#Preserves the state of previously provisioned resources when an operation fails
	Disable rollback [y/N]: N
	Save arguments to configuration file [Y/n]: Y
	SAM configuration file [samconfig.toml]: ./deploy/aws-sam/samconfig.toml
	SAM configuration environment [default]: default

	Looking for resources needed for deployment:

	Managed S3 bucket: aws-sam-cli-managed-default-samclisourcebucket-123456789012
	A different default S3 bucket can be set in samconfig.toml and auto resolution of buckets turned off by setting resolve_s3=False

	Saved arguments to config file
	Running 'sam deploy' for future deployments will use the parameters saved above.
	The above parameters can be changed by modifying samconfig.toml
	Learn more about samconfig.toml syntax at
	https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-config.html

	Uploading to aws-bills-on-slack/6f278f1f8d723a82701a6f561bc16c96  7987705 / 7987705  (100.00%)

	Deploying with following values
	===============================
	Stack name                   : aws-bills-on-slack
	Region                       : ap-northeast-1
	Confirm changeset            : True
	Disable rollback             : False
	Deployment s3 bucket         : aws-sam-cli-managed-default-samclisourcebucket-123456789012
	Capabilities                 : ["CAPABILITY_IAM"]
	Parameter overrides          : {"OuId": "YOUR-OU-ID-HERE", "SlackIncomingWebhookUrl": "YOUR-SLACK-INCOMING-WEBHOOK-URL-HERE"}
	Signing Profiles             : {}

Initiating deployment
=====================

	Uploading to aws-bills-on-slack/c3c2a4d368036749025894bf69b624dd.template  2645 / 2645  (100.00%)


Waiting for changeset to be created..

CloudFormation stack changeset
-----------------------------------------------------------------------------------------------------------------------------
Operation                       LogicalResourceId               ResourceType                    Replacement
-----------------------------------------------------------------------------------------------------------------------------
+ Add                           AwsBillsOnSlackFuncEveryDayWi   AWS::IAM::Role                  N/A
                                thTimeWindowRole
+ Add                           AwsBillsOnSlackFuncEveryDayWi   AWS::Scheduler::Schedule        N/A
                                thTimeWindow
+ Add                           AwsBillsOnSlackFuncRole         AWS::IAM::Role                  N/A
+ Add                           AwsBillsOnSlackFunc             AWS::Lambda::Function           N/A
-----------------------------------------------------------------------------------------------------------------------------


Changeset created successfully. arn:aws:cloudformation:ap-northeast-1:123456789012:changeSet/samcli-deploy1692275109/abcdefgh-1234-5678-ijkl-mnopqr901234


Previewing CloudFormation changeset before deployment
======================================================
Deploy this changeset? [y/N]: y

2023-08-17 21:25:20 - Waiting for stack create/update to complete

CloudFormation events from stack operations (refresh every 5.0 seconds)
-----------------------------------------------------------------------------------------------------------------------------
ResourceStatus                  ResourceType                    LogicalResourceId               ResourceStatusReason
-----------------------------------------------------------------------------------------------------------------------------
CREATE_IN_PROGRESS              AWS::CloudFormation::Stack      aws-bills-on-slack              User Initiated
CREATE_IN_PROGRESS              AWS::IAM::Role                  AwsBillsOnSlackFuncRole         -
CREATE_IN_PROGRESS              AWS::IAM::Role                  AwsBillsOnSlackFuncRole         Resource creation Initiated
CREATE_COMPLETE                 AWS::IAM::Role                  AwsBillsOnSlackFuncRole         -
CREATE_IN_PROGRESS              AWS::Lambda::Function           AwsBillsOnSlackFunc             -
CREATE_IN_PROGRESS              AWS::Lambda::Function           AwsBillsOnSlackFunc             Resource creation Initiated
CREATE_COMPLETE                 AWS::Lambda::Function           AwsBillsOnSlackFunc             -
CREATE_IN_PROGRESS              AWS::IAM::Role                  AwsBillsOnSlackFuncEveryDayWi   -
                                                                thTimeWindowRole
CREATE_IN_PROGRESS              AWS::IAM::Role                  AwsBillsOnSlackFuncEveryDayWi   Resource creation Initiated
                                                                thTimeWindowRole
CREATE_COMPLETE                 AWS::IAM::Role                  AwsBillsOnSlackFuncEveryDayWi   -
                                                                thTimeWindowRole
CREATE_IN_PROGRESS              AWS::Scheduler::Schedule        AwsBillsOnSlackFuncEveryDayWi   -
                                                                thTimeWindow
CREATE_IN_PROGRESS              AWS::Scheduler::Schedule        AwsBillsOnSlackFuncEveryDayWi   Resource creation Initiated
                                                                thTimeWindow
CREATE_COMPLETE                 AWS::Scheduler::Schedule        AwsBillsOnSlackFuncEveryDayWi   -
                                                                thTimeWindow
CREATE_COMPLETE                 AWS::CloudFormation::Stack      aws-bills-on-slack              -
-----------------------------------------------------------------------------------------------------------------------------

CloudFormation outputs from deployed stack
-------------------------------------------------------------------------------------------------------------------------------
Outputs
-------------------------------------------------------------------------------------------------------------------------------
Key                 AwsBillsOnSlackFuncArn
Description         aws-bills-on-slack Lambda Function ARN
Value               arn:aws:lambda:ap-northeast-1:123456789012:function:aws-bills-on-slack-AwsBillsOnSlackFunc-jgiBD64Wc6Xo
-------------------------------------------------------------------------------------------------------------------------------


Successfully created/updated stack - aws-bills-on-slack in ap-northeast-1
```
</details>

**For subsequent deployments:**

You may want to re-deploy once you made some changes on the application code or the template.yaml file. In that case, you can use the following command to deploy the changes. Note that if you want to change the paremeters which you entered in the first deployment, you need to edit the `./deploy/aws-sam/samconfig.toml` file before running the `sam deploy` command.

```shell
cd /path/to/your/aws-bills-on-slack
# edit the code or the template.yaml file
sam build --template ./deploy/aws-sam/template.yaml
sam deploy --config-file ./deploy/aws-sam/samconfig.toml --no-confirm-changeset # `--profile` and/or `--region` options may be also needed
```

### Install using CloudFormation

Todo

### Install using Terraform

Todo

## Contribution

1. Fork ([https://github.com/toricls/aws-bills-on-slack/fork](https://github.com/toricls/aws-bills-on-slack/fork))
4. Create a feature branch
5. Commit your changes
6. Rebase your local changes against the main branch
7. Create a new Pull Request (use [conventional commits] for the title please)

[conventional commits]: https://www.conventionalcommits.org/en/v1.0.0/

## Licence

Distributed under the [Apache-2.0](./LICENSE) license.

## Author

[Tori](https://github.com/toricls)
