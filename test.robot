*** Settings ***
Library     RequestsLibrary

Suite Setup     Create Session      alias=openapisession    verify=True    url=https://localhost:4443

*** Variables ***
&{headers}    name=value

*** Test Cases ***
AppAPI/v1/setok: Set API Success
    ${response}=    POST On Session     openapisession  url=/v1/set?name=Bob            headers=${headers}  expected_status=200
    Should Be Equal As Strings          {'message': 'Hello, Bob'}                       ${response.json()}

AppAPI/v1/getok: Get API Success
    ${response}=    GET On Session      openapisession  url=/v1/get?name=Bob            headers=${headers}  expected_status=200
    Should Be Equal As Strings          {'message': 'Hello, Bob'}                       ${response.json()}


AppAPI/v1/postinvalidpath: POST Request Invalid Path
    ${response}=    POST On Session     openapisession  url=/v1/x                       headers=${headers}  expected_status=404
    Should Be Equal As Strings          {'message': '404 Not Found'}                    ${response.json()}

AppAPI/v1/getinvalidpath: GET Request Invalid Path
    ${response}=    GET On Session      openapisession  url=/v1/x                       headers=${headers}  expected_status=404
    Should Be Equal As Strings          {'message': '404 Not Found'}                    ${response.json()}


AppAPI/v1/setwithunknownparameter: Set API with unknown parameter
    ${response}=    POST On Session     openapisession  url=/v1/set?unknown=unknown     headers=${headers}  expected_status=400
    Should Be Equal As Strings          {'message': '400 Bad Request'}                  ${response.json()}

AppAPI/v1/getwithunknownparameter: Get API with unknown parameter
    ${response}=    GET On Session      openapisession  url=/v1/get?unknown=unknown     headers=${headers}  expected_status=400
    Should Be Equal As Strings          {'message': '400 Bad Request'}                  ${response.json()}


AppAPI/v1/getwithunknownparametervalue: Get API with unknown parameter value
    ${response}=    GET On Session      openapisession  url=/v1/get?name=Unknown        headers=${headers}  expected_status=404
    Should Be Equal As Strings          {'message': '404 Not Found'}                    ${response.json()}


AppAPI/v1/setwithemptyparametervalue: Set API with empty parameter value
    ${response}=    POST On Session     openapisession  url=/v1/set?name=               headers=${headers}  expected_status=400
    Should Be Equal As Strings          {'message': '400 Bad Request'}                  ${response.json()}

AppAPI/v1/getwithemptyparametervalue: Get API with empty parameter value
    ${response}=    GET On Session      openapisession  url=/v1/get?name=               headers=${headers}  expected_status=400
    Should Be Equal As Strings          {'message': '400 Bad Request'}                  ${response.json()}
