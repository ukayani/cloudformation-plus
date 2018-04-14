# CloudFormation Plus

A command line tool to enable usage of more advanced YAML syntax in CloudFormation templates via pre-processing.

## How it works

This tool will read your YAML file with advance syntax unsupported by CloudFormation and produce an equivalent
YAML file which is supported by Cloud Formation.

**Note**: This is a general purpose YAML pre-processor, it is currently not specific to CloudFormation

## Features

- Supports use of YAML anchors and aliases
- Supports use of YAML merge keys (similar to an extends relationship in OOP)


## Example Using Aliases

To support the DRY principle, YAML supports the use of `aliases` and `anchors`.

Anchors provide a way to label nodes in your YAML document so they can be repeated later in your document.
Define anchors using the `&` character followed by the name of the label/anchor

Aliases provide a way to refer to a labeled node.
Define aliases by using the `*` character followed by the name of the anchor you wish to refer to

The below two documents are equivalent in their intended structure:

```yaml
foo: &anchor
 K1: "One"
 K2: "Two"

bar: *anchor
```

```yaml
foo:
 K1: "Hello"
 K2: "World"
bar:
 K1: "Hello"
 K2: "World" 
```

Since CloudFormation does not support the use of anchors we can convert from a document with anchors to one without them.

For example, below we avoid repeating our list of tags on two security group resources by defining an `anchor` labeled `tags`

```yaml
# Existing PARAMS [Project, VPCId, Team, Environment]
SecurityGroup1:
    Type: 'AWS::EC2::SecurityGroup'
    Properties:
      GroupName: 'My First Security Group'
      GroupDescription: 'Allows my instance external access'
      VpcId: !Ref 'VPCId'
      Tags: &tags
        - Key: Project
          Value: !Ref 'Project'
        - Key: Team
          Value: !Ref 'Team'
        - Key: Environment
          Value: !Ref 'Environment'
      SecurityGroupEgress:
        IpProtocol: -1
        CidrIp: "0.0.0.0/0"
SecurityGroup2:
    Type: 'AWS::EC2::SecurityGroup'
    Properties:
      GroupName: 'My First Security Group'
      GroupDescription: 'Allows access to my instance'
      SecurityGroupIngress:
        - CidrIp: "0.0.0.0/0"
          IpProtocol: tcp
          FromPort: 80
          ToPort: 80
      VpcId: !Ref 'VPCId'
      Tags: *tags

```

To convert from a YAML document with anchors/aliases to an equivalent document without we can run the following command:

```bash
$ cf-plus --resolve-aliases myfile.yml
``` 


# YAML Parsing Code

This tool uses code from [go-yaml](https://github.com/go-yaml/yaml) to parse/marshal YAML.

Since this tool needs to work directly with a YAML AST (which is not exposed by go-yaml), it modifies the go-yaml codebase:
 - adds an event initialization function for outputting alias nodes
 - adds additional fields to the AST node to preserve more information about the source document