

[Back](../README.md) | [Chat with us](https://anzoss.slack.com/messages/sysl-catalog/) | [New bug or feature request](https://github.com/anz-bank/sysl-catalog/issues/new)


# ServerPackage

## Integration Diagram
![](integration.svg)








## Database Index
| Database Application Name  | Source Location |
----|----
[RelModel](#Database-RelModel) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  


## Application Index
| Application Name | Method | Source Location |
----|----|----
Server | [Authenticate](#Server-Authenticate) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
Server | [GET /testRestPathParamPrimitive/{primitiveID}](#Server-GETtestRestPathParamPrimitive{primitiveID}) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
Server | [GET /testRestQueryParam](#Server-GETtestRestQueryParam) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
Server | [GET /testRestQueryParamPrimitive](#Server-GETtestRestQueryParamPrimitive) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
Server | [GET /testRestURLParam/{id}](#Server-GETtestRestURLParam{id}) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  
Server | [GET /testReturnNil](#Server-GETtestReturnNil) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|  

## Type Index
| Application Name | Type Name | Source Location |
----|----|----
Server | [Empty](#Server.Empty) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
Server | [Request](#Server.Request) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|
Server | [Response](#Server.Response) | [https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl](https:/github.com/anz-bank/sysl-catalog/blob/master/demo/simple2.sysl)|



# Databases



<details>
<summary>Database RelModel</summary>


![](RelModel/types.svg)
</details>




# Applications







## Application Server

- this is a comment for Server









### <a name=Server-Authenticate></a>Server Authenticate
this is a description of Authenticate

<details>
<summary>Sequence Diagram</summary>

![](Server/authenticate.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types





![](Server/request.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### <a name=Server-GETtestRestPathParamPrimitive{primitiveID}></a>Server GET /testRestPathParamPrimitive/{primitiveID}


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestrestpathparamprimitive{primitiveid}.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types








#### Path Parameter

![](primitive/stringsimple.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### <a name=Server-GETtestRestQueryParam></a>Server GET /testRestQueryParam


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestrestqueryparam.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types










#### Query Parameter

![](Server/request.svg)



#### Query Parameter

![](Server/request.svg)

</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### <a name=Server-GETtestRestQueryParamPrimitive></a>Server GET /testRestQueryParamPrimitive


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestrestqueryparamprimitive.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types










#### Query Parameter

![](primitive/stringsimple.svg)

</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### <a name=Server-GETtestRestURLParam{id}></a>Server GET /testRestURLParam/{id}


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestresturlparam{id}.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types








#### Path Parameter

![](Server/request.svg)



</details>
<details>
<summary>Response types</summary>

#### Response types





![](Server/response.svg)



</details>

---





### <a name=Server-GETtestReturnNil></a>Server GET /testReturnNil


<details>
<summary>Sequence Diagram</summary>

![](Server/gettestreturnnil.svg)
</details>

<details>
<summary>Request types</summary>

#### Request types







</details>
<details>
<summary>Response types</summary>

#### Response types



No Response Types


</details>

---




# Types





<details>
<summary>Server.Empty</summary>

### Server.Empty

- Empty Empty Empty

![](Server/emptysimple.svg)

[Full Diagram](Server/empty.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|

</details>
<details>
<summary>Server.Request</summary>

### Server.Request

- Request Request Request

![](Server/requestsimple.svg)

[Full Diagram](Server/request.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| query | sequence of Response | |

</details>
<details>
<summary>Server.Response</summary>

### Server.Response

- Response Response Response

![](Server/responsesimple.svg)

[Full Diagram](Server/response.svg)

#### Fields

| Field name | Type | Description |
|----|----|----|
| balance | MegaDatabase.Empty | |
| query | MegaDatabase.Money | |

</details>

<div class="footer">

