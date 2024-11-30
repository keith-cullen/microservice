*** Settings ***
Library     RequestsLibrary

Suite Setup     Create Session      alias=openapisession    verify=True    url=https://localhost

*** Variables ***
&{headers}    name=value

*** Test Cases ***
AppAPI/v1/getok: Get Request Success
    ${response}=    GET On Session      openapisession  url=/v1/get?name=Bob    headers=${headers}  expected_status=200
    Should Be Equal As Strings    hello Bob  ${response.json()}

AppAPI/v1/getinvalidmethod: Post Request Invalid Method
    ${response}=    POST On Session     openapisession  url=/v1/get?name=key1   headers=${headers}  expected_status=405

AppAPI/v1/getinvalidpath: Get Request Invalid Path
    ${response}=    Get On Session      openapisession  url=/v1/do?name=key1    headers=${headers}  expected_status=404

AppAPI/v1/getinvalidparam: Get Request Invalid Param
    ${response}=    Get On Session      openapisession  url=/v1/get?param=key1  headers=${headers}  expected_status=400
